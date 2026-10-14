package main

import "testing"

func TestGetPreset(t *testing.T) {
	tests := []struct {
		name     string
		preset   string
		expected Viewport
		exists   bool
	}{
		{
			name:     "Full Mandelbrot",
			preset:   "full",
			expected: Viewport{-0.9, 0.0, 5.2},
			exists:   true,
		},
		{
			name:     "Seahorse Valley",
			preset:   "seahorse",
			expected: Viewport{-0.7435, 0.1314, 0.03},
			exists:   true,
		},
		{
			name:     "Unknown preset",
			preset:   "unknown",
			expected: Viewport{},
			exists:   false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			view, exists := GetPreset(tt.preset)

			if exists != tt.exists {
				t.Errorf("expected exists=%t, got %t", tt.exists, exists)
			}

			if view != tt.expected {
				t.Errorf("expected %+v, got %+v", tt.expected, view)
			}
		})
	}
}
