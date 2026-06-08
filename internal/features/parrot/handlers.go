package parrot

import (
	"fmt"
	"strings"

	"github.com/bachacode/gatoc/internal/bot"
	"github.com/bwmarrin/discordgo"
)

func (f *ParrotFeature) ParrotHandler(s *discordgo.Session, m *discordgo.MessageCreate, ctx *bot.BotContext) error {
	channelID := m.ChannelID
	messages, err := s.ChannelMessages(channelID, 2, "", "", "")
	if err != nil {
		return err
	}

	if len(messages) == 0 {
		return fmt.Errorf("Failed to get enough messages from message history")
	}

	isSameMessage := strings.EqualFold(messages[0].Content, messages[1].Content)
	isDifferentAuthor := messages[0].Author.GlobalName != messages[1].Author.GlobalName
	if isSameMessage && isDifferentAuthor {
		f.messageCount++
	} else {
		f.messageCount = 1
	}

	if f.messageCount >= f.maxMessages {
		f.messageCount = 0
		_, err := s.ChannelMessageSend(channelID, messages[len(messages)-1].Content)
		if err != nil {
			return err
		}
	}

	return nil
}
