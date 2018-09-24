package ansi

import (
	"strconv"
)

// HideCursor hides the cursor
func HideCursor() string {
	return Esc + "[?25l"
}

// ShowCursor shows the cursor
func ShowCursor() string {
	return Esc + "[?25h"
}

// CursorUp moves the cursor up by `n` lines
func CursorUp(n int) string {
	return Esc + "[" + strconv.Itoa(n) + "A"
}

// CursorDown moves the cursor down by `n` lines
func CursorDown(n int) string {
	return Esc + "[" + strconv.Itoa(n) + "B"
}

// CursorForward moves the cursor forward by `n` columns
func CursorForward(n int) string {
	return Esc + "[" + strconv.Itoa(n) + "C"
}

// CursorBackward moves the cursor backwards by `n` columns
func CursorBackward(n int) string {
	return Esc + "[" + strconv.Itoa(n) + "D"
}

// CursorStart moves the cursor to column `n`
func CursorStart(n int) string {
	return Esc + "[" + strconv.Itoa(n) + "G"
}

// SaveCursorPosition saves the current cursor position
func SaveCursorPosition() string {
	return Esc + "[s"
}

// RestoreCursorPosition restores the cursor position
func RestoreCursorPosition() string {
	return Esc + "[u"
}

// GetCursorPosition retrieves the cursor position
func GetCursorPosition() string {
	return Esc + "[6n"
}
