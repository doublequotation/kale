package utils

import (
	"github.com/muesli/termenv"
)

type Colors struct {
	Red     termenv.Color
	Green   termenv.Color
	Yellow  termenv.Color
	Blue    termenv.Color
	Magenta termenv.Color
	Cyan    termenv.Color
	Gray    termenv.Color
}

func InitColors() Colors {
	var colors Colors
	p := termenv.ColorProfile()
	colors.Red = p.Color("#E88388")
	colors.Green = p.Color("#A8CC8C")
	colors.Yellow = p.Color("#DBAB79")
	colors.Blue = p.Color("#71BEF2")
	colors.Magenta = p.Color("#D290E4")
	colors.Cyan = p.Color("#66C2CD")
	colors.Gray = p.Color("#B9BFCA")
	return colors
}
func DelLine() {
	termenv.CursorPrevLine(1)
	termenv.ClearLine()
}
