package dolar

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/bachacode/gatoc/internal/bot"
	"github.com/bwmarrin/discordgo"
)

type StatusResponse struct {
	Status string `json:"estado"`
	Random int    `json:"aleatorio"`
}

var Estado bot.SlashSubcommand = bot.SlashSubcommand{
	Metadata: &discordgo.ApplicationCommandOption{
		Type:        discordgo.ApplicationCommandOptionSubCommand,
		Name:        "estado",
		Description: "Devuelve el estado de la API",
	},
	Handler: func(s *discordgo.Session, i *discordgo.InteractionCreate, ctx *bot.BotContext) error {
		apiUrl := "https://ve.dolarapi.com/v1/estado"

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

		var response StatusResponse
		if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
			bot.EditDeferred(s, i, "Error al decodificar la respuesta de la API.")
			return err
		}

		if resp.StatusCode != 200 {
			bot.GetInteractionFailedResponse(s, i, "Ha ocurrido un error al recibir la respuesta de la API.")
			return fmt.Errorf("API returned non-200 status: %d", resp.StatusCode)
		}

		embed := &discordgo.MessageEmbed{
			Title:       "Estado de la API del Dólar",
			Description: "Información sobre el estado actual de la API del Dólar.",
			Color:       0x00ff00, // Verde
			Fields: []*discordgo.MessageEmbedField{
				{
					Name:   "Estado",
					Value:  response.Status,
					Inline: true,
				},
				{
					Name:   "Código Aleatorio",
					Value:  fmt.Sprintf("%d", response.Random),
					Inline: true,
				},
			},
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
