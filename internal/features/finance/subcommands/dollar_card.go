package subcommands

import (
	"fmt"
	"strings"
	"time"

	"github.com/bachacode/gatoc/internal/imagegen"
)

var spanishMonths = []string{"Ene", "Feb", "Mar", "Abr", "May", "Jun", "Jul", "Ago", "Sep", "Oct", "Nov", "Dic"}

func formatSpanishTimestamp(t time.Time) string {
	return fmt.Sprintf("%d %s %d, %s", t.Day(), spanishMonths[t.Month()-1], t.Year(), t.Format("03:04 PM"))
}

func cardFromRate(rate DolarResponse) (imagegen.Card, error) {
	updatedAt, err := time.Parse(time.RFC3339, rate.UpdatedAt)
	if err != nil {
		return imagegen.Card{}, err
	}

	name := rate.Name
	if name == "" {
		name = "N/A"
	}
	currency := rate.Currency
	if currency == "" {
		currency = "N/A"
	}
	source := rate.Source
	if source == "" {
		source = "N/A"
	}

	average := strings.ReplaceAll(fmt.Sprintf("%.2f", rate.Average), ".", ",")

	return imagegen.Card{
		Label:   fmt.Sprintf("%s (%s)", strings.ToUpper(name), strings.ToUpper(currency)),
		Value:   fmt.Sprintf("Bs. %s", average),
		Caption: fmt.Sprintf("Fuente: %s · %s", strings.ToUpper(source), formatSpanishTimestamp(updatedAt)),
	}, nil
}

func buildCardOptions(results map[string]currencyRatesResult) (cards []imagegen.Card, unavailable []string) {
	for _, currency := range []string{"USD", "EUR"} {
		result, ok := results[currency]
		if !ok || result.Err != nil {
			unavailable = append(unavailable, currency)
			continue
		}

		added := false
		for _, rate := range result.Rates {
			if currency == "USD" && strings.EqualFold(rate.Name, "Bitcoin") {
				continue
			}

			card, err := cardFromRate(rate)
			if err != nil {
				break
			}
			cards = append(cards, card)
			added = true
		}

		if !added {
			unavailable = append(unavailable, currency)
		}
	}
	return cards, unavailable
}

func buildCardListOptions(cards []imagegen.Card, unavailable []string) imagegen.CardListOptions {
	opts := imagegen.CardListOptions{
		Title:    "MONITOR DE DIVISAS VENEZUELA",
		Subtitle: "Cotización del Dólar (Oficial/Paralelo) y Euro Oficial",
		Cards:    cards,
	}
	if len(unavailable) > 0 {
		opts.Note = "No disponible temporalmente: " + strings.Join(unavailable, ", ")
	}
	return opts
}
