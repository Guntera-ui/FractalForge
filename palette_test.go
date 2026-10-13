package main

import "testing"

func TestLerpColor(t *testing.T) {
	tests := []struct {
		name     string
		start    RGB
		end      RGB
		factor   float64
		expected RGB
	}{
		{
			name:     "Start color",
			start:    RGB{50, 100, 150},
			end:      RGB{150, 200, 250},
			factor:   0,
			expected: RGB{50, 100, 150},
		},
		{
			name:     "Midpoint",
			start:    RGB{50, 100, 150},
			end:      RGB{150, 200, 250},
			factor:   0.5,
			expected: RGB{100, 150, 200},
		},
		{
			name:     "End color",
			start:    RGB{50, 100, 150},
			end:      RGB{150, 200, 250},
			factor:   1,
			expected: RGB{150, 200, 250},
		},
		{
			name:     "Reverse interpolation",
			start:    RGB{200, 150, 100},
			end:      RGB{100, 50, 0},
			factor:   0.5,
			expected: RGB{150, 100, 50},
		},
		{
			name:     "Rounding fractional values",
			start:    RGB{0, 0, 0},
			end:      RGB{101, 51, 201},
			factor:   0.5,
			expected: RGB{51, 26, 101},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := LerpColor(tt.start, tt.end, tt.factor)

			if result != tt.expected {
				t.Errorf("expected %+v, got %+v", tt.expected, result)
			}
		})
	}
}
