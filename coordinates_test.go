package main

import "testing"

func TestPixelToComplex(t *testing.T) {
    result := PixelToComplex(
        0, 0,       // x, y
        5, 5,       // width, height
        -2, 2,      // minReal, maxReal
        -2, 2,      // minImag, maxImag
    )

	result2 := PixelToComplex(
		4,0,
		5,5,
		-2,2,
		-2,2,
	)
	
	result3:= PixelToComplex(
		2,2,
		5,5,
		-2,2,
		-2,2,
	)

	result4 := PixelToComplex(
		4,4,
		5,5,
		-2,2,
		-2,2,
	)
    
	expected := complex(-2.0, 2.0)
	expected2 := complex(2.0,2.0)
	expected3 := complex(0.0,0.0)
	expected4 := complex(2.0,-2.0)


    if result != expected {
        t.Errorf("expected %v, got %v", expected, result)
    }
    if result2 != expected2 {
        t.Errorf("expected %v, got %v", expected, result)
	}
	if result3 != expected3 {
        t.Errorf("expected %v, got %v", expected, result)
	}
	if result4 != expected4 {
        t.Errorf("expected %v, got %v", expected, result)
	}
}
