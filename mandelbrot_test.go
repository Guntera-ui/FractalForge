package main

import "testing"
import "math"

func TestMandelbrot(t *testing.T) {
	result := Mandelbrot(3+0i, 100)
	result2 := Mandelbrot(1+0i, 100)
	result3 := Mandelbrot(0+0i, 100)
	result4 := Mandelbrot(-1+0i, 100)
	result5 := Mandelbrot(2+0i, 100)
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

func TestSmoothIteration(t *testing.T) {
	tests := []struct {
		name     string
		input    EscapeResult
		expected float64
	}{
		{
			name: "Exact smooth iteration",
			input: EscapeResult{
				Iterations: 8,
				FinalZ:     4 + 0i,
				Escaped:    true,
			},
			expected: 8.0,
		},
		{
			name: "Fractional smooth iteration",
			input: EscapeResult{
				Iterations: 8,
				FinalZ:     3 + 0i,
				Escaped:    true,
			},
			expected: 8.33555129,
		},
		{
			name: "Sequence did not escape",
			input: EscapeResult{
				Iterations: 100,
				FinalZ:     0 + 0i,
				Escaped:    false,
			},
			expected: 100.0,
		},
		{
			name: "Escape at iteration limit",
			input: EscapeResult{
				Iterations: 3,
				FinalZ:     4 + 0i,
				Escaped:    true,
			},
			expected: 3.0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := SmoothIteration(tt.input)

			if math.IsNaN(result) ||
				math.Abs(result-tt.expected) > 0.000001 {
				t.Errorf(
					"expected %.8f, got %.8f",
					tt.expected,
					result,
				)
			}
		})
	}
}
