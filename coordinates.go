package main

func StepSize(min, max float64, count int) float64 {
    return (max - min) / float64(count-1)
}

func PixelToComplex(
    x, y, width, height int,
    minReal, maxReal, minImag, maxImag float64,
) complex128 {
    realStepSize := StepSize(minReal, maxReal, width)
    imagStepSize := StepSize(minImag, maxImag, height)

    realPart := minReal + float64(x)*realStepSize
    imagPart := maxImag - float64(y)*imagStepSize

    return complex(realPart, imagPart)
}
