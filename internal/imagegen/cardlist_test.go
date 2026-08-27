package imagegen

import (
	"bytes"
	"image/png"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func testCards(n int) []Card {
	cards := make([]Card, n)
	for i := range cards {
		cards[i] = Card{
			Label:   "DÓLAR OFICIAL (USD)",
			Value:   "Bs. 787,52",
			Caption: "Fuente: BCV · 26 Ago 2026, 12:00 AM",
		}
	}
	return cards
}

func TestRenderCardListProducesPNG(t *testing.T) {
	out, err := RenderCardList(CardListOptions{
		Title:    "MONITOR DE DIVISAS VENEZUELA",
		Subtitle: "Cotización del Dólar (Oficial/Paralelo) y Euro Oficial",
		Note:     "No disponible temporalmente: EUR",
		Cards:    testCards(3),
	})
	require.NoError(t, err)

	cfg, err := png.DecodeConfig(bytes.NewReader(out))
	require.NoError(t, err)
	assert.Equal(t, canvasWidth, cfg.Width)
	assert.Equal(t, cardListHeight(3, true, true, true), cfg.Height)
	assert.Equal(t, 686, cfg.Height, "unexpected layout metrics drift")
}

func TestRenderCardListMinimal(t *testing.T) {
	out, err := RenderCardList(CardListOptions{Cards: testCards(1)})
	require.NoError(t, err)

	cfg, err := png.DecodeConfig(bytes.NewReader(out))
	require.NoError(t, err)
	assert.Equal(t, canvasWidth, cfg.Width)
	assert.Equal(t, cardListHeight(1, false, false, false), cfg.Height)
}

func TestRenderCardListRequiresCards(t *testing.T) {
	_, err := RenderCardList(CardListOptions{Title: "Solo título"})
	assert.Error(t, err)
}
