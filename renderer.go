package main

import "fmt"
import "os"

func RenderGrid(width, height int, cellAspect float64, view Viewport, maxIter int, palette Palette) {

	if graphicsRenderer == "kitty" {
		if err := RenderKitty(
			width,
			height,
			cellAspect,
			view,
			maxIter,
			palette,
		); err != nil {
			fmt.Fprintln(os.Stderr, "Kitty renderer:", err)
		}

		return
	}

	minReal, maxReal := RealBounds(view.CenterReal, view.Width)

	realStep := StepSize(minReal, maxReal, width)
	verticalStep := realStep * cellAspect / 2

	sampleHeight := height * 2
	imagRange := verticalStep * float64(sampleHeight-1)

	minImag := view.CenterImag - imagRange/2
	maxImag := view.CenterImag + imagRange/2

	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {

			topC := PixelToComplex(
				x, 2*y,
				width, sampleHeight,
				minReal, maxReal,
				minImag, maxImag,
			)

			bottomC := PixelToComplex(
				x, 2*y+1,
				width, sampleHeight,
				minReal, maxReal,
				minImag, maxImag,
			)

			topResult := MandelbrotDetailed(topC, maxIter)
			bottomResult := MandelbrotDetailed(bottomC, maxIter)

			topColor := ColorFromPalette(topResult, palette, maxIter)
			bottomColor := ColorFromPalette(bottomResult, palette, maxIter)

			PrintHalfBlock(topColor, bottomColor)
		}

		fmt.Print("\r\n")
	}

	fmt.Print("\x1b[0m")
}

func PrintHalfBlock(top, bottom RGB) {
	fmt.Printf(
		"\x1b[38;2;%d;%d;%dm\x1b[48;2;%d;%d;%dm▀", //this is chatgpt im not crazy
		top.R, top.G, top.B,
		bottom.R, bottom.G, bottom.B,
	)
}
