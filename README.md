# FractalForge

**An interactive Mandelbrot set explorer for the terminal, written in Go.**

FractalForge renders the Mandelbrot set using smooth escape-time coloring, customizable palettes, and high-resolution terminal graphics.

![FractalForge Preview](assets/fractalforge-preview.png)

## Features

- Interactive zooming and panning
- High-resolution Kitty Graphics Protocol rendering
- Unicode half-block fallback
- Multicore, tile-based rendering
- Automatic progressive refinement
- Smooth escape-time coloring
- 8 built-in color palettes
- Continuous and cyclic coloring modes
- 6 Mandelbrot location presets
- Automatic terminal sizing

## Installation

Requires Go and a terminal supporting ANSI true colors. Kitty graphics mode requires a compatible terminal such as Kitty or Ghostty.

```bash
git clone https://github.com/Guntera-ui/fractalforge.git
cd fractalforge
go build -o fractalforge .
./fractalforge
```

To install globally for your user:

```bash
mkdir -p ~/.local/bin
go build -o ~/.local/bin/fractalforge .
```

Make sure `~/.local/bin` is in your `PATH`.

Nix users can enter the development environment with:

```bash
nix develop
go run .
```

## Usage

### Unicode rendering

The default backend works in terminals supporting Unicode and 24-bit colors.

```bash
fractalforge
fractalforge -interactive
```

### Kitty graphics

Use high-resolution raster graphics in Kitty, Ghostty, or another compatible terminal:

```bash
fractalforge -renderer kitty
fractalforge -renderer kitty -interactive
```

Adjust resolution:

```bash
fractalforge -renderer kitty -scale 8 -interactive
```

Higher scale values produce finer details but require more computation.

### Progressive rendering

Interactive Kitty graphics use an automatic two-stage rendering process:

1. **Preview:** A lower-resolution image appears progressively as tiles are calculated.
2. **Refinement:** When navigation stops, the application automatically renders a full-quality image.

FractalForge uses Go goroutines to calculate tiles concurrently across available CPU cores. Obsolete rendering jobs are cancelled when the viewport changes.

Completed full-quality tiles are combined into a seamless final image.

No manual refresh is required.

### Interactive controls

| Key | Action |
|---|---|
| `W` / `S` | Pan up / down |
| `A` / `D` | Pan left / right |
| `+` / `-` | Zoom in / out |
| `Q` / `Ctrl+C` | Quit |

### Customization

```bash
# Change palette
fractalforge -palette catppuccin

# Explore Seahorse Valley
fractalforge -preset seahorse -renderer kitty -interactive

# Use cyclic coloring
fractalforge -coloring cyclic -color-period 32

# Increase iteration limit
fractalforge -renderer kitty -interactive -iter 2000
```

**Palettes:** `classic`, `catppuccin`, `tokyonight`, `gruvbox`, `dracula`, `nord`, `rosepine`, `kanagawa`

**Presets:** `full`, `seahorse`, `elephant`, `triple-spiral`, `minibrot`, `scepter`

Other options include `-zoom`, `-x`, `-y`, `-width`, `-height`, and `-aspect`.

Use `fractalforge -help` for the complete list.

## Mathematics

FractalForge evaluates the Mandelbrot recurrence:

$$
z_{n+1}=z_n^2+c,\qquad z_0=0
$$

A sequence is known to escape when:

$$
|z_n|>2
$$

To produce smoother color gradients, the renderer calculates a fractional escape-time value:

$$
\nu=n+1-\log_2(\log_2|z_n|)
$$

This value is mapped to RGB colors using interpolation.

**Continuous mode** uses an aggressive logarithmic gradient that reaches its final color around 128 smooth iterations.

**Cyclic mode** repeats the selected palette using a configurable period, blending smoothly between its last and first colors.

Points that do not escape within the iteration limit receive the palette's interior color.

## Development

```bash
go fmt ./...
go test ./...
go vet ./...
go build .
```

## Roadmap

- [x] Interactive Mandelbrot explorer
- [x] Smooth coloring and multiple palettes
- [x] High-resolution Kitty graphics
- [x] Multicore tile rendering
- [x] Progressive rendering and cancellation
- [ ] Julia sets
- [ ] Burning Ship and Newton fractals
- [ ] Advanced mathematical coloring
- [ ] Save and load custom locations
- [ ] Image export

## License

See [LICENSE](LICENSE).
