package events

import (
	"fmt"
	"net/url"
	"regexp"
	"strings"
	"time"

	"github.com/bachacode/gatoc/internal/bot"
	"github.com/bachacode/gatoc/internal/database"
	"github.com/bwmarrin/discordgo"
	"gorm.io/gorm"
)

func init() {
	bot.RegisterEvent(messageCreate)
}

var messageCount = 0

var messageCreate bot.Event = bot.Event{
	Name: "Message Create",
	Once: false,
	Handler: func(ctx *bot.BotContext) interface{} {
		return func(s *discordgo.Session, m *discordgo.MessageCreate) {
			if m.Author.ID == s.State.User.ID || m.Author.Bot {
				return
			}

			channelID := m.ChannelID
			maxMessages := 3

			if err := HandleResponses(s, m, ctx.DB); err != nil {
				ctx.Logger.Printf("Failed to respond to message: %v", err)
			}

			if err := handleRepeated(channelID, maxMessages, s); err != nil {
				ctx.Logger.Printf("Failed to repeat messages: %v", err)
			}

			if err := handleURLEmbed(s, m); err != nil {
				ctx.Logger.Printf("Failed to fix message embed: %v", err)
			}

		}
	},
}

func HandleResponses(s *discordgo.Session, m *discordgo.MessageCreate, db *gorm.DB) error {
	var response database.ResponseMessage
	result := db.Where("message = ?", strings.ToLower(m.Content)).First(&response)

	if result.Error != nil {
		fmt.Printf("Failed to get response message: %v\n", result.Error)
		return result.Error
	}

	if result.RowsAffected < 1 {
		return nil
	}

	if response.UserID != nil && *response.UserID != m.Author.ID {
		return nil
	}

	msg := &discordgo.MessageSend{
		Content: response.Response,
	}

	if response.UserID != nil && *response.UserID == m.Author.ID {
		msg.Reference = m.Reference()
	}

	if _, err := s.ChannelMessageSendComplex(m.ChannelID, msg); err != nil {
		return err
	}

	return nil
}

func handleRepeated(channelID string, max int, s *discordgo.Session) error {
	messages, err := s.ChannelMessages(channelID, 2, "", "", "")
	if err != nil {
		return err
	}

	if len(messages) == 0 {
		return fmt.Errorf("Failed to get enough messages from message history")
	}

	isSameMessage := strings.ToLower(messages[0].Content) == strings.ToLower(messages[1].Content)
	isDifferentAuthor := messages[0].Author.GlobalName != messages[1].Author.GlobalName
	if isSameMessage && isDifferentAuthor {
		messageCount++
	} else {
		messageCount = 1
	}

	if messageCount >= max {
		messageCount = 0
		_, err := s.ChannelMessageSend(channelID, messages[len(messages)-1].Content)
		if err != nil {
			return err
		}
	}

	return nil
}

func handleURLEmbed(s *discordgo.Session, m *discordgo.MessageCreate) error {
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
	if strings.HasPrefix(trimmedHost, "www.") {
		trimmedHost = strings.TrimPrefix(trimmedHost, "www.")
	}

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
			for i := 0; i < maxRetries; i++ {
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
