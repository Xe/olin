package ansi

// Reset turns off all SGR attributes
func Reset() string {
	return SelectGraphicsRendition(0)
}

// Bold increases font intensity
func Bold() string {
	return SelectGraphicsRendition(1)
}

// BoldOff turns off bold effects
func BoldOff() string {
	// 21 not widely supported, so we're using 22 here instead
	return SelectGraphicsRendition(22)
}

// Faint decreases font intensity
func Faint() string {
	return SelectGraphicsRendition(2)
}

// FaintOff turns off faint effects
func FaintOff() string {
	return SelectGraphicsRendition(22)
}

// Italic makes the font italicized. Sometimes treated as inverse. (Not widely supported)
func Italic() string {
	return SelectGraphicsRendition(3)
}

// ItalicOff turns off italicized effect.
func ItalicOff() string {
	return SelectGraphicsRendition(23)
}

// Underline adds an underline to the text.
func Underline() string {
	return SelectGraphicsRendition(4)
}

// UnderlineOff removes the underline from the text.
func UnderlineOff() string {
	return SelectGraphicsRendition(24)
}

// BlinkSlow applies slow blinking effects.
func BlinkSlow() string {
	return SelectGraphicsRendition(5)
}

// BlinkFast applies fast blinking effects.
func BlinkFast() string {
	return SelectGraphicsRendition(6)
}

// BlinkOff removes all blinking effects.
func BlinkOff() string {
	return SelectGraphicsRendition(25)
}

// ImageNegative applies a negative image effect.
func ImageNegative() string {
	return SelectGraphicsRendition(7)
}

// ImagePosition removes the negative image effect.
func ImagePositive() string {
	return SelectGraphicsRendition(27)
}

// Conceal is not widely supported.
func Conceal() string {
	return SelectGraphicsRendition(8)
}

// Reveal turns off conceal mode.
func Reveal() string {
	return SelectGraphicsRendition(28)
}

// CrossedOut applies a strike-through to the text.
func CrossedOut() string {
	return SelectGraphicsRendition(9)
}

// CrossedOutOff removes the strike-through from the text.
func CrossedOutOff() string {
	return SelectGraphicsRendition(29)
}

// DefaultColor resets the font color.
func DefaultColor() string {
	return SelectGraphicsRendition(39)
}

// Black sets the font color to black.
func Black() string {
	return SelectGraphicsRendition(30)
}

// Red sets the font color to red.
func Red() string {
	return SelectGraphicsRendition(31)
}

// Green sets the font color to green.
func Green() string {
	return SelectGraphicsRendition(32)
}

// Yellow sets the font color to yellow.
func Yellow() string {
	return SelectGraphicsRendition(33)
}

// Blue sets the font color to blue.
func Blue() string {
	return SelectGraphicsRendition(34)
}

// Magenta sets the font color to magenta.
func Magenta() string {
	return SelectGraphicsRendition(35)
}

// Cyan sets the font color to cyan.
func Cyan() string {
	return SelectGraphicsRendition(36)
}

// White sets the font color to white.
func White() string {
	return SelectGraphicsRendition(37)
}

// DefaultBg resets the background color.
func DefaultBg() string {
	return SelectGraphicsRendition(49)
}

// BlackBg sets the background color to black.
func BlackBg() string {
	return SelectGraphicsRendition(40)
}

// RedBg sets the background color to red.
func RedBg() string {
	return SelectGraphicsRendition(41)
}

// GreenBg sets the background color to green.
func GreenBg() string {
	return SelectGraphicsRendition(42)
}

// YellowBg sets the background color to yellow.
func YellowBg() string {
	return SelectGraphicsRendition(43)
}

// BlueBg sets the background color to blue.
func BlueBg() string {
	return SelectGraphicsRendition(44)
}

// MagentaBg sets the background color to magenta.
func MagentaBg() string {
	return SelectGraphicsRendition(45)
}

// CyanBg sets the background color to cyan.
func CyanBg() string {
	return SelectGraphicsRendition(46)
}

// WhiteBg sets the background color to white.
func WhiteBg() string {
	return SelectGraphicsRendition(47)
}
