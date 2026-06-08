package members

import (
	"errors"
	"fmt"
	"math/rand"
	"strings"

	"github.com/bachacode/gatoc/internal/bot"
	"github.com/bachacode/gatoc/internal/database"
	"github.com/bwmarrin/discordgo"
	"gorm.io/gorm"
)

func (f *MembersFeature) BotResponseHandler(s *discordgo.Session, m *discordgo.MessageCreate, ctx *bot.BotContext) error {

	db := ctx.DB

	var response database.ResponseMessage
	result := db.Where("message = ?", strings.ToLower(m.Content)).First(&response)

	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil
	}

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

func (f *MembersFeature) GuildMemberWelcomeHandler(s *discordgo.Session, r *discordgo.GuildMemberAdd, ctx *bot.BotContext) error {
	channelID := ctx.MainChannelID
	emoji := ctx.WelcomeEmoji
	mention := "<@" + r.User.ID + ">"

	embed := discordgo.MessageSend{
		Embeds: []*discordgo.MessageEmbed{
			{
				Title:       "qlq " + emoji + " 🍷",
				Color:       0xFFFFFF,
				Description: mention + " acaba de cometer el error mas grande de su vida entrando a esta tierra profana.",
				Image: &discordgo.MessageEmbedImage{
					URL: "https://media.tenor.com/eH-RoS91Q1gAAAAC/cat.gif",
				},
			},
		},
	}
	_, err := s.ChannelMessageSendComplex(channelID, &embed)

	if err != nil {
		fmt.Printf("Failed to get main channel: %v\n", err)
	}
	return nil
}

func (f *MembersFeature) GuildMemberRemoveMessageHandler(s *discordgo.Session, r *discordgo.GuildMemberRemove, ctx *bot.BotContext) error {
	channelID := ctx.MainChannelID
	memberID := r.User.ID
	mention := "<@" + memberID + ">"
	messages := []LeaveMessage{
		{
			embed: &discordgo.MessageEmbed{
				Title:       "c lo acomodaron por las costillas <:sadcheems:869742943425151087>",
				Color:       0xFFFFFF,
				Description: mention + " no aguanto la pela.",
				Image: &discordgo.MessageEmbedImage{
					URL: "https://media.tenor.com/ww56Kix_vM8AAAAC/seloacomodoporlascostillas.gif",
				},
			},
			filename: "chavez.gif",
		},
		{
			embed: &discordgo.MessageEmbed{
				Title:       "c le fue la luz <:sadcheems:869742943425151087>",
				Color:       0xFFFFFF,
				Description: mention + " no aguanto la pela.",
				Image: &discordgo.MessageEmbedImage{
					URL: "https://media.tenor.com/vHMD9o7RmfYAAAAC/snake-salute.gif",
				},
			},
			filename: "snake.gif",
		},
	}

	randNumber := rand.Intn(len(messages))

	selectedMessage := messages[randNumber]

	embed := discordgo.MessageSend{
		Embeds: []*discordgo.MessageEmbed{
			selectedMessage.embed,
		},
	}
	_, err := s.ChannelMessageSendComplex(channelID, &embed)

	if err != nil {
		fmt.Printf("Failed to get main channel: %v\n", err)
		return err
	}

	return nil
}

func (f *MembersFeature) GuildMemberAddRoleHandler(s *discordgo.Session, r *discordgo.GuildMemberAdd, ctx *bot.BotContext) error {
	db := ctx.DB
	var wRoles []database.WelcomeRole
	result := db.Find(&wRoles)

	if result.Error != nil {
		fmt.Printf("Failed to get welcome roles: %v\n", result.Error)
		return result.Error
	}

	for _, wRole := range wRoles {
		if wRole.UserID == nil {
			if err := s.GuildMemberRoleAdd(r.GuildID, r.Member.User.ID, wRole.RoleID); err != nil {
				fmt.Println("Failed to add role to new member:", err)
			}
			continue
		}

		if *wRole.UserID == r.Member.User.ID {
			if err := s.GuildMemberRoleAdd(r.GuildID, r.Member.User.ID, wRole.RoleID); err != nil {
				fmt.Println("Failed to add role to new member:", err)
			}
		}
	}
	return nil
}

func (f *MembersFeature) MessageReactionAddHandler(s *discordgo.Session, m *discordgo.MessageReactionAdd, ctx *bot.BotContext) error {
	user, err := s.User(m.UserID)
	if err != nil {
		ctx.Logger.Printf("Error fetching user: %v", err)
		return err
	}

	if user.Bot {
		return nil
	}

	targetEmojiID := "957421664738639872"
	channelID := m.ChannelID
	messageID := m.MessageID
	botID := ctx.ClientID
	if m.Emoji.ID != targetEmojiID {
		return nil
	}

	msg, err := s.ChannelMessage(channelID, messageID)
	if err != nil {
		ctx.Logger.Printf("Failed to get message: %v", err)
		return err
	}

	if strings.Contains(msg.Content, "[Fix]") && msg.Author.ID == botID {
		newContent := msg.Content + " "
		_, err = s.ChannelMessageEdit(channelID, messageID, newContent)
		if err != nil {
			ctx.Logger.Printf("Error editing message: %v", err)
			return err
		}
	}
	return nil
}
