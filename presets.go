package main

import (
	"fmt"
	"sort"
)

type Preset struct {
	View        Viewport
	Description string
}

var presets = map[string]Preset{
	"full": {
		View:        Viewport{-0.9, 0, 5.2},
		Description: "The complete Mandelbrot set",
	},
	"seahorse": {
		View:        Viewport{-0.7435, 0.1314, 0.03},
		Description: "Seahorse Valley's curling filaments",
	},
	"elephant": {
		View:        Viewport{0.27205, 0.006118, 0.025},
		Description: "Elephant Valley's repeating bulbs",
	},
	"triple-spiral": {
		View:        Viewport{-0.0888, 0.6547, 0.015},
		Description: "Interwoven spiral structures",
	},
	"minibrot": {
		View:        Viewport{-1.7685, 0.0005, 0.03},
		Description: "Miniature Mandelbrot near the western antenna",
	},
	"scepter": {
		View:        Viewport{-1.36, 0.005, 0.1},
		Description: "Branching structures on the western side",
	},
}

func GetPreset(name string) (Viewport, bool) {
	preset, exists := presets[name]
	return preset.View, exists
}

func PrintPresets() {
	names := make([]string, 0, len(presets))

	for name := range presets {
		names = append(names, name)
	}

	sort.Strings(names)

	fmt.Println("Available FractalForge presets:")

	for _, name := range names {
		preset := presets[name]

		fmt.Printf("  %-14s %s\n", name, preset.Description)

		fmt.Printf(
			"  %-14s center: %.8f %+.8fi | width: %g\n",
			"",
			preset.View.CenterReal,
			preset.View.CenterImag,
			preset.View.Width,
		)
	}
}
