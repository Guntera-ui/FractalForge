package main

import "fmt"

func RenderGrid(width, height int, cellAspect float64) {
	//Something goes here
	minReal := -2.0
	maxReal := 1.0

	realStep := StepSize(minReal, maxReal, width)

	verticalStep := realStep * cellAspect

	imagRange := verticalStep * float64(height-1)


	minImag := -imagRange / 2
	maxImag := imagRange / 2

	for y := 0; y < height; y++{
		for x := 0; x < width; x++{
			c := PixelToComplex(
				x,y,
				width, height,
				minReal,maxReal,
				minImag,maxImag,
		)
			iterations := Mandelbrot(c,100)

			if iterations <= 9 {
				PrintColorCell(59,82,139)
			} else if iterations <= 29 {
				PrintColorCell(155,101,213)
			} else if iterations <= 99 {
				PrintColorCell(233,154,202)
			} else {
				PrintColorCell(15,23,42)
			}
		}
		fmt.Println()
	}
}

func PrintColorCell(r,g,b int) {
	fmt.Printf("\x1b[48;2;%d;%d;%dm ", r, g, b)
}
