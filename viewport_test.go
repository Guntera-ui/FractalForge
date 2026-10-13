package main

import (
	"math"
	"testing"
)

func TestRealBounds(t *testing.T) {
	minReal, maxReal := RealBounds(-0.75, 3.0)

	if minReal != -2.25 {
		t.Errorf("Expected minReal -2.25, got %f", minReal)
	}

	if maxReal != 0.75 {
		t.Errorf("Expected maxReal 0.75, got %f", maxReal)
	}
}

func TestZoom(t *testing.T) {
	original := Viewport{-0.75, 0.125, 4.0}

	tests := []struct {
		name          string
		factor        float64
		expectedWidth float64
	}{
		{"Zoom in 2x", 2.0, 2.0},
		{"Zoom in 4x", 4.0, 1.0},
		{"Zoom out 2x", 0.5, 8.0},
		{"No zoom", 1.0, 4.0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := Zoom(original, tt.factor)

			if math.Abs(result.Width-tt.expectedWidth) > 1e-9 {
				t.Errorf(
					"Expected width %f, got %f",
					tt.expectedWidth,
					result.Width,
				)
			}

			if result.CenterReal != original.CenterReal ||
				result.CenterImag != original.CenterImag {
				t.Error("Zoom should not change the viewport center")
			}
		})
	}

	if original.Width != 4.0 {
		t.Error("Zoom modified the original viewport")
	}
}

func TestPan(t *testing.T) {
	original := Viewport{-0.75, 0.0, 3.0}

	tests := []struct {
		name         string
		deltaReal    float64
		deltaImag    float64
		expectedReal float64
		expectedImag float64
	}{
		{"Move right", 0.25, 0.0, -0.5, 0.0},
		{"Move left", -0.25, 0.0, -1.0, 0.0},
		{"Move up", 0.0, 0.25, -0.75, 0.25},
		{"Move down", 0.0, -0.25, -0.75, -0.25},
		{"Move diagonally", 0.25, 0.5, -0.5, 0.5},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := Pan(original, tt.deltaReal, tt.deltaImag)

			if result.CenterReal != tt.expectedReal {
				t.Errorf(
					"Expected real center %f, got %f",
					tt.expectedReal,
					result.CenterReal,
				)
			}

			if result.CenterImag != tt.expectedImag {
				t.Errorf(
					"Expected imaginary center %f, got %f",
					tt.expectedImag,
					result.CenterImag,
				)
			}

			if result.Width != original.Width {
				t.Error("Pan should not change the viewport width")
			}
		})
	}

	if original != (Viewport{-0.75, 0.0, 3.0}) {
		t.Error("Pan modified the original viewport")
	}
}
