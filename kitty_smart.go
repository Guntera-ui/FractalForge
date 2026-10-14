package main

import (
	"bufio"
	"context"
	"errors"
	"fmt"
	"image"
	"math"
	"os"
	"time"

	"golang.org/x/term"
)

func RunSmartKitty(
	view Viewport,
	columns, rows int,
	aspect float64,
	maxIter int,
	palette Palette,
) error {
	if columns < 2 || rows < 1 ||
		graphicsScale < 1 || graphicsScale > 8 ||
		aspect <= 0 ||
		math.IsNaN(aspect) ||
		math.IsInf(aspect, 0) {
		return errors.New("invalid Kitty renderer dimensions")
	}

	fullWidth := columns * graphicsScale

	fullHeight := int(math.Round(
		float64(rows*graphicsScale) * aspect,
	))

	if fullWidth < 2 || fullHeight < 2 ||
		int64(fullWidth)*int64(fullHeight) > 8_000_000 {
		return errors.New(
			"Kitty image too large; reduce scale or terminal dimensions",
		)
	}

	fd := int(os.Stdin.Fd())

	if !term.IsTerminal(fd) {
		return errors.New("interactive mode requires a terminal")
	}

	oldState, err := term.MakeRaw(fd)
	if err != nil {
		return err
	}
	defer term.Restore(fd, oldState)

	out := bufio.NewWriter(os.Stdout)

	if _, err := out.WriteString(
		"\x1b[?1049h\x1b[?25l\x1b[H",
	); err != nil {
		return err
	}

	if err := out.Flush(); err != nil {
		return err
	}

	const previewLimit = 4

	previewScale := min(previewLimit, graphicsScale)

	previewTiles := makeFractalTiles(
		columns, rows, previewScale, aspect,
	)

	fullTiles := makeFractalTiles(
		columns, rows, graphicsScale, aspect,
	)

	idBase := (uint32(os.Getpid()) & 0xffff) << 16

	imageID := func(index int) uint32 {
		return idBase + uint32(index+1)
	}

	fullImageID := imageID(len(fullTiles))

	clearTiles := func() error {
		for _, t := range fullTiles {
			if _, err := fmt.Fprintf(
				out,
				"\x1b_Ga=d,d=I,i=%d,q=2;\x1b\\",
				imageID(t.id),
			); err != nil {
				return err
			}
		}

		return out.Flush()
	}

	clearAll := func() error {
		if err := clearTiles(); err != nil {
			return err
		}

		if _, err := fmt.Fprintf(
			out,
			"\x1b_Ga=d,d=I,i=%d,q=2;\x1b\\",
			fullImageID,
		); err != nil {
			return err
		}

		return out.Flush()
	}

	defer func() {
		_ = clearAll()

		_, _ = out.WriteString(
			"\x1b[?25h\x1b[?1049l",
		)

		_ = out.Flush()
	}()

	keys := make(chan byte, 128)
	readErrors := make(chan error, 1)
	stopReading := make(chan struct{})

	defer close(stopReading)

	go func() {
		var buf [64]byte

		for {
			n, readErr := os.Stdin.Read(buf[:])

			for _, key := range buf[:n] {
				select {
				case keys <- key:
				case <-stopReading:
					return
				}
			}

			if readErr != nil {
				select {
				case readErrors <- readErr:
				case <-stopReading:
				}
				return
			}
		}
	}()

	results := make(chan finishedTile, 64)
	composed := make(chan completedFrame, 2)

	idleTimer := time.NewTimer(time.Hour)

	if !idleTimer.Stop() {
		select {
		case <-idleTimer.C:
		default:
		}
	}

	defer idleTimer.Stop()

	var idle <-chan time.Time
	var cancel context.CancelFunc

	defer func() {
		if cancel != nil {
			cancel()
		}
	}()

	var jobID uint64

	var currentCanvas *image.RGBA
	var currentCtx context.Context

	quality := ""
	finishedCount := 0
	totalTiles := 1
	idleReached := false

	showStatus := func() error {
		progress := 100 * finishedCount /
			max(1, totalTiles)

		status := fmt.Sprintf(
			" FractalForge | %s %d%% | x:%+.6f y:%+.6f | span:%.4g | iter:%d | WASD +/- Q",
			quality,
			progress,
			view.CenterReal,
			view.CenterImag,
			view.Width,
			maxIter,
		)

		if len(status) > columns {
			status = status[:columns]
		}

		if _, err := fmt.Fprintf(
			out,
			"\x1b[%d;1H\x1b[0m\x1b[2K\x1b[7m%-*s\x1b[0m",
			rows+1,
			columns,
			status,
		); err != nil {
			return err
		}

		return out.Flush()
	}

	launch := func(
		stage string,
		scale, iterations int,
	) error {
		if cancel != nil {
			cancel()
		}

		ctx, stop := context.WithCancel(
			context.Background(),
		)

		cancel = stop
		currentCtx = ctx

		jobID++
		quality = stage
		finishedCount = 0

		tiles := fullTiles

		if stage == "preview" {
			tiles = previewTiles
		}

		totalTiles = len(tiles)
		currentCanvas = nil

		if stage == "full" {
			currentCanvas = image.NewRGBA(
				image.Rect(
					0, 0,
					fullWidth,
					fullHeight,
				),
			)
		}

		cfg := tileConfig{
			columns: columns,
			rows:    rows,
			scale:   scale,
			aspect:  aspect,
			view:    view,
			maxIter: iterations,
			palette: palette,
			canvas:  currentCanvas,
		}

		go calculateTiles(
			ctx,
			jobID,
			tiles,
			cfg,
			results,
		)

		return showStatus()
	}

	restart := func() error {
		if cancel != nil {
			cancel()
		}

		if err := clearAll(); err != nil {
			return err
		}

		idleReached = false

		if !idleTimer.Stop() {
			select {
			case <-idleTimer.C:
			default:
			}
		}

		idleTimer.Reset(250 * time.Millisecond)
		idle = idleTimer.C

		return launch(
			"preview",
			previewScale,
			min(250, maxIter),
		)
	}

	if err := restart(); err != nil {
		return err
	}

	for {
		select {
		case key := <-keys:
			changed := false

		drain:
			for {
				if key == 'q' ||
					key == 'Q' ||
					key == 3 {
					return nil
				}

				nextView, keepRunning := HandleKey(
					view,
					key,
				)

				if !keepRunning {
					return nil
				}

				if nextView != view {
					view = nextView
					changed = true
				}

				select {
				case key = <-keys:
				default:
					break drain
				}
			}

			if changed {
				if err := restart(); err != nil {
					return err
				}
			}

		case <-idle:
			idle = nil
			idleReached = true

			if quality == "preview" &&
				finishedCount == totalTiles {
				if err := launch(
					"full",
					graphicsScale,
					maxIter,
				); err != nil {
					return err
				}
			}

		case tile := <-results:
			if tile.jobID != jobID {
				continue
			}

			if tile.err != nil {
				if errors.Is(
					tile.err,
					context.Canceled,
				) {
					continue
				}

				return tile.err
			}

			id := imageID(tile.tile.id)

			if _, err := fmt.Fprintf(
				out,
				"\x1b_Ga=d,d=I,i=%d,q=2;\x1b\\",
				id,
			); err != nil {
				return err
			}

			if _, err := fmt.Fprintf(
				out,
				"\x1b[%d;%dH",
				tile.tile.row+1,
				tile.tile.col+1,
			); err != nil {
				return err
			}

			if err := writeKittyPNG(
				out,
				tile.png,
				tile.tile.cols,
				tile.tile.rows,
				id,
			); err != nil {
				return err
			}

			if err := out.Flush(); err != nil {
				return err
			}

			finishedCount++

			if finishedCount == totalTiles ||
				finishedCount%max(
					1, totalTiles/10,
				) == 0 {
				if err := showStatus(); err != nil {
					return err
				}
			}

			if quality == "preview" &&
				finishedCount == totalTiles &&
				idleReached {
				if err := launch(
					"full",
					graphicsScale,
					maxIter,
				); err != nil {
					return err
				}
			}

			if quality == "full" &&
				finishedCount == totalTiles {
				snapshot := currentCanvas
				renderCtx := currentCtx
				finishedJob := jobID

				go func() {
					data, encodeErr := encodeWholeFrame(
						renderCtx,
						snapshot,
					)

					select {
					case composed <- completedFrame{
						jobID: finishedJob,
						png:   data,
						err:   encodeErr,
					}:
					case <-renderCtx.Done():
					}
				}()
			}

		case frame := <-composed:
			if frame.jobID != jobID ||
				quality != "full" {
				continue
			}

			if frame.err != nil {
				if errors.Is(
					frame.err,
					context.Canceled,
				) {
					continue
				}

				return frame.err
			}

			if _, err := out.WriteString(
				"\x1b[H",
			); err != nil {
				return err
			}

			if err := writeKittyPNG(
				out,
				frame.png,
				columns,
				rows,
				fullImageID,
			); err != nil {
				return err
			}

			if err := out.Flush(); err != nil {
				return err
			}

			if err := clearTiles(); err != nil {
				return err
			}

			if err := showStatus(); err != nil {
				return err
			}

		case readErr := <-readErrors:
			return fmt.Errorf(
				"keyboard input: %w",
				readErr,
			)
		}
	}
}
