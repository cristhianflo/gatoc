package members

import "github.com/bwmarrin/discordgo"

type LeaveMessage struct {
	embed    *discordgo.MessageEmbed
	filename string
}
