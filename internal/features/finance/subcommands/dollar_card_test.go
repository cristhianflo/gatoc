package subcommands

import (
	"errors"
	"testing"
	"time"

	"github.com/bachacode/gatoc/internal/imagegen"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestFormatSpanishTimestamp(t *testing.T) {
	tests := []struct {
		name string
		time time.Time
		want string
	}{
		{
			name: "midnight UTC",
			time: time.Date(2026, time.August, 26, 0, 0, 0, 0, time.UTC),
			want: "26 Ago 2026, 12:00 AM",
		},
		{
			name: "evening",
			time: time.Date(2026, time.August, 26, 21, 1, 0, 0, time.UTC),
			want: "26 Ago 2026, 09:01 PM",
		},
		{
			name: "different month",
			time: time.Date(2025, time.February, 3, 14, 30, 0, 0, time.UTC),
			want: "3 Feb 2025, 02:30 PM",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, formatSpanishTimestamp(tt.time))
		})
	}
}

func TestCardFromRate(t *testing.T) {
	rate := DolarResponse{
		Currency:  "usd",
		Source:    "bcv",
		Name:      "Dólar Oficial",
		Average:   787.52,
		UpdatedAt: "2026-08-26T00:00:00.000Z",
	}

	card, err := cardFromRate(rate)
	require.NoError(t, err)
	assert.Equal(t, "DÓLAR OFICIAL (USD)", card.Label)
	assert.Equal(t, "Bs. 787,52", card.Value)
	assert.Equal(t, "Fuente: BCV · 26 Ago 2026, 12:00 AM", card.Caption)
}

func TestCardFromRateEmptyFieldsFallBackToNA(t *testing.T) {
	card, err := cardFromRate(DolarResponse{Average: 1, UpdatedAt: "2026-08-26T00:00:00.000Z"})
	require.NoError(t, err)
	assert.Equal(t, "N/A (N/A)", card.Label)
	assert.Equal(t, "Fuente: N/A · 26 Ago 2026, 12:00 AM", card.Caption)
}

func TestCardFromRateInvalidTimestamp(t *testing.T) {
	_, err := cardFromRate(DolarResponse{UpdatedAt: "not-a-time"})
	assert.Error(t, err)
}

func usdRates() []DolarResponse {
	return []DolarResponse{
		{
			Currency:  "usd",
			Source:    "bcv",
			Name:      "Dólar Oficial",
			Average:   787.52,
			UpdatedAt: "2026-08-26T00:00:00.000Z",
		},
		{
			Currency:  "usd",
			Source:    "enparalelovzla",
			Name:      "Dólar Paralelo",
			Average:   931.28,
			UpdatedAt: "2026-08-26T21:01:00.000Z",
		},
		{
			Currency:  "usd",
			Source:    "binance",
			Name:      "Bitcoin",
			Average:   1.00,
			UpdatedAt: "2026-08-26T21:01:00.000Z",
		},
	}
}

func eurRates() []DolarResponse {
	return []DolarResponse{
		{
			Currency:  "eur",
			Source:    "bcv",
			Name:      "Euro Oficial",
			Average:   919.15,
			UpdatedAt: "2026-08-26T00:00:00.000Z",
		},
	}
}

func TestBuildCardOptionsBothSourcesAvailable(t *testing.T) {
	cards, unavailable := buildCardOptions(map[string]currencyRatesResult{
		"USD": {Rates: usdRates()},
		"EUR": {Rates: eurRates()},
	})

	require.Len(t, cards, 3)
	assert.Equal(t, "DÓLAR OFICIAL (USD)", cards[0].Label)
	assert.Equal(t, "DÓLAR PARALELO (USD)", cards[1].Label)
	assert.Equal(t, "EURO OFICIAL (EUR)", cards[2].Label)
	assert.Empty(t, unavailable)
}

func TestBuildCardOptionsEURUnavailable(t *testing.T) {
	cards, unavailable := buildCardOptions(map[string]currencyRatesResult{
		"USD": {Rates: usdRates()},
		"EUR": {Err: errors.New("boom")},
	})

	require.Len(t, cards, 2)
	assert.Equal(t, []string{"EUR"}, unavailable)
}

func TestBuildCardOptionsUSDUnavailable(t *testing.T) {
	cards, unavailable := buildCardOptions(map[string]currencyRatesResult{
		"USD": {Err: errors.New("boom")},
		"EUR": {Rates: eurRates()},
	})

	require.Len(t, cards, 1)
	assert.Equal(t, "EURO OFICIAL (EUR)", cards[0].Label)
	assert.Equal(t, []string{"USD"}, unavailable)
}

func TestBuildCardOptionsUSDOnlyBitcoin(t *testing.T) {
	cards, unavailable := buildCardOptions(map[string]currencyRatesResult{
		"USD": {Rates: usdRates()[2:]},
		"EUR": {Rates: eurRates()},
	})

	require.Len(t, cards, 1)
	assert.Equal(t, []string{"USD"}, unavailable)
}

func TestBuildCardListOptionsWithoutUnavailable(t *testing.T) {
	cards, _ := buildCardOptions(map[string]currencyRatesResult{
		"USD": {Rates: usdRates()},
		"EUR": {Rates: eurRates()},
	})

	opts := buildCardListOptions(cards, nil)
	assert.Equal(t, "MONITOR DE DIVISAS VENEZUELA", opts.Title)
	assert.Equal(t, "Cotización del Dólar (Oficial/Paralelo) y Euro Oficial", opts.Subtitle)
	assert.Empty(t, opts.Note)
	require.Len(t, opts.Cards, 3)
}

func TestBuildCardListOptionsWithUnavailable(t *testing.T) {
	opts := buildCardListOptions([]imagegen.Card{{Label: "x"}}, []string{"EUR"})
	assert.Equal(t, "No disponible temporalmente: EUR", opts.Note)
}
