package subcommands

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/bachacode/gatoc/internal/bot"
	"github.com/bwmarrin/discordgo"
)

var ConvertParalelo bot.SlashSubcommand = bot.SlashSubcommand{
	Metadata: &discordgo.ApplicationCommandOption{
		Type:        discordgo.ApplicationCommandOptionSubCommand,
		Name:        "paralelo",
		Description: "Converts dollars from Parallel rate to BCV rate",
		Options: []*discordgo.ApplicationCommandOption{
			{
				Type:        discordgo.ApplicationCommandOptionNumber,
				Name:        "amount",
				Description: "Amount in USD (Parallel) to convert",
				Required:    true,
			},
		},
	},
	Handler: func(s *discordgo.Session, i *discordgo.InteractionCreate, ctx *bot.BotContext) error {
		options := i.ApplicationCommandData().Options

		optionMap := make(map[string]*discordgo.ApplicationCommandInteractionDataOption)
		for _, opt := range options[0].Options {
			optionMap[opt.Name] = opt
		}

		amountOption, ok := optionMap["amount"]
		if !ok {
			bot.GetInteractionFailedResponse(s, i, "Amount not provided.")
			return nil
		}
		amount := amountOption.FloatValue()

		apiUrl := "https://ve.dolarapi.com/v1/dolares"

		if err := bot.DeferReply(s, i); err != nil {
			bot.GetInteractionFailedResponse(s, i, "")
			return err
		}

		resp, err := http.Get(apiUrl)
		if err != nil {
			bot.EditDeferred(s, i, "Error al conectar con la API.")
			return err
		}
		defer resp.Body.Close()

		var response []DolarResponse
		if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
			bot.EditDeferred(s, i, "Error al decodificar la respuesta de la API.")
			return err
		}

		if resp.StatusCode != 200 {
			bot.EditDeferred(s, i, "Ha ocurrido un error al recibir la respuesta de la API.")
			return fmt.Errorf("API returned non-200 status: %d", resp.StatusCode)
		}

		var bcvRate, paraleloRate float32
		for _, dolar := range response {
			if dolar.Source == string(OfficialDollarSource) {
				bcvRate = dolar.Average
			}
			if dolar.Source == string(ParallelDollarSource) {
				paraleloRate = dolar.Average
			}
		}

		if bcvRate == 0 || paraleloRate == 0 {
			bot.EditDeferred(s, i, "No se pudieron obtener las tasas necesarias para la conversión.")
			return nil
		}

		bolivares := float64(paraleloRate) * amount
		resultadoBCV := bolivares / float64(bcvRate)

		formatVez := func(f float64) string {
			return strings.ReplaceAll(fmt.Sprintf("%.2f", f), ".", ",")
		}

		embed := &discordgo.MessageEmbed{
			Title: "Conversión: Paralelo ➔ BCV",
			Color: 0x00ff00,
			Fields: []*discordgo.MessageEmbedField{
				{
					Name:   "Monto Original (Paralelo)",
					Value:  fmt.Sprintf("```fix\n$ %s\n```", formatVez(amount)),
					Inline: true,
				},
				{
					Name:   "Tasa Paralelo",
					Value:  fmt.Sprintf("```fix\nBs. %s\n```", formatVez(float64(paraleloRate))),
					Inline: true,
				},
				{
					Name:   "Equivalente en Bolívares",
					Value:  fmt.Sprintf("```fix\nBs. %s\n```", formatVez(bolivares)),
					Inline: false,
				},
				{
					Name:   "Tasa BCV",
					Value:  fmt.Sprintf("```fix\nBs. %s\n```", formatVez(float64(bcvRate))),
					Inline: true,
				},
				{
					Name:   "Resultado (BCV)",
					Value:  fmt.Sprintf("```fix\n$ %s\n```", formatVez(resultadoBCV)),
					Inline: true,
				},
			},
			Footer: &discordgo.MessageEmbedFooter{
				Text: "gatoc - Conversión de divisas",
			},
		}

		if _, err := s.InteractionResponseEdit(i.Interaction, &discordgo.WebhookEdit{
			Embeds: &[]*discordgo.MessageEmbed{embed},
		}); err != nil {
			return fmt.Errorf("Error responding to interaction: %v", err)
		}

		return nil
	},
}
