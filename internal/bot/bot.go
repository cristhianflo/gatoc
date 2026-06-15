package bot

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/bachacode/gatoc/internal/config"
	"github.com/bwmarrin/discordgo"
	"github.com/go-redis/redis/v8"
	"gorm.io/gorm"
)

type BotContext struct {
	*config.BotConfig
	DB     *gorm.DB
	Redis  *redis.Client
	Logger *log.Logger
}

type bot struct {
	session     *discordgo.Session
	intents     discordgo.Intent
	commands    map[string]SlashCommand
	eventRouter *EventRouter
	*BotContext
}

func (b *bot) Setup() {
	b.Logger.Println("INFO: Setting bot events...")
	b.SetupEvents()

	b.Logger.Println("INFO: Registering commands...")
	b.registerCommands()
}

func (b *bot) isDuplicateEvent(ctx context.Context, eventKey string) bool {
	// Attempt to set the key with a 10-second TTL.
	// NX means "Only set the key if it does not already exist."
	success, err := b.Redis.SetNX(ctx, "event_lock:"+eventKey, "1", 10*time.Second).Result()
	if err != nil {
		b.Logger.Printf("WARN: Redis error: %v", err)
		return false // Fail-safe: allow processing if Redis goes down
	}
	// If success is false, the key already existed -> it's a duplicate
	return !success
}

func (b *bot) SetupEvents() {
	// 1. Ready Event Dispatcher
	b.session.AddHandlerOnce(func(s *discordgo.Session, r *discordgo.Ready) {
		log.Printf("Logged in as: %v#%v", s.State.User.Username, s.State.User.Discriminator)
		for _, handler := range b.eventRouter.readyHandlers {
			if err := handler(s, r, b.BotContext); err != nil {
				b.Logger.Printf("Middleware error on Ready: %v", err)
			}
		}
	})

	// 2. Core Interaction Create Dispatcher
	b.session.AddHandler(func(s *discordgo.Session, i *discordgo.InteractionCreate) {
		if i.Type != discordgo.InteractionApplicationCommand {
			return
		}

		if b.isDuplicateEvent(context.Background(), "msg:"+i.ID) {
			return
		}

		cmd, ok := b.commands[i.ApplicationCommandData().Name]
		if !ok {
			b.Logger.Printf("WARN: Received interaction for unknown command: %s", i.ApplicationCommandData().Name)
			return
		}

		if err := cmd.Handler(s, i, b.BotContext); err != nil {
			fmt.Printf("Failed to run interaction: %v\n", err)
		}

		for _, handler := range b.eventRouter.interactionCreateHandlers {
			if err := handler(s, i, b.BotContext); err != nil {
				b.Logger.Printf("Middleware error on InteractionCreate: %v", err)
			}
		}
	})

	// 3. Core Message Create Dispatcher
	b.session.AddHandler(func(s *discordgo.Session, m *discordgo.MessageCreate) {
		if m.Author.ID == s.State.User.ID || m.Author.Bot {
			return
		}

		if b.isDuplicateEvent(context.Background(), "interaction:"+m.ID) {
			return
		}

		for _, handler := range b.eventRouter.messageCreateHandlers {
			if err := handler(s, m, b.BotContext); err != nil {
				b.Logger.Printf("Middleware error on MessageCreate: %v", err)
			}
		}
	})

	// 4. Core Guild Member Add Dispatcher
	b.session.AddHandler(func(s *discordgo.Session, m *discordgo.GuildMemberAdd) {
		key := fmt.Sprintf("member_add:%s:%s", m.GuildID, m.User.ID)
		if b.isDuplicateEvent(context.Background(), key) {
			return
		}

		for _, handler := range b.eventRouter.guildMemberAddHandlers {
			if err := handler(s, m, b.BotContext); err != nil {
				b.Logger.Printf("Middleware error on GuildMemberAdd: %v", err)
			}
		}
	})

	// 5. Core Guild Member Remove Dispatcher
	b.session.AddHandler(func(s *discordgo.Session, m *discordgo.GuildMemberRemove) {
		key := fmt.Sprintf("member_remove:%s:%s", m.GuildID, m.User.ID)
		if b.isDuplicateEvent(context.Background(), key) {
			return
		}

		for _, handler := range b.eventRouter.guildMemberRemoveHandlers {
			if err := handler(s, m, b.BotContext); err != nil {
				b.Logger.Printf("Middleware error on GuildMemberRemove: %v", err)
			}
		}
	})

	// 6. Core Message Reaction Add Dispatcher
	b.session.AddHandler(func(s *discordgo.Session, r *discordgo.MessageReactionAdd) {
		key := fmt.Sprintf("reaction_add:%s:%s", r.GuildID, r.UserID)
		if b.isDuplicateEvent(context.Background(), key) {
			return
		}

		for _, handler := range b.eventRouter.messageReactionAddHandlers {
			if err := handler(s, r, b.BotContext); err != nil {
				b.Logger.Printf("Middleware error on MessageReactionAdd: %v", err)
			}
		}
	})
}

func (b *bot) registerCommands() {
	total := len(b.commands)
	count := 0
	for _, cmd := range b.commands {
		_, err := b.session.ApplicationCommandCreate(b.ClientID, b.GuildID, cmd.Metadata)

		if err != nil {
			b.Logger.Printf("WARN: Failed to register command %s: %v\n", cmd.Metadata.Name, err)
		} else {
			count++
			b.Logger.Printf("INFO: Registered command: %s\n", cmd.Metadata.Name)
		}
	}
	b.Logger.Printf("INFO: %d commands of %d were register successfully!", count, total)
}

func (b *bot) UnregisterCommands() {
	commands, err := b.session.ApplicationCommands(b.ClientID, b.GuildID)

	if err != nil {
		b.Logger.Printf("WARN: Failed to fetch applications commands: %v", err)
		b.Logger.Println("WARN: Skipping commands removal...")
		return
	}

	total := len(commands)
	count := 0

	for _, cmd := range commands {
		err := b.session.ApplicationCommandDelete(b.ClientID, b.GuildID, cmd.ID)
		if err != nil {
			b.Logger.Printf("WARN: Failed to unregister command %s: %v\n", cmd.Name, err)
		} else {
			count++
			b.Logger.Printf("INFO: Unregistered command: %s\n", cmd.Name)
		}
	}
	b.Logger.Printf("%d commands of %d were unregister successfully!\n", count, total)
}

func (b *bot) Run() error {

	b.session.Identify.Intents = b.intents
	b.Logger.Println("INFO: Starting bot session...")
	if err := b.session.Open(); err != nil {
		return fmt.Errorf("Error starting bot session: %v", err)
	}

	b.Logger.Println("INFO: Bot is now running. Press CTRL-C to exit.")

	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	<-sc

	// b.Logger.Println("INFO: Unregistering commands...")
	// b.UnregisterCommands()

	b.Logger.Println("INFO: Closing bot session...")
	b.session.Close()
	if err := b.session.Open(); err != nil {
		return fmt.Errorf("Error closing bot session: %v", err)
	}

	return nil
}
