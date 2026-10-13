package main

import (
	"flag"
	"fmt"
	"os"

	"golang.org/x/term"
)

func main() {
	zoom := flag.Float64("zoom", 1.0, "Zoom multiplier")
	centerReal := flag.Float64("x", -0.9, "Real coordinate of viewport center")
	centerImag := flag.Float64("y", 0.0, "Imaginary coordinate of viewport center")

	canvasWidth := flag.Int("width", 0, "Canvas width (0 = automatic)")
	canvasHeight := flag.Int("height", 0, "Canvas height (0 = automatic)")

	cellAspect := flag.Float64("aspect", 2.0, "Terminal cell height-to-width ratio")
	maxIter := flag.Int("iter", 10000, "Maximum Mandelbrot iterations")

	interactive := flag.Bool("interactive", false, "Enable interactive navigation")

	flag.Parse()

	if *zoom <= 0 || *cellAspect <= 0 {
		fmt.Fprintln(os.Stderr, "zoom and aspect must be positive")
		os.Exit(1)
	}

	if *maxIter <= 0 {
		fmt.Fprintln(os.Stderr, "iter must be positive")
		os.Exit(1)
	}

	if *canvasWidth < 0 || *canvasHeight < 0 {
		fmt.Fprintln(os.Stderr, "canvas dimensions cannot be negative")
		os.Exit(1)
	}

	// Detect terminal dimensions
	if *canvasWidth == 0 || *canvasHeight == 0 {
		columns, rows, err := term.GetSize(int(os.Stdout.Fd()))

		if err != nil {
			fmt.Fprintln(os.Stderr, "Could not detect terminal size")
			fmt.Fprintln(os.Stderr, "Specify -width and -height manually")
			os.Exit(1)
		}

		if *canvasWidth == 0 {
			*canvasWidth = columns - 1
		}

		if *canvasHeight == 0 {
			*canvasHeight = rows - 2
		}
	}

	if *canvasWidth < 2 || *canvasHeight < 1 {
		fmt.Fprintln(os.Stderr, "Canvas dimensions are too small")
		os.Exit(1)
	}

	view := Viewport{*centerReal, *centerImag, 5.2}
	view = Zoom(view, *zoom)

	if *interactive {
		err := RunInteractive(
			view,
			*canvasWidth,
			*canvasHeight,
			*cellAspect,
			*maxIter,
		)

		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}

		return
	}

	RenderGrid(*canvasWidth, *canvasHeight, *cellAspect, view, *maxIter)
}
