package main

import (
	"log"
	"os"

	"github.com/bachacode/gatoc/internal/bot"
	"github.com/bachacode/gatoc/internal/config"
	"github.com/bachacode/gatoc/internal/database"
	"github.com/bachacode/gatoc/internal/features/commander"
	"github.com/bachacode/gatoc/internal/features/embedfixer"
	"github.com/bachacode/gatoc/internal/features/finance"
	"github.com/bachacode/gatoc/internal/features/members"
	"github.com/bachacode/gatoc/internal/features/parrot"
	"github.com/bachacode/gatoc/internal/features/ping"
	"github.com/bwmarrin/discordgo"
)

func main() {
	logger := log.New(os.Stdout, "[MAIN] ", log.LstdFlags|log.Lshortfile)
	cfg := config.LoadConfig()

	db, err := database.New(cfg.DbConfig)
	if err != nil {
		logger.Fatalf("ERROR: Failed to connect to db: %v\n", err)
		return
	}

	if err := database.Migrate(db); err != nil {
		logger.Fatalf("ERROR: Failed to migrate tables to db: %v\n", err)
	}

	features := []bot.Feature{
		ping.NewFeature(),
		members.NewFeature(),
		finance.NewFeature(),
		parrot.NewFeature(),
		embedfixer.NewFeature(),
		commander.NewFeature(),
	}

	bb := bot.NewBotBuilder(cfg.BotConfig)
	bb.WithDatabase(db)
	bb.WithIntents(
		discordgo.IntentsGuilds |
			discordgo.IntentsGuildVoiceStates |
			discordgo.IntentsMessageContent |
			discordgo.IntentGuildMessages |
			discordgo.IntentsGuildMembers |
			discordgo.IntentsGuildMessageReactions,
	)
	bb.WithLogger(logger)
	bb.WithFeatures(features)
	bot, err := bb.Build()
	if err != nil {
		logger.Fatalf("ERROR: Failed to create bot instance: %v\n", err)
		return
	}

	bot.Setup()

	if err := bot.Run(); err != nil {
		logger.Fatalf("ERROR: Failed to run bot instance: %v\n", err)
		return
	}
}
