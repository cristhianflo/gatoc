package subcommands

import (
	"fmt"
	"time"

	"github.com/bachacode/gatoc/internal/bot"
	"github.com/bachacode/gatoc/internal/database"
	"github.com/bwmarrin/discordgo"
)

var WelcomeRoleList bot.SlashSubcommand = bot.SlashSubcommand{
	Metadata: &discordgo.ApplicationCommandOption{
		Type:        discordgo.ApplicationCommandOptionSubCommand,
		Name:        "lista",
		Description: "Ver todos los roles de bienvenida registrados",
	},
	Handler: func(s *discordgo.Session, i *discordgo.InteractionCreate, ctx *bot.BotContext) error {

		if err := bot.DeferReply(s, i); err != nil {
			bot.GetInteractionFailedResponse(s, i, "")
			return err
		}

		db := ctx.DB
		content := ""

		var wRoles []database.WelcomeRole
		result := db.Find(&wRoles)

		if result.Error != nil {
			content = "Ha ocurrido un error para obtener la lista de roles"
			bot.EditDeferred(s, i, content)
			return fmt.Errorf("Error responding to interaction: %v\n", result.Error)
		}

		fields := []*discordgo.MessageEmbedField{
			{
				Name:   "ID",
				Value:  "",
				Inline: true,
			},
			{
				Name:   "Rol",
				Value:  "",
				Inline: true,
			},
			{
				Name:   "Otorgado a",
				Value:  "",
				Inline: true,
			},
		}

		for _, wRole := range wRoles {
			role, err := s.State.Role(i.GuildID, wRole.RoleID)
			if err != nil {
				content = "Ha ocurrido un error obteniendo un rol"
				bot.EditDeferred(s, i, content)
				return fmt.Errorf("Error getting a role: %s from the guild: %s\n%v", wRole.RoleID, i.GuildID, err)
			}

			fields[0].Value += fmt.Sprintf("%d\n", wRole.ID)
			fields[1].Value += role.Name + "\n"
			if wRole.UserID == nil {
				fields[2].Value += "Todos\n"
				continue
			}

			user, err := s.User(*wRole.UserID)

			if err != nil {
				content = "Ha ocurrido un error obteniendo un usuario"
				bot.EditDeferred(s, i, content)
				return fmt.Errorf("Error getting an user: %s from the guild: %s\n%v", *wRole.UserID, i.GuildID, err)
			}

			if user.GlobalName != "" {
				fields[2].Value += user.GlobalName + "\n"
			} else {
				fields[2].Value += user.Username + "\n"
			}
		}

		embed := discordgo.MessageEmbed{
			Title:       "Roles de bienvenida",
			Description: "Roles dados a los miembros al unirse al servidor.",
			Fields:      fields,
			Color:       0xFFFFFF,
			Footer: &discordgo.MessageEmbedFooter{
				Text: fmt.Sprintf("Generated at %s", time.Now().Format("2006-01-02 15:04:05")),
			},
		}

		// Send the embed to the channel
		_, err := s.InteractionResponseEdit(
			i.Interaction,
			&discordgo.WebhookEdit{
				Embeds: &[]*discordgo.MessageEmbed{&embed},
			},
		)
		if err != nil {
			return fmt.Errorf("Error sending embed: %v\n", err)
		}

		return nil
	},
}
