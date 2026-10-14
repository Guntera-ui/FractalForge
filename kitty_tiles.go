package main

import (
	"bytes"
	"context"
	"errors"
	"image"
	"image/color"
	"image/png"
	"math"
	"runtime"
	"sort"
	"sync"
	"sync/atomic"
)

type fractalTile struct {
	id                      int
	col, row                int
	cols, rows              int
	px, py                  int
	pixelWidth, pixelHeight int
}

type tileConfig struct {
	columns, rows, scale int
	aspect               float64
	view                 Viewport
	maxIter              int
	palette              Palette
	canvas               *image.RGBA
}

type finishedTile struct {
	jobID uint64
	tile  fractalTile
	png   []byte
	err   error
}

func makeFractalTiles(
	columns, rows, scale int,
	aspect float64,
) []fractalTile {
	const tileColumns = 10
	const tileRows = 5

	var tiles []fractalTile

	for row := 0; row < rows; row += tileRows {
		for col := 0; col < columns; col += tileColumns {
			cw := min(tileColumns, columns-col)
			rh := min(tileRows, rows-row)

			y0 := int(math.Round(float64(row*scale) * aspect))
			y1 := int(math.Round(float64((row+rh)*scale) * aspect))

			tiles = append(tiles, fractalTile{
				id:          len(tiles),
				col:         col,
				row:         row,
				cols:        cw,
				rows:        rh,
				px:          col * scale,
				py:          y0,
				pixelWidth:  cw * scale,
				pixelHeight: y1 - y0,
			})
		}
	}

	// Calculate the center of the viewport first.
	sort.Slice(tiles, func(i, j int) bool {
		a, b := tiles[i], tiles[j]

		ax := float64(a.col) + float64(a.cols)/2 -
			float64(columns)/2
		ay := float64(a.row) + float64(a.rows)/2 -
			float64(rows)/2

		bx := float64(b.col) + float64(b.cols)/2 -
			float64(columns)/2
		by := float64(b.row) + float64(b.rows)/2 -
			float64(rows)/2

		da := ax*ax + ay*ay
		db := bx*bx + by*by

		if da == db {
			return a.id < b.id
		}

		return da < db
	})

	return tiles
}

// Calculate tiles concurrently using available CPU workers.
func calculateTiles(
	ctx context.Context,
	jobID uint64,
	tiles []fractalTile,
	cfg tileConfig,
	finished chan<- finishedTile,
) {
	var next atomic.Int64
	var wg sync.WaitGroup

	workerCount := min(
		runtime.GOMAXPROCS(0),
		len(tiles),
	)

	for w := 0; w < workerCount; w++ {
		wg.Add(1)

		go func() {
			defer wg.Done()

			for {
				if ctx.Err() != nil {
					return
				}

				index := int(next.Add(1)) - 1

				if index >= len(tiles) {
					return
				}

				tile := tiles[index]

				data, err := renderOneTile(
					ctx,
					tile,
					cfg,
				)

				if errors.Is(err, context.Canceled) {
					return
				}

				select {
				case finished <- finishedTile{
					jobID: jobID,
					tile:  tile,
					png:   data,
					err:   err,
				}:
				case <-ctx.Done():
					return
				}

				if err != nil {
					return
				}
			}
		}()
	}

	wg.Wait()
}

// Render one tile while also filling the shared canvas.
// Each worker owns a separate region of the canvas.
func renderOneTile(
	ctx context.Context,
	tile fractalTile,
	cfg tileConfig,
) ([]byte, error) {
	width := cfg.columns * cfg.scale

	height := int(math.Round(
		float64(cfg.rows*cfg.scale) * cfg.aspect,
	))

	if width < 2 || height < 2 ||
		tile.pixelWidth < 1 || tile.pixelHeight < 1 {
		return nil, errors.New("invalid tile dimensions")
	}

	img := image.NewRGBA(
		image.Rect(
			0, 0,
			tile.pixelWidth,
			tile.pixelHeight,
		),
	)

	// IMPORTANT:
	// Viewport spans must not depend on render scale.
	// This prevents shifting between preview and full quality.

	realSpan := cfg.view.Width

	imagSpan := realSpan *
		float64(cfg.rows) *
		cfg.aspect /
		float64(cfg.columns)

	realStep := realSpan / float64(width)
	imagStep := imagSpan / float64(height)

	minReal := cfg.view.CenterReal - realSpan/2
	maxImag := cfg.view.CenterImag + imagSpan/2

	for y := 0; y < tile.pixelHeight; y++ {
		if err := ctx.Err(); err != nil {
			return nil, err
		}

		// Position within the complete raster.
		iy := tile.py + y

		// Sample at the center of each pixel.
		imag := maxImag -
			(float64(iy)+0.5)*imagStep

		for x := 0; x < tile.pixelWidth; x++ {
			if x%16 == 0 {
				if err := ctx.Err(); err != nil {
					return nil, err
				}
			}

			ix := tile.px + x

			c := complex(
				minReal+(float64(ix)+0.5)*realStep,
				imag,
			)

			result := MandelbrotDetailed(
				c,
				cfg.maxIter,
			)

			pixel := ColorFromPalette(
				result,
				cfg.palette,
				cfg.maxIter,
			)

			rgba := color.RGBA{
				R: uint8(max(0, min(255, pixel.R))),
				G: uint8(max(0, min(255, pixel.G))),
				B: uint8(max(0, min(255, pixel.B))),
				A: 255,
			}

			// Temporary tile for progressive display.
			img.SetRGBA(x, y, rgba)

			// Full-resolution canvas for seamless display.
			// Tiles occupy non-overlapping rectangles.
			if cfg.canvas != nil {
				cfg.canvas.SetRGBA(ix, iy, rgba)
			}
		}
	}

	var buffer bytes.Buffer

	encoder := png.Encoder{
		CompressionLevel: png.BestSpeed,
	}

	if err := encoder.Encode(&buffer, img); err != nil {
		return nil, err
	}

	if err := ctx.Err(); err != nil {
		return nil, err
	}

	return buffer.Bytes(), nil
}
