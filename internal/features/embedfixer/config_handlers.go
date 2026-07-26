package embedfixer

import (
	"fmt"
	"strings"

	"github.com/bachacode/gatoc/internal/bot"
	"github.com/bwmarrin/discordgo"
)

func handleConfigShow(s *discordgo.Session, i *discordgo.InteractionCreate, ctx *bot.BotContext) error {
	if err := bot.DeferReply(s, i); err != nil {
		bot.GetInteractionFailedResponse(s, i, "")
		return err
	}

	fields := make([]*discordgo.MessageEmbedField, 0, len(supportedPlatforms))
	for _, platform := range supportedPlatforms {
		domain, isCustom, err := resolveReplacementDomain(ctx.DB, i.GuildID, platform)
		if err != nil {
			bot.EditDeferred(s, i, "Failed to load embedfixer config")
			return err
		}

		mode := "default"
		if isCustom {
			mode = "custom"
		}

		fields = append(fields, &discordgo.MessageEmbedField{
			Name: platform.DisplayName,
			Value: fmt.Sprintf(
				"Hosts: `%s`\nDomain: `%s`\nMode: `%s`",
				strings.Join(platform.SourceHosts, "`, `"),
				domain,
				mode,
			),
			Inline: false,
		})
	}

	embed := &discordgo.MessageEmbed{
		Title:       "EmbedFixer Config",
		Description: "Current replacement-domain configuration by supported platform.",
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
