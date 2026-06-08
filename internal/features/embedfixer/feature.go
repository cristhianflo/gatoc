package embedfixer

import "github.com/bachacode/gatoc/internal/bot"

type EmbedFixerFeature struct{}

func NewFeature() *EmbedFixerFeature {
	return &EmbedFixerFeature{}
}

func (f *EmbedFixerFeature) SlashCommands() []bot.SlashCommand {
	return nil
}

func (f *EmbedFixerFeature) Models() []interface{} {
	return nil
}

func (f *EmbedFixerFeature) RegisterEvents(router *bot.EventRouter) {
	router.OnMessageCreate(f.EmbedFixerHandler)
}
