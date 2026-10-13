package main

type Viewport struct {
	CenterReal float64
	CenterImag float64
	Width      float64
}

func RealBounds(center, width float64) (float64, float64) {
	minReal := center - (width / 2)
	maxReal := center + (width / 2)

	return minReal, maxReal
}

func Zoom(v Viewport, factor float64) Viewport {
	newWidth := v.Width / factor
	return Viewport{v.CenterReal, v.CenterImag, newWidth}
}

func Pan(v Viewport, deltaReal, deltaImag float64) Viewport {
	newCenterReal := v.CenterReal + deltaReal
	newCenterImag := v.CenterImag + deltaImag
	return Viewport{newCenterReal, newCenterImag, v.Width}

}
