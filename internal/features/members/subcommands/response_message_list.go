package subcommands

import (
	"fmt"
	"time"

	"github.com/bachacode/gatoc/internal/bot"
	"github.com/bachacode/gatoc/internal/database"
	"github.com/bwmarrin/discordgo"
)

var ResponseMessageList bot.SlashSubcommand = bot.SlashSubcommand{
	Metadata: &discordgo.ApplicationCommandOption{
		Type:        discordgo.ApplicationCommandOptionSubCommand,
		Name:        "lista",
		Description: "Ver todas las respuestas a mensajes",
	},
	Handler: func(s *discordgo.Session, i *discordgo.InteractionCreate, ctx *bot.BotContext) error {

		if err := bot.DeferReply(s, i); err != nil {
			bot.GetInteractionFailedResponse(s, i, "")
			return err
		}

		db := ctx.DB
		content := ""

		var responses []database.ResponseMessage
		result := db.Find(&responses)

		if result.Error != nil {
			content = "Ha ocurrido un error para obtener la lista de respuestas"
			bot.EditDeferred(s, i, content)
			return fmt.Errorf("Error responding to interaction: %v\n", result.Error)
		}

		fields := []*discordgo.MessageEmbedField{
			{
				Name:   "ID",
				Value:  "",
				Inline: true,
			},
			{
				Name:   "Para",
				Value:  "",
				Inline: true,
			},
			{
				Name:   "Mensaje",
				Value:  "",
				Inline: true,
			},
			{
				Name:   "Respuesta",
				Value:  "",
				Inline: true,
			},
		}

		for _, response := range responses {

			fields[0].Value += fmt.Sprintf("%d\n", response.ID)
			fields[2].Value += response.Message + "\n"
			fields[3].Value += response.Response + "\n"

			if response.UserID == nil {
				fields[1].Value += "Todos\n"
				continue
			}

			user, err := s.User(*response.UserID)

			if err != nil {
				content = "Ha ocurrido un error obteniendo un usuario"
				bot.EditDeferred(s, i, content)
				return fmt.Errorf("Error getting an user: %s from the guild: %s\n%v", *response.UserID, i.GuildID, err)
			}

			if user.GlobalName != "" {
				fields[1].Value += user.GlobalName + "\n"
			} else {
				fields[1].Value += user.Username + "\n"
			}
		}

		embed := discordgo.MessageEmbed{
			Title:       "Respuestas",
			Description: "Respuestas a mensajes especificos",
			Fields:      fields,
			Color:       0xFFFFFF,
			Footer: &discordgo.MessageEmbedFooter{
				Text: fmt.Sprintf("Generated at %s", time.Now().Format("2006-01-02 15:04:05")),
			},
		}

		// Send the embed to the channel
		_, err := s.InteractionResponseEdit(
			i.Interaction,
			&discordgo.WebhookEdit{
				Embeds: &[]*discordgo.MessageEmbed{&embed},
			},
		)
		if err != nil {
			return fmt.Errorf("Error sending embed: %v\n", err)
		}

		return nil
	},
}
