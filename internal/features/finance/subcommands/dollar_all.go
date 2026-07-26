package subcommands

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/bachacode/gatoc/internal/bot"
	"github.com/bwmarrin/discordgo"
)

const (
	usdRatesURL = "https://ve.dolarapi.com/v1/dolares"
	eurRateURL  = "https://ve.dolarapi.com/v1/euros/oficial"
)

type currencyRatesResult struct {
	Rates []DolarResponse
	Err   error
}

func fetchRates(client *http.Client, apiURL string) currencyRatesResult {
	resp, err := client.Get(apiURL)
	if err != nil {
		return currencyRatesResult{Err: err}
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return currencyRatesResult{Err: fmt.Errorf("API returned non-200 status: %d", resp.StatusCode)}
	}
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return currencyRatesResult{Err: err}
	}

	var rates []DolarResponse
	if err := json.Unmarshal(body, &rates); err == nil {
		return currencyRatesResult{Rates: rates}
	}

	var singleRate DolarResponse
	if err := json.Unmarshal(body, &singleRate); err == nil {
		return currencyRatesResult{Rates: []DolarResponse{singleRate}}
	}

	return currencyRatesResult{Err: fmt.Errorf("unsupported response shape from %s", apiURL)}
}

func formatRateFields(rate DolarResponse) ([]*discordgo.MessageEmbedField, error) {
	updatedAt, err := time.Parse(time.RFC3339, rate.UpdatedAt)
	if err != nil {
		return nil, err
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
	header := fmt.Sprintf("%s (%s) - %s", rate.Name, currency, strings.ToUpper(source))

	return []*discordgo.MessageEmbedField{
		{
			Name:   fmt.Sprintf("%s - Promedio", header),
			Value:  fmt.Sprintf("```fix\nBs. %s\n```", average),
			Inline: true,
		},
		{
			Name:   fmt.Sprintf("%s - Última Actualización", header),
			Value:  fmt.Sprintf("```fix\n%s\n```", updatedAt.Format("02/01/2006 03:04:05 PM")),
			Inline: true,
		},
	}, nil
}

var DollarAll bot.SlashSubcommand = bot.SlashSubcommand{
	Metadata: &discordgo.ApplicationCommandOption{
		Type:        discordgo.ApplicationCommandOptionSubCommand,
		Name:        "all",
		Description: "Returns all the different Dollar exchange rates in Venezuela",
	},
	Handler: func(s *discordgo.Session, i *discordgo.InteractionCreate, ctx *bot.BotContext) error {
		if err := bot.DeferReply(s, i); err != nil {
			bot.GetInteractionFailedResponse(s, i, "")
			return err
		}

		httpClient := &http.Client{Timeout: 10 * time.Second}
		results := make(map[string]currencyRatesResult)
		var mu sync.Mutex
		var wg sync.WaitGroup

		endpoints := map[string]string{
			"USD": usdRatesURL,
			"EUR": eurRateURL,
		}

		for currency, apiURL := range endpoints {
			wg.Add(1)
			go func(currency string, apiURL string) {
				defer wg.Done()

				result := fetchRates(httpClient, apiURL)

				mu.Lock()
				results[currency] = result
				mu.Unlock()
			}(currency, apiURL)
		}

		wg.Wait()
		var embedFields []*discordgo.MessageEmbedField
		var unavailable []string

		for _, currency := range []string{"USD", "EUR"} {
			result, ok := results[currency]
			if !ok || result.Err != nil {
				unavailable = append(unavailable, currency)
				continue
			}

			for _, rate := range result.Rates {
				if currency == "USD" && strings.EqualFold(rate.Name, "Bitcoin") {
					continue
				}

				fields, err := formatRateFields(rate)
				if err != nil {
					unavailable = append(unavailable, currency)
					break
				}

				embedFields = append(embedFields, fields...)
				embedFields = append(embedFields, &discordgo.MessageEmbedField{
					Name:   "\u200b",
					Value:  "\u200b",
					Inline: false,
				})
			}
		}

		if len(embedFields) > 0 {
			embedFields = embedFields[:len(embedFields)-1]
		}

		if len(embedFields) == 0 {
			bot.EditDeferred(s, i, "No se pudo obtener cotizaciones en este momento.")
			return fmt.Errorf("all upstream sources failed")
		}

		description := "Cotización del Dólar (Oficial/Paralelo) y Euro Oficial en Venezuela"
		if len(unavailable) > 0 {
			description = fmt.Sprintf("%s\n\nNo disponible temporalmente: %s", description, strings.Join(unavailable, ", "))
		}

		embed := &discordgo.MessageEmbed{
			Title:       "Cotización",
			Description: description,
			Color:       0x00ff00,
			Fields:      embedFields,
		}

		if _, err := s.InteractionResponseEdit(i.Interaction, &discordgo.WebhookEdit{
			Embeds: &[]*discordgo.MessageEmbed{
				embed,
			},
		}); err != nil {
			bot.GetInteractionFailedResponse(s, i, "")
			return fmt.Errorf("Error responding to interaction: %v", err)

		}

		return nil
	},
}
