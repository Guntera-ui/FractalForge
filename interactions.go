package main

import (
	"fmt"
	"os"

	"golang.org/x/term"
)

func RunInteractive(view Viewport, width, height int, cellAspect float64, maxIter int) error {
	fd := int(os.Stdin.Fd())

	if !term.IsTerminal(fd) || !term.IsTerminal(int(os.Stdout.Fd())) {
		return fmt.Errorf("interactive mode requires a terminal")
	}

	oldState, err := term.MakeRaw(fd)
	if err != nil {
		return err
	}
	defer term.Restore(fd, oldState)

	fmt.Print("\x1b[?1049h\x1b[?25l")
	defer fmt.Print("\x1b[0m\x1b[?25h\x1b[?1049l")

	var key [1]byte

	for {
		fmt.Print("\x1b[H\x1b[2J")

		RenderGrid(width, height, cellAspect, view, maxIter)

		fmt.Print("\r\nWASD: Move | +/-: Zoom | Q: Quit")

		_, err := os.Stdin.Read(key[:])
		if err != nil {
			return err
		}

		if key[0] == 3 {
			return nil
		}

		newView, running := HandleKey(view, key[0])

		if !running {
			return nil
		}

		view = newView
	}
}
