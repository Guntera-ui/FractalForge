package main

import "math"
import "math/cmplx"

type EscapeResult struct {
	Iterations int
	FinalZ     complex128
	Escaped    bool
}

func MandelbrotDetailed(c complex128, maxIter int) EscapeResult {
	z := complex(0.0, 0.0)

	for n := 0; n < maxIter; n++ {
		z = z*z + c

		realPart := real(z)
		imagPart := imag(z)

		magnitudeSquared := realPart*realPart + imagPart*imagPart

		if magnitudeSquared > 4 {
			return EscapeResult{
				Iterations: n + 1,
				FinalZ:     z,
				Escaped:    true,
			}
		}
	}

	return EscapeResult{
		Iterations: maxIter,
		FinalZ:     z,
		Escaped:    false,
	}
}

func SmoothIteration(result EscapeResult) float64 {
	n := float64(result.Iterations)
	z := result.FinalZ

	if !result.Escaped {
		return float64(result.Iterations)
	}

	magnitude := cmplx.Abs(z)
	innerLog := math.Log2(magnitude)
	outerLog := math.Log2(innerLog)

	smoothValue := n + 1 - outerLog

	return smoothValue
}

func Mandelbrot(c complex128, maxIter int) int {
	result := MandelbrotDetailed(c, maxIter)
	return result.Iterations
}
