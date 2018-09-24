package ansi

import (
	"strconv"
)

// EraseDisplay clears the screen
func EraseDisplay(code int) string {
	return Esc + "[" + strconv.Itoa(code) + "2J"
}

// EraseLine clears part of the line.
// If `n` is zero, clears the cursor to the end of the line.
// If `n` is one, clear from cursor to beginning of the line.
// If `n` is two, clear entire line.
// Cursor position does not change
func EraseLine(code int) string {
	return Esc + "[" + strconv.Itoa(code) + "K"
}

// SelectGraphicsRendition sets the SGR parameters.
func SelectGraphicsRendition(code int) string {
	return Esc + "[" + strconv.Itoa(code) + "m"
}
