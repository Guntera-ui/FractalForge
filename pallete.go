package main

import (
	"fmt"
	"math"
	"sort"
)

type RGB struct {
	R int
	G int
	B int
}

func Lerp(start, end, t float64) float64 {
	// R = start + t(end - start)
	interpolated := start + (t * (end - start))
	return interpolated
}

func LerpColor(start, end RGB, t float64) RGB {
	R := int(math.Round(Lerp(float64(start.R), float64(end.R), t)))
	G := int(math.Round(Lerp(float64(start.G), float64(end.G), t)))
	B := int(math.Round(Lerp(float64(start.B), float64(end.B), t)))

	result := RGB{R, G, B}
	return result
}

func ColorFromIterations(
	iterations, maxIter int,
	startColor, endColor RGB,
) RGB {
	if iterations == maxIter {
		return RGB{4, 8, 20}
	}

	t := math.Log1p(float64(iterations)) /
		math.Log1p(float64(maxIter))

	return LerpColor(startColor, endColor, t)
}

type ColoringMode string

const (
	ColoringContinuous ColoringMode = "continuous"
	ColoringCyclic     ColoringMode = "cyclic"
)

type Palette struct {
	Colors   []RGB
	Interior RGB
	Mode     ColoringMode
	Period   float64
}

func ColorFromClassicPalette(result EscapeResult, maxIter ...int) RGB {
	limit := 10000

	if len(maxIter) > 0 {
		limit = maxIter[0]
	}

	return ColorFromPalette(result, palettes["classic"], limit)
}

func ColorFromPalette(result EscapeResult, palette Palette, maxIter int) RGB {
	if !result.Escaped {
		return palette.Interior
	}

	if len(palette.Colors) == 0 {
		return palette.Interior
	}

	if len(palette.Colors) == 1 {
		return palette.Colors[0]
	}

	iterations := math.Max(0, SmoothIteration(result))
	colors := palette.Colors

	var t float64

	switch palette.Mode {
	case ColoringCyclic:
		period := palette.Period

		if period <= 0 || math.IsNaN(period) || math.IsInf(period, 0) {
			period = 32
		}

		t = math.Mod(iterations, period) / period

		position := t * float64(len(colors))
		index := int(position) % len(colors)
		next := (index + 1) % len(colors)

		fraction := position - math.Floor(position)

		return LerpColor(
			colors[index],
			colors[next],
			fraction,
		)

	default:
		_ = maxIter

		t = math.Log2(math.Max(1, iterations)) / 7.0
		t = math.Max(0, math.Min(1, t))
	}

	position := t * float64(len(colors)-1)
	index := int(position)

	if index >= len(colors)-1 {
		return colors[len(colors)-1]
	}

	fraction := position - float64(index)

	return LerpColor(
		colors[index],
		colors[index+1],
		fraction,
	)
}

func GetPalette(name string) (Palette, bool) {
	palette, exists := palettes[name]
	return palette, exists
}

func PrintPalettes() {
	names := make([]string, 0, len(palettes))

	for name := range palettes {
		names = append(names, name)
	}

	sort.Strings(names)

	fmt.Println("Available FractalForge palettes:")

	for _, name := range names {
		palette := palettes[name]

		fmt.Printf(
			"  %-14s %d colors\n",
			name,
			len(palette.Colors),
		)
	}
}

var palettes = map[string]Palette{
	"classic": {
		Colors: []RGB{
			{94, 68, 239},
			{101, 52, 222},
			{109, 35, 206},
			{117, 19, 190},
			{126, 4, 175},
			{134, 10, 158},
			{141, 27, 141},
			{159, 58, 111},
			{171, 88, 75},
			{198, 145, 58},
			{222, 193, 41},
			{255, 218, 119},
			{255, 241, 183},
			{182, 207, 255},
		},
		Interior: RGB{0, 0, 0},
	},

	"catppuccin": {
		Colors: []RGB{
			{137, 180, 250}, // Blue
			{116, 199, 236}, // Sapphire
			{137, 220, 235}, // Sky
			{148, 226, 213}, // Teal
			{166, 227, 161}, // Green
			{249, 226, 175}, // Yellow
			{250, 179, 135}, // Peach
			{235, 160, 172}, // Maroon
			{243, 139, 168}, // Red
			{245, 194, 231}, // Pink
			{203, 166, 247}, // Mauve
			{180, 190, 254}, // Lavender
			{242, 205, 205}, // Flamingo
			{245, 224, 220}, // Rosewater
		},
		Interior: RGB{30, 30, 46},
	},

	"tokyonight": {
		Colors: []RGB{
			{122, 162, 247}, // Blue
			{125, 207, 255}, // Cyan
			{42, 195, 222},  // Blue1
			{26, 188, 156},  // Teal
			{158, 206, 106}, // Green
			{224, 175, 104}, // Yellow
			{255, 158, 100}, // Orange
			{247, 118, 142}, // Red
			{187, 154, 247}, // Magenta
			{157, 124, 216}, // Purple
			{192, 202, 245}, // Foreground
		},
		Interior: RGB{26, 27, 38},
	},

	"gruvbox": {
		Colors: []RGB{
			{69, 133, 136},  // Neutral Blue
			{131, 165, 152}, // Bright Blue
			{142, 192, 124}, // Bright Aqua
			{184, 187, 38},  // Bright Green
			{250, 189, 47},  // Bright Yellow
			{215, 153, 33},  // Neutral Yellow
			{254, 128, 25},  // Bright Orange
			{251, 73, 52},   // Bright Red
			{204, 36, 29},   // Neutral Red
			{211, 134, 155}, // Bright Purple
			{177, 98, 134},  // Neutral Purple
			{104, 157, 106}, // Neutral Aqua
		},
		Interior: RGB{40, 40, 40},
	},

	"dracula": {
		Colors: []RGB{
			{68, 71, 90},    // Current Line
			{98, 114, 164},  // Comment
			{189, 147, 249}, // Purple
			{255, 121, 198}, // Pink
			{255, 85, 85},   // Red
			{255, 184, 108}, // Orange
			{241, 250, 140}, // Yellow
			{80, 250, 123},  // Green
			{139, 233, 253}, // Cyan
			{248, 248, 242}, // Foreground
		},
		Interior: RGB{40, 42, 54},
	},

	"nord": {
		Colors: []RGB{
			{94, 129, 172},  // Nord10
			{129, 161, 193}, // Nord9
			{136, 192, 208}, // Nord8
			{143, 188, 187}, // Nord7
			{163, 190, 140}, // Nord14
			{235, 203, 139}, // Nord13
			{208, 135, 112}, // Nord12
			{191, 97, 106},  // Nord11
			{180, 142, 173}, // Nord15
			{216, 222, 233}, // Nord4
			{229, 233, 240}, // Nord5
			{236, 239, 244}, // Nord6
		},
		Interior: RGB{46, 52, 64},
	},

	"rosepine": {
		Colors: []RGB{
			{82, 79, 103},   // Highlight High
			{110, 106, 134}, // Muted
			{144, 140, 170}, // Subtle
			{49, 116, 143},  // Pine
			{156, 207, 216}, // Foam
			{149, 177, 172}, // Leaf
			{196, 167, 231}, // Iris
			{235, 188, 186}, // Rose
			{235, 111, 146}, // Love
			{246, 193, 119}, // Gold
			{224, 222, 244}, // Text
		},
		Interior: RGB{25, 23, 36},
	},

	"kanagawa": {
		Colors: []RGB{
			{45, 79, 103},   // Wave Blue 2
			{101, 133, 148}, // Dragon Blue
			{126, 156, 216}, // Crystal Blue
			{127, 180, 202}, // Spring Blue
			{122, 168, 159}, // Wave Aqua 2
			{152, 187, 108}, // Spring Green
			{230, 195, 132}, // Carp Yellow
			{255, 160, 102}, // Surimi Orange
			{228, 104, 118}, // Wave Red
			{210, 126, 153}, // Sakura Pink
			{149, 127, 184}, // Oni Violet
			{220, 215, 186}, // Fuji White
		},
		Interior: RGB{31, 31, 40},
	},
}
