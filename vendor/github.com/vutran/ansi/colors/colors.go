/*
Apply custom color to your string
*/
package colors

import (
	"github.com/vutran/ansi"
)

func Black(value string) string {
	return ansi.Black() + value + ansi.DefaultColor()
}

func Red(value string) string {
	return ansi.Red() + value + ansi.DefaultColor()
}

func Green(value string) string {
	return ansi.Green() + value + ansi.DefaultColor()
}

func Yellow(value string) string {
	return ansi.Yellow() + value + ansi.DefaultColor()
}

func Blue(value string) string {
	return ansi.Blue() + value + ansi.DefaultColor()
}

func Magenta(value string) string {
	return ansi.Magenta() + value + ansi.DefaultColor()
}

func Cyan(value string) string {
	return ansi.Cyan() + value + ansi.DefaultColor()
}

func White(value string) string {
	return ansi.White() + value + ansi.DefaultColor()
}

func BlackBg(value string) string {
	return ansi.BlackBg() + value + ansi.DefaultBg()
}

func RedBg(value string) string {
	return ansi.RedBg() + value + ansi.DefaultBg()
}

func GreenBg(value string) string {
	return ansi.GreenBg() + value + ansi.DefaultBg()
}

func YellowBg(value string) string {
	return ansi.YellowBg() + value + ansi.DefaultBg()
}

func BlueBg(value string) string {
	return ansi.BlueBg() + value + ansi.DefaultBg()
}

func MagentaBg(value string) string {
	return ansi.MagentaBg() + value + ansi.DefaultBg()
}

func CyanBg(value string) string {
	return ansi.CyanBg() + value + ansi.DefaultBg()
}

func WhiteBg(value string) string {
	return ansi.WhiteBg() + value + ansi.DefaultBg()
}
