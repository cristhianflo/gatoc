package bot

import "github.com/bwmarrin/discordgo"

type MessageCreateMiddleware func(s *discordgo.Session, m *discordgo.MessageCreate, ctx *BotContext) error
type GuildMemberAddMiddleware func(s *discordgo.Session, m *discordgo.GuildMemberAdd, ctx *BotContext) error
type GuildMemberRemoveMiddleware func(s *discordgo.Session, m *discordgo.GuildMemberRemove, ctx *BotContext) error
type ReadyMiddleware func(s *discordgo.Session, r *discordgo.Ready, ctx *BotContext) error
type InteractionCreateMiddleware func(s *discordgo.Session, i *discordgo.InteractionCreate, ctx *BotContext) error
type MessageReactionAddMiddleware func(s *discordgo.Session, r *discordgo.MessageReactionAdd, ctx *BotContext) error

type EventRouter struct {
	messageCreateHandlers      []MessageCreateMiddleware
	guildMemberAddHandlers     []GuildMemberAddMiddleware
	guildMemberRemoveHandlers  []GuildMemberRemoveMiddleware
	readyHandlers              []ReadyMiddleware
	interactionCreateHandlers  []InteractionCreateMiddleware
	messageReactionAddHandlers []MessageReactionAddMiddleware
}

func NewEventRouter() *EventRouter {
	return &EventRouter{}
}

func (r *EventRouter) OnMessageCreate(handler MessageCreateMiddleware) {
	r.messageCreateHandlers = append(r.messageCreateHandlers, handler)
}

func (r *EventRouter) OnGuildMemberAdd(handler GuildMemberAddMiddleware) {
	r.guildMemberAddHandlers = append(r.guildMemberAddHandlers, handler)
}

func (r *EventRouter) OnGuildMemberRemove(handler GuildMemberRemoveMiddleware) {
	r.guildMemberRemoveHandlers = append(r.guildMemberRemoveHandlers, handler)
}

func (r *EventRouter) OnReady(handler ReadyMiddleware) {
	r.readyHandlers = append(r.readyHandlers, handler)
}

func (r *EventRouter) OnInteractionCreate(handler InteractionCreateMiddleware) {
	r.interactionCreateHandlers = append(r.interactionCreateHandlers, handler)
}

func (r *EventRouter) OnMessageReactionAdd(handler MessageReactionAddMiddleware) {
	r.messageReactionAddHandlers = append(r.messageReactionAddHandlers, handler)
}
