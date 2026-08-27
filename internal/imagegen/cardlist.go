package imagegen

import (
	"bytes"
	"fmt"
	"image"
	"image/color"
	"image/png"

	"github.com/fogleman/gg"
)

// Canvas geometry (pixels). Heights derive from the font sizes below.
const (
	canvasWidth = 1024

	padX      = 40.0
	padTop    = 44
	padBottom = 44

	titleSize   = 40
	titleLineH  = 52
	gapTitle    = 10
	subtitleSze = 22
	subtitleLnH = 30
	gapToCards  = 30

	cardW       = 944 // canvasWidth - 2*padX
	cardH       = 128
	cardGap     = 20
	cardRadius  = 14
	cardPadX    = 28.0
	cardPadY    = 15
	labelSize   = 19
	labelLineH  = 26
	valueSize   = 36
	valueLineH  = 46
	captionSize = 18
	captionLnH  = 26

	gapToNote = 26
	noteSze   = 18
	noteLineH = 26
)

// cardListHeight returns the canvas height for a card-list image.
func cardListHeight(cardCount int, withTitle, withSubtitle, withNote bool) int {
	h := padTop
	if withTitle {
		h += titleLineH + gapTitle
	}
	if withSubtitle {
		h += subtitleLnH + gapToCards
	}
	h += cardCount*cardH + (cardCount-1)*cardGap
	if withNote {
		h += gapToNote + noteLineH
	}
	return h + padBottom
}

// newContext returns a gg context of the given size filled with bg.
func newContext(w, h int, bg color.RGBA) *gg.Context {
	dc := gg.NewContext(w, h)
	dc.SetColor(bg)
	dc.Clear()
	return dc
}

// RenderCardList renders a CardList image and returns it encoded as PNG.
func RenderCardList(opts CardListOptions) ([]byte, error) {
	if len(opts.Cards) == 0 {
		return nil, fmt.Errorf("imagegen: at least one card is required")
	}

	img, err := renderCardList(opts)
	if err != nil {
		return nil, err
	}

	var buf bytes.Buffer
	if err := png.Encode(&buf, img); err != nil {
		return nil, fmt.Errorf("imagegen: encode png: %w", err)
	}
	return buf.Bytes(), nil
}

func renderCardList(opts CardListOptions) (image.Image, error) {
	regular, bold, err := loadFonts()
	if err != nil {
		return nil, fmt.Errorf("imagegen: %w", err)
	}

	theme := opts.Theme
	if theme == nil {
		theme = DefaultTheme()
	}

	withTitle := opts.Title != ""
	withSubtitle := opts.Subtitle != ""
	withNote := opts.Note != ""

	dc := newContext(canvasWidth, cardListHeight(len(opts.Cards), withTitle, withSubtitle, withNote), theme.Background)

	y := float64(padTop)
	if withTitle {
		dc.SetFontFace(newFace(bold, titleSize))
		dc.SetColor(theme.TextPrimary)
		dc.DrawStringAnchored(opts.Title, canvasWidth/2, y+titleLineH/2, 0.5, 0.5)
		y += titleLineH + gapTitle
	}
	if withSubtitle {
		dc.SetFontFace(newFace(regular, subtitleSze))
		dc.SetColor(theme.TextSecondary)
		dc.DrawStringAnchored(opts.Subtitle, canvasWidth/2, y+subtitleLnH/2, 0.5, 0.5)
		y += subtitleLnH + gapToCards
	}

	for _, card := range opts.Cards {
		dc.DrawRoundedRectangle(padX, y, cardW, cardH, cardRadius)
		dc.SetColor(theme.CardBG)
		dc.Fill()
		dc.SetColor(theme.Border)
		dc.SetLineWidth(2)
		dc.Stroke()

		x := padX + cardPadX
		dc.SetFontFace(newFace(regular, labelSize))
		dc.SetColor(theme.TextSecondary)
		dc.DrawStringAnchored(card.Label, x, y+cardPadY+labelLineH/2, 0, 0.5)

		dc.SetFontFace(newFace(bold, valueSize))
		dc.SetColor(theme.Accent)
		dc.DrawStringAnchored(card.Value, x, y+cardPadY+labelLineH+valueLineH/2, 0, 0.5)

		dc.SetFontFace(newFace(regular, captionSize))
		dc.SetColor(theme.TextSecondary)
		dc.DrawStringAnchored(card.Caption, x, y+cardPadY+labelLineH+valueLineH+captionLnH/2, 0, 0.5)

		y += cardH + cardGap
	}

	if withNote {
		dc.SetFontFace(newFace(regular, noteSze))
		dc.SetColor(theme.Note)
		dc.DrawStringAnchored(opts.Note, canvasWidth/2, y+gapToNote+noteLineH/2, 0.5, 0.5)
	}

	return dc.Image(), nil
}
