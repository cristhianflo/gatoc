package finance

import (
	"fmt"

	"github.com/bachacode/gatoc/internal/bot"
	"github.com/bachacode/gatoc/internal/features/finance/subcommands"
	"github.com/bwmarrin/discordgo"
)

func (f *FinanceFeature) dollarCommand() bot.SlashCommand {
	return bot.SlashCommand{
		Metadata: &discordgo.ApplicationCommand{
			Name:        "dollar",
			Description: "Dollar to Bolivares exchange rates",

			Options: []*discordgo.ApplicationCommandOption{
				subcommands.DollarAll.Metadata,
				subcommands.DollarStatus.Metadata,
			},
		},
		Handler: func(s *discordgo.Session, i *discordgo.InteractionCreate, ctx *bot.BotContext) error {
			options := i.ApplicationCommandData().Options

			switch options[0].Name {
			case subcommands.DollarStatus.Metadata.Name:
				return subcommands.DollarStatus.Handler(s, i, ctx)
			case subcommands.DollarAll.Metadata.Name:
				return subcommands.DollarAll.Handler(s, i, ctx)
			default:
				bot.GetInteractionFailedResponse(s, i, "El subcomando llamado no existe.")
				return fmt.Errorf("Subcommand doesn't exist\n")
			}
		},
	}
}

func (f *FinanceFeature) convertCommand() bot.SlashCommand {
	return bot.SlashCommand{
		Metadata: &discordgo.ApplicationCommand{
			Name: "test",
		},
		Handler: func(s *discordgo.Session, i *discordgo.InteractionCreate, ctx *bot.BotContext) error {
			return nil
		},
	}
}
