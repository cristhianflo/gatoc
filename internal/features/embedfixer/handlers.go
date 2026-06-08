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

func fixTwitterEmbed(m *discordgo.MessageCreate, url string, fixedUrl string) string {
	authorName := strings.Split(strings.Split(url, "/status")[0], ".com/")[1]
	author := "<" + strings.Split(url, "/status")[0] + ">"
	mention := fmt.Sprintf("<@%s>", m.Author.ID)
	fxtwitterURL := strings.Replace(url, "twitter.com", "fxtwitter.com", 1)
	fxtwitterURL = strings.Replace(fxtwitterURL, "x.com", "fxtwitter.com", 1)

	fixedEmbedMessageContent := fmt.Sprintf("[Tweet](%s) • [%s](%s) • [Fix](%s) • Enviado por %s ", fixedUrl, authorName, author, fxtwitterURL, mention)
	return fixedEmbedMessageContent
}

func fixRedditEmbed(m *discordgo.MessageCreate, url string, fixedUrl string) string {
	authorName := strings.Split(strings.Split(url, "/comments")[0], "r/")[1]
	author := "<" + strings.Split(url, "/comments")[0] + ">"
	mention := fmt.Sprintf("<@%s>", m.Author.ID)
	vxredditURL := strings.Replace(url, "reddit.com", "vxreddit.com", 1)

	fixedEmbedMessageContent := fmt.Sprintf("[Reddit](%s) • [%s](%s) • [Fix](%s) • Enviado por %s ", fixedUrl, authorName, author, vxredditURL, mention)
	return fixedEmbedMessageContent
}

func fixInstagramEmbed(m *discordgo.MessageCreate, url string, fixedUrl string) string {
	mention := fmt.Sprintf("<@%s>", m.Author.ID)
	fixedInstagramURL := strings.Replace(url, "instagram.com", "kkinstagram.com", 1)

	fixedEmbedMessageContent := fmt.Sprintf("[Instagram](%s) • [Fix](%s) • Enviado por %s ", fixedUrl, fixedInstagramURL, mention)
	return fixedEmbedMessageContent
}

func (f *EmbedFixerFeature) EmbedFixerHandler(s *discordgo.Session, m *discordgo.MessageCreate, ctx *bot.BotContext) error {
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

	fixableHosts := map[string]func(m *discordgo.MessageCreate, url string, fixedUrl string) string{
		"twitter.com":   fixTwitterEmbed,
		"x.com":         fixTwitterEmbed,
		"reddit.com":    fixRedditEmbed,
		"instagram.com": fixInstagramEmbed,
	}

	if handler, ok := fixableHosts[trimmedHost]; ok {
		if (trimmedHost == "twitter.com" || trimmedHost == "x.com") && !strings.Contains(url.Path, "/status") {
			return nil
		}

		supressedUrl := "<" + urlToParse + ">"
		fixedEmbedMessageContent := handler(m, urlToParse, supressedUrl)

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
	}

	return nil

}
