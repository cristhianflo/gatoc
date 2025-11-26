package commands

import (
	"fmt"

	"github.com/bachacode/gatoc/internal/bot"
	"github.com/bwmarrin/discordgo"
)

func init() {
	bot.RegisterCommand(ping.Metadata.Name, ping)
}

var ping bot.SlashCommand = bot.SlashCommand{
	Metadata: &discordgo.ApplicationCommand{
		Name:        "gatoping",
		Description: "Devuelve la latencia en MS",
	},
	Handler: func(s *discordgo.Session, i *discordgo.InteractionCreate, ctx *bot.BotContext) error {
		latency := s.HeartbeatLatency().Milliseconds()

		// Follow up with the actual latency
		embed := &discordgo.MessageEmbed{
			Title: "🏓 | GatoPong!",
			Fields: []*discordgo.MessageEmbedField{
				{
					Name:   "GatoLatencia",
					Value:  fmt.Sprintf("```fix\n⚡ | %dms\n```", latency),
					Inline: true,
				},
				{
					Name:   "GatoVersión",
					Value:  "```fix\n1.2.1\n```",
					Inline: true,
				},
			},
		}
		err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Embeds: []*discordgo.MessageEmbed{
					embed,
				},
			},
		})
		if err != nil {
			bot.GetInteractionFailedResponse(s, i, "")
			return fmt.Errorf("Error responding to interaction: %v", err)
		}
		return nil
	},
}
