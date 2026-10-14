package main

import (
	"flag"
	"fmt"
	"math"
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

	presetName := flag.String("preset", "full", "Fractal location preset")
	listPresets := flag.Bool("list-presets", false, "List available fractal presets")

	paletteName := flag.String("palette", "classic", "Fractal color palette")
	listPalettes := flag.Bool("list-palettes", false, "List available color palettes")

	coloringName := flag.String(
		"coloring",
		"continuous",
		"Coloring mode: continuous or cyclic",
	)

	colorPeriod := flag.Float64(
		"color-period",
		32,
		"Smooth iterations per color cycle",
	)

	renderer := flag.String(
		"renderer",
		"unicode",
		"Rendering backend: unicode or kitty",
	)

	scale := flag.Int(
		"scale",
		8,
		"Kitty raster resolution multiplier (1-8)",
	)

	flag.Parse()

	if *renderer != "unicode" && *renderer != "kitty" {
		fmt.Fprintln(
			os.Stderr,
			"renderer must be unicode or kitty",
		)
		os.Exit(1)
	}

	if *scale < 1 || *scale > 8 {
		fmt.Fprintln(
			os.Stderr,
			"scale must be between 1 and 8",
		)
		os.Exit(1)
	}

	graphicsRenderer = *renderer
	graphicsScale = *scale
	graphicsInteractive = *interactive

	if *listPresets {
		PrintPresets()
		return
	}

	if *listPalettes {
		PrintPalettes()
		return
	}

	if *zoom <= 0 || math.IsNaN(*zoom) || math.IsInf(*zoom, 0) ||
		*cellAspect <= 0 || math.IsNaN(*cellAspect) || math.IsInf(*cellAspect, 0) {
		fmt.Fprintln(os.Stderr, "zoom and aspect must be finite and positive")
		os.Exit(1)
	}

	if math.IsNaN(*centerReal) || math.IsInf(*centerReal, 0) ||
		math.IsNaN(*centerImag) || math.IsInf(*centerImag, 0) {
		fmt.Fprintln(os.Stderr, "viewport center coordinates must be finite")
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

	if *colorPeriod <= 0 ||
		math.IsNaN(*colorPeriod) ||
		math.IsInf(*colorPeriod, 0) {
		fmt.Fprintln(os.Stderr, "color-period must be finite and positive")
		os.Exit(1)
	}

	mode := ColoringMode(*coloringName)

	if mode != ColoringContinuous && mode != ColoringCyclic {
		fmt.Fprintf(
			os.Stderr,
			"unknown coloring mode %q (use continuous or cyclic)\n",
			*coloringName,
		)
		os.Exit(1)
	}

	view, exists := GetPreset(*presetName)

	if !exists {
		fmt.Fprintf(os.Stderr, "unknown preset: %s\n", *presetName)
		fmt.Fprintln(os.Stderr, "Use -list-presets to see available locations")
		os.Exit(1)
	}

	// Override preset coordinates only when explicitly provided.
	flag.Visit(func(f *flag.Flag) {
		switch f.Name {
		case "x":
			view.CenterReal = *centerReal
		case "y":
			view.CenterImag = *centerImag
		}
	})

	view = Zoom(view, *zoom)

	palette, exists := GetPalette(*paletteName)

	if !exists {
		fmt.Fprintf(os.Stderr, "unknown palette: %s\n", *paletteName)
		fmt.Fprintln(os.Stderr, "Use -list-palettes to see available palettes")
		os.Exit(1)
	}

	// Store coloring settings in Palette.
	// The renderer signature stays unchanged.
	palette.Mode = mode
	palette.Period = *colorPeriod

	if *canvasWidth == 0 || *canvasHeight == 0 {
		columns, rows, err := term.GetSize(int(os.Stdout.Fd()))

		if err != nil {
			fmt.Fprintln(
				os.Stderr,
				"could not detect terminal size; specify -width and -height",
			)
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
		fmt.Fprintln(os.Stderr, "canvas dimensions are too small")
		os.Exit(1)
	}

	if *interactive {
		var err error

		if graphicsRenderer == "kitty" {
			err = RunSmartKitty(
				view,
				*canvasWidth,
				*canvasHeight,
				*cellAspect,
				*maxIter,
				palette,
			)
		} else {
			err = RunInteractive(
				view,
				*canvasWidth,
				*canvasHeight,
				*cellAspect,
				*maxIter,
				palette,
			)
		}

		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}

		return
	}

	RenderGrid(
		*canvasWidth,
		*canvasHeight,
		*cellAspect,
		view,
		*maxIter,
		palette,
	)
}
