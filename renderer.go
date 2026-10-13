package main

import "fmt"

func RenderGrid(width, height int, cellAspect float64) {
	//Some questionable code
	minReal := -3.5
	maxReal := 1.7
	maxIter := 1000000

	realStep := StepSize(minReal, maxReal, width)
	verticalStep := realStep * cellAspect / 2

	sampleHeight := height * 2
	imagRange := verticalStep * float64(sampleHeight-1)

	minImag := -imagRange / 2
	maxImag := imagRange / 2

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

			topColor := ColorFromClassicPalette(topResult)
			bottomColor := ColorFromClassicPalette(bottomResult)

			PrintHalfBlock(topColor, bottomColor)
		}

		fmt.Println()
	}

	fmt.Print("\x1b[0m")
	//End of miy questionable code
}

func PrintHalfBlock(top, bottom RGB) {
	fmt.Printf(
		"\x1b[38;2;%d;%d;%dm\x1b[48;2;%d;%d;%dm▀", //this is chatgpt im not crazy
		top.R, top.G, top.B,
		bottom.R, bottom.G, bottom.B,
	)
}
