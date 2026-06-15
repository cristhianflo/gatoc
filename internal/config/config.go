package config

import (
	"os"
	"strings"
)

type Config struct {
	*BotConfig
	*DbConfig
	*RedisConfig
}

type BotConfig struct {
	Token         string
	ClientID      string
	GuildID       string
	MainChannelID string
	WelcomeEmoji  string
	GoodbyeEmoji  string
}

type DbConfig struct {
	DbHost  string
	DbUser  string
	DbPass  string
	DbName  string
	DbPort  string
	SslMode string
}

type RedisConfig struct {
	Url string
}

func LoadConfig() *Config {
	cfg := &Config{
		BotConfig: &BotConfig{
			Token:         getEnv("TOKEN", ""),
			ClientID:      getEnv("CLIENT_ID", ""),
			GuildID:       getEnv("GUILD_ID", ""),
			MainChannelID: getEnv("MAIN_CHANNEL_ID", ""),
			WelcomeEmoji:  getEnv("WELCOME_EMOJI", "<:gatoc:1356257759850663976>"),
			GoodbyeEmoji:  getEnv("GOODBYE_EMOJI", "<:sadcheems:1383154335675973834>"),
		},
		DbConfig: &DbConfig{
			DbHost:  getEnv("DB_HOST", "localhost"),
			DbUser:  getEnv("DB_USER", "postgres"),
			DbPass:  getEnv("DB_PASS", "password"),
			DbName:  getEnv("DB_NAME", "discordbot"),
			DbPort:  getEnv("DB_PORT", "5432"),
			SslMode: getEnv("DB_SSL", "disable"),
		},
		RedisConfig: &RedisConfig{
			Url: getEnv("REDIS_URL", "redis://localhost:6379"),
		},
	}

	return cfg
}

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}

	file, ok := os.LookupEnv(key + "_FILE")
	if !ok {
		return fallback
	}

	data, err := os.ReadFile(file)
	if err != nil {
		return fallback
	}

	return strings.TrimSpace(string(data))
}
