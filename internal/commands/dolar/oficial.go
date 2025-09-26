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

type DolarResponse struct {
	Source    string  `json:"fuente"`
	Name      string  `json:"nombre"`
	Buy       int     `json:"compra"`
	Sell      int     `json:"venta"`
	Average   float32 `json:"promedio"`
	UpdatedAt string  `json:"fechaActualizacion"`
}

var Oficial bot.SlashSubcommand = bot.SlashSubcommand{
	Metadata: &discordgo.ApplicationCommandOption{
		Type:        discordgo.ApplicationCommandOptionSubCommand,
		Name:        "oficial",
		Description: "Devuelve la cotización del Dólar Oficial",
	},
	Handler: func(s *discordgo.Session, i *discordgo.InteractionCreate, ctx *bot.BotContext) error {
		apiUrl := "https://ve.dolarapi.com/v1/dolares/oficial"

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

		var response DolarResponse
		if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
			bot.EditDeferred(s, i, "Error al decodificar la respuesta de la API.")
			return err
		}

		if resp.StatusCode != 200 {
			bot.GetInteractionFailedResponse(s, i, "Ha ocurrido un error al recibir la respuesta de la API.")
			return fmt.Errorf("API returned non-200 status: %d", resp.StatusCode)
		}

		updatedAt, err := time.Parse(time.RFC3339, response.UpdatedAt)
		average := strings.ReplaceAll(fmt.Sprintf("%.2f", response.Average), ".", ",")
		if err != nil {
			bot.GetInteractionFailedResponse(s, i, "Ha ocurrido un error interno.")
			return err
		}

		embed := &discordgo.MessageEmbed{
			Title:       "Cotización",
			Description: "Cotización del Dólar Oficial (BCV) en Venezuela",
			Color:       0x00ff00,
			Fields: []*discordgo.MessageEmbedField{
				{
					Name:   "Promedio",
					Value:  fmt.Sprintf("```fix\nBs. %s\n```", average),
					Inline: true,
				},
				{
					Name:   "Última Actualización",
					Value:  fmt.Sprintf("```fix\n%s\n```", updatedAt.Format("02/01/2006 03:04:05 PM")),
					Inline: false,
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
