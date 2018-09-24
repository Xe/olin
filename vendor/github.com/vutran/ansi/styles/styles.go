/*
Apply custom font treatment to your strings
*/
package styles

import (
	"github.com/vutran/ansi"
)

func Bold(value string) string {
	return ansi.Bold() + value + ansi.BoldOff()
}

func Faint(value string) string {
	return ansi.Faint() + value + ansi.FaintOff()
}

func Italic(value string) string {
	return ansi.Italic() + value + ansi.ItalicOff()
}

func Underline(value string) string {
	return ansi.Underline() + value + ansi.UnderlineOff()
}

func CrossedOut(value string) string {
	return ansi.CrossedOut() + value + ansi.CrossedOutOff()
}
