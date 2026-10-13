# FractalForge

A terminal-based Mandelbrot set renderer written in Go.

FractalForge uses Unicode half-block characters, 24-bit RGB colors, and smooth escape-time coloring to display fractals directly in the terminal.

## Features

- Mandelbrot set rendering
- Smooth escape-time coloring
- Classic purple-and-gold palette
- Unicode half-block rendering for increased resolution
- Adjustable canvas size and iteration limit

## Requirements

- Go
- A terminal supporting Unicode and 24-bit colors (such as Kitty)

## Running

Clone the repository:

```bash
git clone https://github.com/Guntera-ui/fractalforge.git
cd fractalforge
go run .
```

If you're using Nix, enter the development environment first:

```bash
nix develop
go run .
```

## Mathematics

### The Mandelbrot Set

The Mandelbrot set is defined by the recurrence relation:

**zₙ₊₁ = zₙ² + c**

Where:

- `z₀ = 0` is the initial value.
- `c` is a fixed complex number representing a point in the complex plane.
- `zₙ` is the result after `n` iterations.

For each point `c`, we repeatedly apply the equation and observe whether the sequence remains bounded.

A point belongs to the Mandelbrot set if its sequence remains bounded indefinitely.

### Escape-Time Algorithm

For the Mandelbrot set, a sequence is guaranteed to diverge if its magnitude exceeds 2.

**|zₙ| > 2**

Instead of calculating the square root required for the magnitude, FractalForge checks:

**Re(zₙ)² + Im(zₙ)² > 4**

This is mathematically equivalent and avoids an unnecessary square-root operation.

Since we cannot run infinitely many iterations, the renderer uses a maximum iteration limit (`maxIter`).

- If a point escapes, we record its iteration count and final complex value.
- If it does not escape within the limit, we render it black.

Not escaping within the limit does not necessarily prove that a point belongs to the Mandelbrot set.

### Smooth Escape-Time Coloring

Using only integer escape counts produces visible bands of color.

To reduce this effect, FractalForge calculates a normalized iteration count:

**ν = n + 1 − log₂(log₂|zₙ|)**

Where:

- `n` is the number of iterations before escape.
- `|zₙ|` is the magnitude of the final complex value.
- `ν` is the smooth iteration value.

Unlike ordinary escape counts, this produces fractional values, allowing smoother transitions between neighboring colors.

The resulting value is mapped onto a color palette using linear interpolation.

The colors themselves are a visualization choice; they are not intrinsic to the Mandelbrot set.

### Terminal Rendering

Each terminal character cell uses the Unicode upper-half-block character (`▀`).

The foreground color represents the upper pixel, while the background color represents the lower pixel.

This allows FractalForge to display two vertically stacked samples per terminal character cell.

## Tests

```bash
go test ./...
```

## Status

Work in progress.

Planned features include zoom controls, additional fractals, customizable palettes, and higher-resolution rendering.
