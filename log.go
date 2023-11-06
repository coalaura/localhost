package main

import (
	"os"

	"github.com/gookit/color"
)

func InfoGreen(msg, green string) {
	green = color.RenderCode("38;5;34", green)

	color.Fprint(os.Stdout, " - "+msg)
	color.Fprintln(os.Stdout, green)

	_, _ = color.Reset()
}

func InfoPlain(msg, plain string) {
	plain = color.RenderCode("38;5;111", plain)

	color.Fprint(os.Stdout, " - "+msg)
	color.Fprintln(os.Stdout, plain)

	_, _ = color.Reset()
}

func InfoRed(msg, red string) {
	red = color.RenderCode("38;5;124", red)

	color.Fprint(os.Stdout, " - "+msg)
	color.Fprintln(os.Stdout, red)

	_, _ = color.Reset()
}
