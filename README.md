# FractalForge

A terminal-based Mandelbrot set renderer written in Go.

![FractalForge Preview](assets/fractalforge-preview.png)

FractalForge renders the Mandelbrot set directly in your terminal using Unicode half-block characters, 24-bit RGB colors, and smooth escape-time coloring.

## Features

- Mandelbrot set rendering
- Smooth escape-time coloring
- Classic purple-and-gold color palette
- Unicode half-block rendering for increased resolution
- Configurable canvas dimensions and iteration limit
- Automated tests

## Requirements

- Go
- A terminal supporting Unicode and true-color (24-bit RGB), such as Kitty

## Installation

Clone the repository:

```bash
git clone https://github.com/Guntera-ui/fractalforge.git
cd fractalforge
```

Run:

```bash
go run .
```

For Nix users, the project includes a development shell:

```bash
nix develop
go run .
```

## Configuration

FractalForge currently uses source-code configuration.

### Canvas Size

In `main.go`:

```go
RenderGrid(360, 85, 2.0)
```

| Parameter | Description |
|---|---|
| `360` | Width in terminal columns |
| `85` | Height in terminal rows |
| `2.0` | Character height-to-width ratio for aspect correction |

Check your terminal dimensions before rendering:

```bash
stty size
```

The output shows rows followed by columns.

### Mandelbrot Settings

In `renderer.go`:

```go
minReal := -3.5
maxReal := 1.7
maxIter := 1000000
```

| Variable | Description |
|---|---|
| `minReal` | Left boundary of the complex-plane viewport |
| `maxReal` | Right boundary of the complex-plane viewport |
| `maxIter` | Maximum number of iterations per point |

The imaginary-axis boundaries are calculated automatically based on canvas dimensions and aspect ratio.

Higher iteration limits can significantly increase rendering time.

### Colors

The color gradient is defined inside `ColorFromClassicPalette()` in `pallete.go`.

Modify the RGB values to customize the palette.

## Mathematics

### Mandelbrot Set

The Mandelbrot set is defined by the recurrence relation:

$$
z_{n+1}=z_n^2+c
$$

with the initial condition:

$$
z_0=0
$$

Here, $c$ is a fixed complex number and $z_n$ represents the value after $n$ iterations.

A point belongs to the Mandelbrot set if its sequence remains bounded indefinitely.

### Escape-Time Algorithm

A sequence is guaranteed to diverge when:

$$
|z_n|>2
$$

FractalForge uses the equivalent squared-magnitude condition:

$$
\mathrm{Re}(z_n)^2 + \mathrm{Im}(z_n)^2 > 4
$$

This avoids unnecessary square-root calculations.

Points that do not escape within the iteration limit are rendered black. This does not necessarily prove that they belong to the Mandelbrot set.

### Smooth Escape-Time Coloring

Integer iteration counts produce distinct color bands.

To produce smoother transitions, FractalForge uses the normalized iteration count:

$$
\nu=n+1-\log_2(\log_2|z_n|)
$$

Where:

- $n$ is the escape iteration count.
- $|z_n|$ is the final magnitude.
- $\nu$ is the smooth iteration value.

The fractional value is mapped onto the color palette using linear RGB interpolation.

Colors are used for visualization and are not intrinsic to the Mandelbrot set.

### Terminal Rendering

FractalForge uses the Unicode upper-half-block character `▀`.

- Foreground color represents the upper sample.
- Background color represents the lower sample.

Each terminal cell therefore displays two vertically stacked fractal samples.

## Tests

Run:

```bash
go test ./...
```

## Roadmap

- [x] Mandelbrot set rendering
- [x] Smooth escape-time coloring
- [x] RGB color interpolation
- [x] Unicode half-block renderer
- [ ] Interior detection and performance optimizations
- [ ] Automatic terminal sizing
- [ ] Zoom and pan controls
- [ ] Customizable palettes
- [ ] Julia sets
- [ ] Kitty graphics-protocol rendering

## Status

FractalForge is a work in progress, built as a Go programming and mathematical visualization project.
