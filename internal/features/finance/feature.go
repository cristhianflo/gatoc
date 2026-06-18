package finance

import "github.com/bachacode/gatoc/internal/bot"

type FinanceFeature struct{}

func NewFeature() *FinanceFeature {
	return &FinanceFeature{}
}

func (f *FinanceFeature) SlashCommands() []bot.SlashCommand {
	return []bot.SlashCommand{
		f.dollarCommand(),
		f.convertCommand(),
	}
}

func (f *FinanceFeature) Models() []interface{} {
	return nil
}

func (f *FinanceFeature) RegisterEvents(router *bot.EventRouter) {
	// No events to register for this feature
}
