package subcommands

import (
	"errors"
	"fmt"

	"github.com/bachacode/gatoc/internal/bot"
	"github.com/bachacode/gatoc/internal/database"
	"github.com/bwmarrin/discordgo"
	"gorm.io/gorm"
)

var WelcomeRoleAdd bot.SlashSubcommand = bot.SlashSubcommand{
	Metadata: &discordgo.ApplicationCommandOption{
		Type:        discordgo.ApplicationCommandOptionSubCommand,
		Name:        "añadir",
		Description: "Añade un nuevo rol de bienvenida",
		Options: []*discordgo.ApplicationCommandOption{
			{
				Type:        discordgo.ApplicationCommandOptionRole, // This is the key!
				Name:        "role",
				Description: "Rol a añadir",
				Required:    true,
			},
			{
				Type:        discordgo.ApplicationCommandOptionUser,
				Name:        "user",
				Description: "Usuario para añadir rol especifico (opcional)",
				Required:    false,
			},
		},
	},
	Handler: func(s *discordgo.Session, i *discordgo.InteractionCreate, ctx *bot.BotContext) error {
		options := i.ApplicationCommandData().Options

		optionMap := make(map[string]*discordgo.ApplicationCommandInteractionDataOption, len(options))
		for _, opt := range options[0].Options {
			optionMap[opt.Name] = opt
		}

		if err := bot.DeferReply(s, i); err != nil {
			bot.GetInteractionFailedResponse(s, i, "")
			return err
		}

		db := ctx.DB
		guildID := ctx.GuildID
		var selectedRole *discordgo.Role
		var targetUser *discordgo.User
		var content string
		if roleOption, ok := optionMap["role"]; ok {
			roleID := roleOption.Value.(string)
			if i.ApplicationCommandData().Resolved != nil && i.ApplicationCommandData().Resolved.Roles != nil {
				selectedRole = i.ApplicationCommandData().Resolved.Roles[roleID]
			}
		} else {
			content = "Ha ocurrido un error para obtener el rol"
			bot.EditDeferred(s, i, content)
			return fmt.Errorf("Error responding to interaction\n")
		}

		// Get the "user" option (optional)
		if userOption, ok := optionMap["user"]; ok {
			userID := userOption.Value.(string)
			if i.ApplicationCommandData().Resolved != nil && i.ApplicationCommandData().Resolved.Users != nil {
				targetUser = i.ApplicationCommandData().Resolved.Users[userID]
			}
		}

		welcomeRole := database.WelcomeRole{
			GuildID: guildID,
			RoleID:  selectedRole.ID,
		}

		content = fmt.Sprintf("El rol `%s` será asignado a los nuevos miembros", selectedRole.Name)
		if targetUser != nil {
			content = fmt.Sprintf("El rol `%s` será asignado al usuario `%s`", selectedRole.Name, targetUser.Username)
			welcomeRole.UserID = &targetUser.ID
		}

		if result := db.Create(&welcomeRole); result.Error != nil {
			if errors.Is(result.Error, gorm.ErrDuplicatedKey) {
				content = "Ese rol ya esta registrado en los roles de bienvenidas para ese/esos usuario/s"
			} else {
				content = "Ha ocurrido un error agregando el rol a los roles de bienvenida"
			}
			bot.EditDeferred(s, i, content)
			return fmt.Errorf("Error creating welcome role: %s\n%v", selectedRole.Name, result.Error)
		}

		bot.EditDeferred(s, i, content)

		return nil
	},
}
