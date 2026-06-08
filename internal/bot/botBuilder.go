package bot

import (
	"errors"
	"fmt"
	"log"
	"os"

	"github.com/bachacode/gatoc/internal/config"
	"github.com/bwmarrin/discordgo"
	"gorm.io/gorm"
)

type BotBuilder struct {
	// Mandatory configuration
	cfg *config.BotConfig

	// Optional dependencies/configurations
	db          *gorm.DB
	logger      *log.Logger
	intents     discordgo.Intent
	commands    map[string]SlashCommand
	eventRouter *EventRouter
}

func NewBotBuilder(cfg *config.BotConfig) *BotBuilder {
	return &BotBuilder{
		cfg:         cfg,
		logger:      log.New(os.Stdout, "[DEFAULT_BOT] ", log.LstdFlags|log.Lshortfile),
		intents:     discordgo.IntentGuildMessages,
		commands:    make(map[string]SlashCommand),
		eventRouter: NewEventRouter(),
	}
}

func (bb *BotBuilder) WithDatabase(db *gorm.DB) *BotBuilder {
	bb.db = db
	return bb
}

func (bb *BotBuilder) WithLogger(logger *log.Logger) *BotBuilder {
	bb.logger = logger
	return bb
}

func (bb *BotBuilder) WithIntents(intents discordgo.Intent) *BotBuilder {
	bb.intents = intents
	return bb
}

func (bb *BotBuilder) WithFeatures(f []Feature) *BotBuilder {
	for _, feature := range f {
		for _, cmd := range feature.SlashCommands() {
			bb.commands[cmd.Metadata.Name] = cmd
		}

		feature.RegisterEvents(bb.eventRouter)
	}
	return bb
}

func (bb *BotBuilder) Build() (*bot, error) {
	// Validate that base config and token exist
	if bb.cfg == nil {
		return nil, errors.New("bot configuration is required")
	}
	if bb.cfg.Token == "" {
		return nil, errors.New("Discord bot token is required in the configuration")
	}

	// Create new discord session
	s, err := discordgo.New("Bot " + bb.cfg.Token)
	if err != nil {
		return nil, fmt.Errorf("failed to create Discord session: %w", err)
	}

	// Create bot context
	botCtx := &BotContext{
		BotConfig: bb.cfg,
		DB:        bb.db,
		Logger:    bb.logger,
	}

	// Create new bot struct
	b := &bot{
		session:     s,
		intents:     bb.intents,
		BotContext:  botCtx,
		eventRouter: bb.eventRouter,
	}

	// Create commands
	if len(bb.commands) > 0 {
		botCtx.Logger.Println("INFO: Creating commands...")
		b.commands = bb.commands
	} else {
		botCtx.Logger.Println("WARN: No commands provided, bot will start without any registered commands")
	}

	botCtx.Logger.Println("INFO: Bot instance successfully built")
	return b, nil
}
