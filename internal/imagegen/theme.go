package imagegen

import "image/color"

// Theme holds the colors used when rendering an image.
type Theme struct {
	Background    color.RGBA
	CardBG        color.RGBA
	Border        color.RGBA
	TextPrimary   color.RGBA
	TextSecondary color.RGBA
	Accent        color.RGBA
	Note          color.RGBA
}

// DefaultTheme returns the standard dark theme: near-black background,
// green accent values and gray secondary text.
func DefaultTheme() *Theme {
	return &Theme{
		Background:    color.RGBA{R: 0x16, G: 0x19, B: 0x1f, A: 0xff},
		CardBG:        color.RGBA{R: 0x1f, G: 0x24, B: 0x2c, A: 0xff},
		Border:        color.RGBA{R: 0x3b, G: 0x41, B: 0x4b, A: 0xff},
		TextPrimary:   color.RGBA{R: 0xf2, G: 0xf4, B: 0xf7, A: 0xff},
		TextSecondary: color.RGBA{R: 0x9a, G: 0xa3, B: 0xaf, A: 0xff},
		Accent:        color.RGBA{R: 0x4a, G: 0xde, B: 0x80, A: 0xff},
		Note:          color.RGBA{R: 0xfb, G: 0xbf, B: 0x24, A: 0xff},
	}
}
