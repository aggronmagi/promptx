// Code generated by "gogen import"; DO NOT EDIT.
// Exec: "gogen import ./output -t Color -t ConsoleWriter -o gen_output.go"
// Version: 0.0.2

package promptx

import output "github.com/aggronmagi/promptx/output"

// Color represents color on terminal.
type Color = output.Color

const (
	// Black represents a black.
	Black = output.Black
	// Blue represents a blue.
	Blue = output.Blue
	// Brown represents a brown.
	Brown = output.Brown
	// Cyan represents a cyan.
	Cyan = output.Cyan
	// DarkBlue represents a dark blue.
	DarkBlue = output.DarkBlue
	// DarkGray represents a dark gray.
	DarkGray = output.DarkGray
	// DarkGreen represents a dark green.
	DarkGreen = output.DarkGreen
	// DarkRed represents a dark red.
	DarkRed = output.DarkRed
	// DefaultColor represents a default color.
	DefaultColor = output.DefaultColor
	// Fuchsia represents a fuchsia.
	Fuchsia = output.Fuchsia
	// Green represents a green.
	Green = output.Green
	// LightGray represents a light gray.
	LightGray = output.LightGray
	// Purple represents a purple.
	Purple = output.Purple
	// Red represents a red.
	Red = output.Red
	// Turquoise represents a turquoise.
	Turquoise = output.Turquoise
	// White represents a white.
	White = output.White
	// Yellow represents a yellow.
	Yellow = output.Yellow
)

// ConsoleWriter is an interface to abstract output layer.
type ConsoleWriter = output.ConsoleWriter
