package main

import (
	"bufio"
	"bytes"
	"encoding/base64"
	"fmt"
	"image"
	"image/color"
	"image/png"
	"io"
	"math"
	"os"
)

// Graphics settings. The Unicode renderer remains the default.
var graphicsRenderer = "unicode"
var graphicsScale = 2
var graphicsInteractive bool

func RenderKitty(
	columns, rows int,
	cellAspect float64,
	view Viewport,
	maxIter int,
	palette Palette,
) error {
	if columns < 2 || rows < 1 ||
		graphicsScale < 1 || graphicsScale > 8 {
		return fmt.Errorf("invalid graphics dimensions or scale")
	}

	// More samples per terminal character cell.
	pixelWidth := columns * graphicsScale
	pixelHeight := int(math.Round(
		float64(rows*graphicsScale) * cellAspect,
	))

	if pixelHeight < 2 || pixelWidth < 2 ||
		int64(pixelWidth)*int64(pixelHeight) > 8_000_000 {
		return fmt.Errorf(
			"image too large; reduce -scale, -width or -height",
		)
	}

	img := image.NewRGBA(
		image.Rect(0, 0, pixelWidth, pixelHeight),
	)

	// Preserve complex-plane proportions using square
	// raster samples instead of terminal character samples.
	minReal := view.CenterReal - view.Width/2
	realStep := view.Width / float64(pixelWidth-1)

	imagRange := realStep * float64(pixelHeight-1)
	maxImag := view.CenterImag + imagRange/2

	for y := 0; y < pixelHeight; y++ {
		imag := maxImag - float64(y)*realStep

		for x := 0; x < pixelWidth; x++ {
			c := complex(
				minReal+float64(x)*realStep,
				imag,
			)

			result := MandelbrotDetailed(c, maxIter)
			pixel := ColorFromPalette(
				result,
				palette,
				maxIter,
			)

			img.SetRGBA(x, y, color.RGBA{
				R: uint8(max(0, min(255, pixel.R))),
				G: uint8(max(0, min(255, pixel.G))),
				B: uint8(max(0, min(255, pixel.B))),
				A: 255,
			})
		}
	}

	// Encode the image as PNG entirely in memory.
	var data bytes.Buffer

	encoder := png.Encoder{
		CompressionLevel: png.BestSpeed,
	}

	if err := encoder.Encode(&data, img); err != nil {
		return fmt.Errorf("encode PNG: %w", err)
	}

	out := bufio.NewWriter(os.Stdout)

	if graphicsInteractive {
		// Anchor each frame at the upper-left corner.
		if _, err := out.WriteString("\x1b[H"); err != nil {
			return err
		}
	} else {
		// Reserve enough terminal rows for the image
		// so the shell prompt appears below it.
		for i := 0; i < rows; i++ {
			if _, err := out.WriteString("\r\n"); err != nil {
				return err
			}
		}

		if _, err := fmt.Fprintf(
			out, "\x1b[%dA\r", rows,
		); err != nil {
			return err
		}
	}

	// Use a process-specific image ID.
	id := uint32(os.Getpid())
	if id == 0 {
		id = 1
	}

	// Delete our previous image before redrawing.
	if _, err := fmt.Fprintf(
		out,
		"\x1b_Ga=d,d=I,i=%d,q=2;\x1b\\",
		id,
	); err != nil {
		return err
	}

	if err := writeKittyPNG(
		out,
		data.Bytes(),
		columns,
		rows,
		id,
	); err != nil {
		return err
	}

	if !graphicsInteractive {
		if _, err := fmt.Fprintf(
			out, "\x1b[%dB\r", rows,
		); err != nil {
			return err
		}
	}

	return out.Flush()
}

// Send the PNG directly to the terminal in Base64 chunks.
// Each chunk is at most 4096 bytes.
func writeKittyPNG(
	w io.Writer,
	pngBytes []byte,
	columns, rows int,
	id uint32,
) error {
	encoded := base64.StdEncoding.EncodeToString(pngBytes)

	for start := 0; start < len(encoded); start += 4096 {
		end := min(start+4096, len(encoded))

		more := 0
		if end < len(encoded) {
			more = 1
		}

		control := fmt.Sprintf("m=%d,q=2", more)

		if start == 0 {
			control = fmt.Sprintf(
				"a=T,f=100,t=d,i=%d,c=%d,r=%d,C=1,m=%d,q=2",
				id,
				columns,
				rows,
				more,
			)
		}

		_, err := fmt.Fprintf(
			w,
			"\x1b_G%s;%s\x1b\\",
			control,
			encoded[start:end],
		)

		if err != nil {
			return err
		}
	}

	return nil
}
