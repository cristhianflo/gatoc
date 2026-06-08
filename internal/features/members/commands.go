package members

import (
	"fmt"

	"github.com/bachacode/gatoc/internal/bot"
	"github.com/bachacode/gatoc/internal/features/members/subcommands"
	"github.com/bwmarrin/discordgo"
)

var defaultMemberPermissions int64 = discordgo.PermissionManageServer

var welcomeRole bot.SlashCommand = bot.SlashCommand{
	Metadata: &discordgo.ApplicationCommand{
		Name:                     "rol-de-bienvenida",
		Description:              "Gestiona los roles otorgados a los nuevos miembros",
		DefaultMemberPermissions: &defaultMemberPermissions,
		Options: []*discordgo.ApplicationCommandOption{
			subcommands.WelcomeRoleAdd.Metadata,
			subcommands.WelcomeRoleList.Metadata,
			subcommands.WelcomeRoleDelete.Metadata,
		},
	},
	Handler: func(s *discordgo.Session, i *discordgo.InteractionCreate, ctx *bot.BotContext) error {
		options := i.ApplicationCommandData().Options

		switch options[0].Name {
		case "añadir":
			return subcommands.WelcomeRoleAdd.Handler(s, i, ctx)
		case "lista":
			return subcommands.WelcomeRoleList.Handler(s, i, ctx)
		case "eliminar":
			return subcommands.WelcomeRoleDelete.Handler(s, i, ctx)
		default:
			bot.GetInteractionFailedResponse(s, i, "El subcomando llamado no existe.")
			return fmt.Errorf("Subcommand doesn't exist\n")
		}
	},
}

var responseMessage bot.SlashCommand = bot.SlashCommand{
	Metadata: &discordgo.ApplicationCommand{
		Name:                     "respuestas",
		Description:              "Gestiona las respuestas a mensajes especificos",
		DefaultMemberPermissions: &defaultMemberPermissions,
		Options: []*discordgo.ApplicationCommandOption{
			subcommands.ResponseMessageAdd.Metadata,
			subcommands.ResponseMessageList.Metadata,
			subcommands.ResponseMessageDelete.Metadata,
		},
	},
	Handler: func(s *discordgo.Session, i *discordgo.InteractionCreate, ctx *bot.BotContext) error {
		options := i.ApplicationCommandData().Options

		switch options[0].Name {
		case "añadir":
			return subcommands.ResponseMessageAdd.Handler(s, i, ctx)
		case "lista":
			return subcommands.ResponseMessageList.Handler(s, i, ctx)
		case "eliminar":
			return subcommands.ResponseMessageDelete.Handler(s, i, ctx)
		default:
			bot.GetInteractionFailedResponse(s, i, "El subcomando llamado no existe.")
			return fmt.Errorf("Subcommand doesn't exist\n")
		}
	},
}
