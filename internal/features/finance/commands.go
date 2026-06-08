package finance

import (
	"fmt"

	"github.com/bachacode/gatoc/internal/bot"
	"github.com/bachacode/gatoc/internal/features/finance/subcommands"
	"github.com/bwmarrin/discordgo"
)

var dollar bot.SlashCommand = bot.SlashCommand{
	Metadata: &discordgo.ApplicationCommand{
		Name:        "dolar",
		Description: "Cotización del dolar a bolívares",
		Options: []*discordgo.ApplicationCommandOption{
			subcommands.DollarAll.Metadata,
			subcommands.DollarStatus.Metadata,
			subcommands.DollarBcv.Metadata,
			subcommands.DollarParalelo.Metadata,
		},
	},
	Handler: func(s *discordgo.Session, i *discordgo.InteractionCreate, ctx *bot.BotContext) error {
		options := i.ApplicationCommandData().Options

		switch options[0].Name {
		case "estado":
			return subcommands.DollarStatus.Handler(s, i, ctx)
		case "oficial":
			return subcommands.DollarBcv.Handler(s, i, ctx)
		case "paralelo":
			return subcommands.DollarParalelo.Handler(s, i, ctx)
		case "all":
			return subcommands.DollarAll.Handler(s, i, ctx)
		default:
			bot.GetInteractionFailedResponse(s, i, "El subcomando llamado no existe.")
			return fmt.Errorf("Subcommand doesn't exist\n")
		}
	},
}
