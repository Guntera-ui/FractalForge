package main

import "testing"

func TestLerp(t *testing.T) {
	result := Lerp(50,150,0)
	result2 := Lerp(50,150,0.2)
	result3 := Lerp(50,150,0.5)
	result4 := Lerp(200,100,0.5)
	result5 := Lerp(50,150,1)

	if result != 50 {
		t.Errorf("expected 50, got %f", result)
	}
	if result2 != 70 {
		t.Errorf("expected 70, got %f", result2)
	}
	if result3 != 100 {
		t.Errorf("expected 100, got %f", result3)
	}
	if result4 != 150 {
		t.Errorf("expected 150, got %f", result4)
	}
	if result5 != 150 {
		t.Errorf("expected 150, got %f", result5)
	}
	

}
