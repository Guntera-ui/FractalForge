package main

import "math"

type RGB struct {
	R int
	G int
	B int
}

func Lerp(start, end, t float64) float64 {
	// R = start + t(end - start)
	interpolated := start + (t * (end - start))
	return interpolated
}

func LerpColor(start, end RGB, t float64) RGB {
	R := int(math.Round(Lerp(float64(start.R), float64(end.R), t)))
	G := int(math.Round(Lerp(float64(start.G), float64(end.G), t)))
	B := int(math.Round(Lerp(float64(start.B), float64(end.B), t)))

	result := RGB{R, G, B}

	return result
}

func ColorFromIterations(
	iterations, maxIter int,
	startColor, endColor RGB,
) RGB {
	if iterations == maxIter {
		return RGB{4, 8, 20}
	}

	t := math.Log1p(float64(iterations)) /
		math.Log1p(float64(maxIter))

	return LerpColor(startColor, endColor, t)
}

func ColorFromClassicPalette(result EscapeResult) RGB {
	if !result.Escaped {
		return RGB{0, 0, 0}
	}

	iterations := SmoothIteration(result)

	colors := []RGB{
		{94, 68, 239}, // Violet-blue
		{101, 52, 222},
		{109, 35, 206},
		{117, 19, 190},
		{126, 4, 175},
		{134, 10, 158},
		{141, 27, 141},
		{159, 58, 111},
		{171, 88, 75},
		{198, 145, 58},
		{222, 193, 41},
		{255, 218, 119},
		{255, 241, 183},
		{182, 207, 255},
	}

	t := math.Log2(math.Max(1, iterations)) / 7.0
	t = math.Max(0, math.Min(1, t))

	position := t * float64(len(colors)-1)
	index := int(position)

	if index >= len(colors)-1 {
		return colors[len(colors)-1]
	}

	fraction := position - float64(index)

	return LerpColor(
		colors[index],
		colors[index+1],
		fraction,
	)
}
