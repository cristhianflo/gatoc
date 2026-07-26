package embedfixer

import (
	"fmt"
	"strings"

	"github.com/bachacode/gatoc/internal/bot"
	"github.com/bwmarrin/discordgo"
)

func modeDisplay(isCustom bool) string {
	if isCustom {
		return "CUSTOM"
	}

	return "DEFAULT"
}

func handleConfigShow(s *discordgo.Session, i *discordgo.InteractionCreate, ctx *bot.BotContext) error {
	if err := bot.DeferReply(s, i); err != nil {
		bot.GetInteractionFailedResponse(s, i, "")
		return err
	}

	fields := make([]*discordgo.MessageEmbedField, 0, len(supportedPlatforms)*4)
	for index, platform := range supportedPlatforms {
		domain, isCustom, err := resolveReplacementDomain(ctx.DB, i.GuildID, platform)
		if err != nil {
			bot.EditDeferred(s, i, "Failed to load embedfixer config")
			return err
		}

		hosts := strings.Join(platform.SourceHosts, "\n")
		mode := modeDisplay(isCustom)

		fields = append(fields,
			&discordgo.MessageEmbedField{
				Name:   fmt.Sprintf("%s - Hosts", platform.DisplayName),
				Value:  fmt.Sprintf("```fix\n%s\n```", hosts),
				Inline: true,
			},
			&discordgo.MessageEmbedField{
				Name:   fmt.Sprintf("%s - Domain", platform.DisplayName),
				Value:  fmt.Sprintf("```fix\n%s\n```", domain),
				Inline: true,
			},
			&discordgo.MessageEmbedField{
				Name:   fmt.Sprintf("%s - Mode", platform.DisplayName),
				Value:  fmt.Sprintf("```fix\n%s\n```", mode),
				Inline: true,
			},
		)

		if index < len(supportedPlatforms)-1 {
			fields = append(fields, &discordgo.MessageEmbedField{
				Name:   "\u200b",
				Value:  "\u200b",
				Inline: false,
			})
		}
	}

	embed := &discordgo.MessageEmbed{
		Title:       "EmbedFixer Config",
		Description: "Estado de configuración por plataforma (tabla).",
		Color:       0x00ff00,
		Fields:      fields,
	}

	_, err := s.InteractionResponseEdit(i.Interaction, &discordgo.WebhookEdit{Embeds: &[]*discordgo.MessageEmbed{embed}})
	return err
}

func handleConfigSet(s *discordgo.Session, i *discordgo.InteractionCreate, ctx *bot.BotContext, subcommand *discordgo.ApplicationCommandInteractionDataOption) error {
	if err := bot.DeferReply(s, i); err != nil {
		bot.GetInteractionFailedResponse(s, i, "")
		return err
	}

	optionMap := map[string]*discordgo.ApplicationCommandInteractionDataOption{}
	for _, opt := range subcommand.Options {
		optionMap[opt.Name] = opt
	}

	platformValue, ok := optionMap["platform"]
	if !ok {
		bot.EditDeferred(s, i, "Missing platform option")
		return fmt.Errorf("missing platform option")
	}

	domainValue, ok := optionMap["domain"]
	if !ok {
		bot.EditDeferred(s, i, "Missing domain option")
		return fmt.Errorf("missing domain option")
	}

	platform, found := platformByKey(platformValue.StringValue())
	if !found {
		bot.EditDeferred(s, i, fmt.Sprintf("Unsupported platform. Use: %s", buildPlatformChoices()))
		return fmt.Errorf("unsupported platform: %s", platformValue.StringValue())
	}

	normalizedDomain, err := normalizeDomainHost(domainValue.StringValue())
	if err != nil {
		bot.EditDeferred(s, i, "Invalid domain. Provide only a host, for example: fxtwitter.com")
		return err
	}

	if err := setCustomDomain(ctx.DB, i.GuildID, platform.Key, normalizedDomain); err != nil {
		bot.EditDeferred(s, i, "Failed to save domain config")
		return err
	}

	bot.EditDeferred(s, i, fmt.Sprintf("%s now uses custom domain `%s`", platform.DisplayName, normalizedDomain))
	return nil
}

func handleConfigReset(s *discordgo.Session, i *discordgo.InteractionCreate, ctx *bot.BotContext, subcommand *discordgo.ApplicationCommandInteractionDataOption) error {
	if err := bot.DeferReply(s, i); err != nil {
		bot.GetInteractionFailedResponse(s, i, "")
		return err
	}

	optionMap := map[string]*discordgo.ApplicationCommandInteractionDataOption{}
	for _, opt := range subcommand.Options {
		optionMap[opt.Name] = opt
	}

	platformValue, ok := optionMap["platform"]
	if !ok {
		bot.EditDeferred(s, i, "Missing platform option")
		return fmt.Errorf("missing platform option")
	}

	platform, found := platformByKey(platformValue.StringValue())
	if !found {
		bot.EditDeferred(s, i, fmt.Sprintf("Unsupported platform. Use: %s", buildPlatformChoices()))
		return fmt.Errorf("unsupported platform: %s", platformValue.StringValue())
	}

	if err := resetCustomDomain(ctx.DB, i.GuildID, platform.Key); err != nil {
		bot.EditDeferred(s, i, "Failed to reset domain config")
		return err
	}

	bot.EditDeferred(s, i, fmt.Sprintf("%s reset to default domain `%s`", platform.DisplayName, platform.DefaultDomain))
	return nil
}
