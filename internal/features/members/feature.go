package members

import "github.com/bachacode/gatoc/internal/bot"

type MembersFeature struct{}

func NewFeature() *MembersFeature {
	return &MembersFeature{}
}

func (f *MembersFeature) SlashCommands() []bot.SlashCommand {
	return []bot.SlashCommand{
		welcomeRole,
		responseMessage,
	}
}

func (f *MembersFeature) Models() []interface{} {
	return nil
}

func (f *MembersFeature) RegisterEvents(router *bot.EventRouter) {
	router.OnMessageCreate(f.BotResponseHandler)
	router.OnGuildMemberAdd(f.GuildMemberAddRoleHandler)
	router.OnGuildMemberAdd(f.GuildMemberWelcomeHandler)
	router.OnGuildMemberRemove(f.GuildMemberRemoveMessageHandler)
}
