package parrot

import (
	"github.com/bachacode/gatoc/internal/bot"
)

type ParrotFeature struct {
	messageCount int
	maxMessages  int
}

func NewFeature() *ParrotFeature {
	return &ParrotFeature{
		messageCount: 0,
		maxMessages:  3, // Set the threshold for how many times a message should be repeated before parroting
	}
}

func (f *ParrotFeature) SlashCommands() []bot.SlashCommand { return nil }
func (f *ParrotFeature) Models() []interface{}             { return nil }

func (f *ParrotFeature) RegisterEvents(router *bot.EventRouter) {
	// Attach the parrot logic to the MessageCreate event!
	router.OnMessageCreate(f.ParrotHandler)
}
