package commander

import "github.com/bachacode/gatoc/internal/bot"

type CommanderFeature struct{}

func NewFeature() *CommanderFeature {
	return &CommanderFeature{}
}

func (f *CommanderFeature) SlashCommands() []bot.SlashCommand { return nil }

func (f *CommanderFeature) Models() []interface{} { return nil }

func (f *CommanderFeature) RegisterEvents(router *bot.EventRouter) {
	router.OnInteractionCreate(f.CommandHandler)
}
