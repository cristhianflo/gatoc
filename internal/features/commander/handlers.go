package commander

import (
	"fmt"

	"github.com/bachacode/gatoc/internal/bot"
	"github.com/bwmarrin/discordgo"
)

func (f *CommanderFeature) CommandHandler(s *discordgo.Session, i *discordgo.InteractionCreate, ctx *bot.BotContext) error {
	cmd := bot.GetCommand(i.ApplicationCommandData().Name)
	if err := cmd.Handler(s, i, ctx); err != nil {
		fmt.Printf("Failed to run interaction: %v\n", err)
	}
	return nil
}
