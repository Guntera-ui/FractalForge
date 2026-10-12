package main

import "testing"

func TestMandelbrot(t *testing.T) {
	result := Mandelbrot(3+0i,100)
	result2 := Mandelbrot(1+0i,100)
	result3 := Mandelbrot(0+0i, 100)
	result4 := Mandelbrot(-1+0i,100)
	result5 := Mandelbrot(2+0i,100)
	if result != 1 {
		t.Errorf("expected 1, got %d", result)
	}
	if result2 != 3 {
		t.Errorf("expected 3, got %d", result2)
	}
	if result3 != 100 {
		t.Errorf("expected 100, got %d", result3)
	}
	if result4 != 100 {
		t.Errorf("expected 100, got %d", result4)
	}
	if result5 != 2 {
		t.Errorf("expected 1, got %d", result5)
	}





}
