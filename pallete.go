package main

import "math"

type RGB struct {
	R int
	G int
	B int
}
func Lerp(start,end,t float64) float64 {
	// R = start + t(end - start)
	interpolated := start + (t*(end - start))
	return interpolated
}

func LerpColor(start, end RGB, t float64) RGB {
    R := int(math.Round(Lerp(float64(start.R), float64(end.R), t)))
    G := int(math.Round(Lerp(float64(start.G), float64(end.G), t)))
    B := int(math.Round(Lerp(float64(start.B), float64(end.B), t)))

    result := RGB{R, G, B}

    return result
}
	

