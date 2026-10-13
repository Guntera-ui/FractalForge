# FractalForge

A terminal-based Mandelbrot set renderer written in Go. It uses Unicode half-block characters, 24-bit RGB colors, and smooth escape-time coloring.

![FractalForge rendering the Mandelbrot set](assets/fractalforge-preview.png)

## Features

- Mandelbrot set rendering
- Smooth escape-time coloring
- Classic purple-and-gold palette
- Two vertically stacked samples per terminal cell using Unicode half-blocks
- Configurable canvas dimensions, viewport, and maximum iteration count

## Requirements

- Go
- A terminal supporting Unicode and 24-bit RGB colors (such as Kitty)

## Configuration

The current version uses variables in the source code rather than command-line options. Set these before running.

**In `main.go`:**

`RenderGrid(360, 85, 2.0)`

| Argument | Meaning |
|---|---|
| `360` | Canvas width in terminal columns |
| `85` | Canvas height in terminal rows |
| `2.0` | Approximate terminal character height-to-width ratio (`cellAspect`), used to correct the image's proportions |

Check your terminal's available dimensions with `stty size` (rows first, then columns). Make sure the canvas fits your terminal.

**In `renderer.go`, inside `RenderGrid()`:**

`minReal := -3.5`

`maxReal := 1.7`

`maxIter := 1000000`

| Variable | Meaning |
|---|---|
| `minReal` | Left edge of the complex-plane viewport |
| `maxReal` | Right edge of the complex-plane viewport |
| `maxIter` | Maximum iterations per sample; larger values can considerably increase rendering time |

The imaginary-axis range is calculated automatically from the canvas dimensions and `cellAspect`, keeping it centered at zero.

To customize the colors, edit the `colors` slice inside `ColorFromClassicPalette()` in `pallete.go`.

## Running

Clone the repository:

`git clone https://github.com/Guntera-ui/fractalforge.git`

`cd fractalforge`

Run the program:

`go run .`

If you're using Nix, enter the development shell first:

`nix develop`

## Mathematics

### The Mandelbrot Set

The Mandelbrot set is defined by the recurrence:

\[
z_{n+1}=z_n^2+c,\qquad z_0=0
\]

Here, `c` is a fixed complex number. A point belongs to the Mandelbrot set if its resulting sequence remains bounded indefinitely.

### Escape-Time Algorithm

A sequence is guaranteed to diverge if:

\[
|z_n|>2
\]

FractalForge checks the equivalent squared-magnitude condition:

\[
\operatorname{Re}(z_n)^2+\operatorname{Im}(z_n)^2>4
\]

This avoids calculating square roots.

Points that do not escape within `maxIter` are colored black. This does not necessarily prove that they belong to the Mandelbrot set.

### Smooth Coloring

To reduce color banding, escaped points use a normalized iteration count:

\[
\nu=n+1-\log_2(\log_2|z_n|)
\]

This produces fractional escape values, which are mapped onto the color palette using linear RGB interpolation.

The colors are a visualization choice, not an intrinsic property of the Mandelbrot set.

## Tests

`go test ./...`

## Status

Work in progress. Planned improvements include zoom controls, additional palettes, other fractals, and higher-resolution rendering.
