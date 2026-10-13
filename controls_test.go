package main

import (
	"math"
	"testing"
)

func TestHandleKey(t *testing.T) {
	original := Viewport{-0.75, 0.125, 4.0}

	tests := []struct {
		name      string
		key       byte
		expected  Viewport
		shouldRun bool
	}{
		{
			name:      "Move up",
			key:       'w',
			expected:  Viewport{-0.75, 0.525, 4.0},
			shouldRun: true,
		},
		{
			name:      "Move down",
			key:       's',
			expected:  Viewport{-0.75, -0.275, 4.0},
			shouldRun: true,
		},
		{
			name:      "Move left",
			key:       'a',
			expected:  Viewport{-1.15, 0.125, 4.0},
			shouldRun: true,
		},
		{
			name:      "Move right",
			key:       'd',
			expected:  Viewport{-0.35, 0.125, 4.0},
			shouldRun: true,
		},
		{
			name:      "Zoom in",
			key:       '+',
			expected:  Viewport{-0.75, 0.125, 2.0},
			shouldRun: true,
		},
		{
			name:      "Zoom out",
			key:       '-',
			expected:  Viewport{-0.75, 0.125, 8.0},
			shouldRun: true,
		},
		{
			name:      "Quit",
			key:       'q',
			expected:  original,
			shouldRun: false,
		},
		{
			name:      "Unknown key",
			key:       'z',
			expected:  original,
			shouldRun: true,
		},
	}

	const tolerance = 1e-9

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, running := HandleKey(original, tt.key)

			if math.Abs(result.CenterReal-tt.expected.CenterReal) > tolerance {
				t.Errorf(
					"Expected CenterReal %f, got %f",
					tt.expected.CenterReal,
					result.CenterReal,
				)
			}

			if math.Abs(result.CenterImag-tt.expected.CenterImag) > tolerance {
				t.Errorf(
					"Expected CenterImag %f, got %f",
					tt.expected.CenterImag,
					result.CenterImag,
				)
			}

			if math.Abs(result.Width-tt.expected.Width) > tolerance {
				t.Errorf(
					"Expected Width %f, got %f",
					tt.expected.Width,
					result.Width,
				)
			}

			if running != tt.shouldRun {
				t.Errorf(
					"Expected running=%t, got %t",
					tt.shouldRun,
					running,
				)
			}
		})
	}

	if original != (Viewport{-0.75, 0.125, 4.0}) {
		t.Error("HandleKey modified the original viewport")
	}
}
