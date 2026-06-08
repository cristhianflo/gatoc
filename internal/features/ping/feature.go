package ping

import "github.com/bachacode/gatoc/internal/bot"

type PingFeature struct{}

func NewFeature() *PingFeature {
	return &PingFeature{}
}

func (f *PingFeature) SlashCommands() []bot.SlashCommand {
	return []bot.SlashCommand{
		ping,
	}
}

func (f *PingFeature) Models() []interface{} {
	return nil
}

func (f *PingFeature) RegisterEvents(router *bot.EventRouter) {
	// No events to register for this feature
}
