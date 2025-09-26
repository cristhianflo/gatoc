package commands

import (
	"fmt"

	"github.com/bachacode/gatoc/internal/bot"
	subcommands "github.com/bachacode/gatoc/internal/commands/dolar"
	"github.com/bwmarrin/discordgo"
)

func init() {
	bot.RegisterCommand(dolar.Metadata.Name, dolar)
}

var dolar bot.SlashCommand = bot.SlashCommand{
	Metadata: &discordgo.ApplicationCommand{
		Name:        "dolar",
		Description: "Cotización del dolar a bolívares",
		Options: []*discordgo.ApplicationCommandOption{
			subcommands.Estado.Metadata,
			subcommands.Oficial.Metadata,
			subcommands.Paralelo.Metadata,
		},
	},
	Handler: func(s *discordgo.Session, i *discordgo.InteractionCreate, ctx *bot.BotContext) error {
		options := i.ApplicationCommandData().Options

		switch options[0].Name {
		case "estado":
			return subcommands.Estado.Handler(s, i, ctx)
		case "oficial":
			return subcommands.Oficial.Handler(s, i, ctx)
		case "paralelo":
			return subcommands.Paralelo.Handler(s, i, ctx)
		default:
			bot.GetInteractionFailedResponse(s, i, "El subcomando llamado no existe.")
			return fmt.Errorf("Subcommand doesn't exist\n")
		}
	},
}
