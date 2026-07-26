package embedfixer

import (
	"fmt"
	"net/url"
	"regexp"
	"strings"
	"time"

	"github.com/bachacode/gatoc/internal/bot"
	"github.com/bwmarrin/discordgo"
)

func hasNoFixTag(content string) bool {
	parts := strings.Fields(strings.ToLower(content))
	for _, part := range parts {
		if part == "#nofix" {
			return true
		}
	}

	return false
}

func fixedURLFromDomain(rawURL string, replacementDomain string) (string, error) {
	parsedURL, err := url.ParseRequestURI(rawURL)
	if err != nil {
		return "", err
	}

	parsedURL.Host = replacementDomain
	return parsedURL.String(), nil
}

func buildFixedEmbedMessage(m *discordgo.MessageCreate, platform platformConfig, sourceURL string, fixedURL string) string {
	mention := fmt.Sprintf("<@%s>", m.Author.ID)
	suppressedURL := "<" + sourceURL + ">"
	return fmt.Sprintf("[%s](%s) • [Fix](%s) • Enviado por %s", platform.LinkLabel, suppressedURL, fixedURL, mention)
}

func (f *EmbedFixerFeature) EmbedFixerHandler(s *discordgo.Session, m *discordgo.MessageCreate, ctx *bot.BotContext) error {
	if hasNoFixTag(m.Content) {
		return nil
	}

	re := regexp.MustCompile(`https?://[^\s]+`)

	urlToParse := re.FindString(m.Content)

	if urlToParse == "" {
		return nil
	}

	url, err := url.ParseRequestURI(urlToParse)
	if err != nil {
		return nil
	}

	trimmedHost := strings.ToLower(url.Host)
	trimmedHost = strings.TrimPrefix(trimmedHost, "www.")

	platform, ok := platformByHost(trimmedHost)
	if !ok {
		return nil
	}

	if platform.Key == "twitter" && !strings.Contains(url.Path, "/status") {
		return nil
	}

	replacementDomain, _, err := resolveReplacementDomain(ctx.DB, m.GuildID, platform)
	if err != nil {
		return err
	}

	fixedURL, err := fixedURLFromDomain(urlToParse, replacementDomain)
	if err != nil {
		return nil
	}

	fixedEmbedMessageContent := buildFixedEmbedMessage(m, platform, urlToParse, fixedURL)

	go func() {
		maxRetries := 5
		for range maxRetries {
			s.ChannelMessageEditComplex(&discordgo.MessageEdit{
				ID:      m.ID,
				Channel: m.ChannelID,
				Flags:   discordgo.MessageFlagsSuppressEmbeds,
			})

			time.Sleep(300 * time.Millisecond)
		}
	}()

	if _, err := s.ChannelMessageSendComplex(m.ChannelID, &discordgo.MessageSend{
		Content: fixedEmbedMessageContent,
		AllowedMentions: &discordgo.MessageAllowedMentions{
			Parse: []discordgo.AllowedMentionType{},
		},
	}); err != nil {
		return err
	}

	return nil
}
