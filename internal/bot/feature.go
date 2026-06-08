package bot

type Feature interface {
	SlashCommands() []SlashCommand
	Models() []interface{}

	RegisterEvents(router *EventRouter)
}
