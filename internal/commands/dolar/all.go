package dolar

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/bachacode/gatoc/internal/bot"
	"github.com/bwmarrin/discordgo"
)

var All bot.SlashSubcommand = bot.SlashSubcommand{
	Metadata: &discordgo.ApplicationCommandOption{
		Type:        discordgo.ApplicationCommandOptionSubCommand,
		Name:        "all",
		Description: "Devuelve las distintas cotizaciones del Dólar en Venezuela",
	},
	Handler: func(s *discordgo.Session, i *discordgo.InteractionCreate, ctx *bot.BotContext) error {
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
			bot.GetInteractionFailedResponse(s, i, "Ha ocurrido un error al recibir la respuesta de la API.")
			return fmt.Errorf("API returned non-200 status: %d", resp.StatusCode)
		}
		var embedFields []*discordgo.MessageEmbedField

		for _, dolar := range response {
			if dolar.Name == "Bitcoin" {
				continue
			}

			updatedAt, err := time.Parse(time.RFC3339, dolar.UpdatedAt)
			average := strings.ReplaceAll(fmt.Sprintf("%.2f", dolar.Average), ".", ",")
			if err != nil {
				bot.GetInteractionFailedResponse(s, i, "Ha ocurrido un error interno.")
				return err
			}

			embedFields = append(embedFields, &discordgo.MessageEmbedField{
				Name:   fmt.Sprintf("%s - Promedio", dolar.Name),
				Value:  fmt.Sprintf("```fix\nBs. %s\n```", average),
				Inline: true,
			})

			embedFields = append(embedFields, &discordgo.MessageEmbedField{
				Name:   fmt.Sprintf("%s - Última Actualización", dolar.Name),
				Value:  fmt.Sprintf("```fix\n%s\n```", updatedAt.Format("02/01/2006 03:04:05 PM")),
				Inline: true,
			})

			if dolar.Name == "Paralelo" {
				continue
			}

			embedFields = append(embedFields, &discordgo.MessageEmbedField{
				Name:   "\u200b",
				Value:  "\u200b",
				Inline: false,
			})
		}

		embed := &discordgo.MessageEmbed{
			Title:       "Cotización",
			Description: "Cotización del Dólar Paralelo y Oficial (BCV) en Venezuela",
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
