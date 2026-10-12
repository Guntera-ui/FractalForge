package main


func Mandelbrot(c complex128, maxIter int) int {
	z := complex(0.0,0.0) // first z always equals to 0 btw
	for n := 0; n < maxIter; n++ {
		z = z*z + c

		realPart := real(z)
		imagPart := imag(z)
		
		magnitudeSquared := realPart*realPart + imagPart*imagPart
		
		if magnitudeSquared > 4{
			return n + 1
		}

	}

	return maxIter
}
