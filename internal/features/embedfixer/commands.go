package embedfixer

import (
	"fmt"
	"strings"

	"github.com/bachacode/gatoc/internal/bot"
	"github.com/bwmarrin/discordgo"
)

var embedFixerDefaultPermissions int64 = discordgo.PermissionManageServer

var embedFixerConfigCommand = bot.SlashCommand{
	Metadata: &discordgo.ApplicationCommand{
		Name:                     "embedfixer",
		Description:              "Configure social embed fixer domains",
		DefaultMemberPermissions: &embedFixerDefaultPermissions,
		Options: []*discordgo.ApplicationCommandOption{
			{
				Type:        discordgo.ApplicationCommandOptionSubCommand,
				Name:        "show",
				Description: "Show current platform domain config",
			},
			{
				Type:        discordgo.ApplicationCommandOptionSubCommand,
				Name:        "nofix",
				Description: "Explain how to bypass fixing with #nofix",
			},
			{
				Type:        discordgo.ApplicationCommandOptionSubCommand,
				Name:        "set",
				Description: "Set custom domain for a supported platform",
				Options: []*discordgo.ApplicationCommandOption{
					{
						Type:        discordgo.ApplicationCommandOptionString,
						Name:        "platform",
						Description: "Supported platform (twitter, reddit, instagram)",
						Required:    true,
					},
					{
						Type:        discordgo.ApplicationCommandOptionString,
						Name:        "domain",
						Description: "Replacement domain host (example: fxtwitter.com)",
						Required:    true,
					},
				},
			},
			{
				Type:        discordgo.ApplicationCommandOptionSubCommand,
				Name:        "reset",
				Description: "Reset platform custom domain to default",
				Options: []*discordgo.ApplicationCommandOption{
					{
						Type:        discordgo.ApplicationCommandOptionString,
						Name:        "platform",
						Description: "Supported platform (twitter, reddit, instagram)",
						Required:    true,
					},
				},
			},
		},
	},
	Handler: func(s *discordgo.Session, i *discordgo.InteractionCreate, ctx *bot.BotContext) error {
		options := i.ApplicationCommandData().Options
		if len(options) == 0 {
			bot.GetInteractionFailedResponse(s, i, "Subcommand is required")
			return fmt.Errorf("embedfixer subcommand is required")
		}

		subcommand := options[0]
		switch subcommand.Name {
		case "show":
			return handleConfigShow(s, i, ctx)
		case "nofix":
			return handleNoFixInfo(s, i)
		case "set":
			return handleConfigSet(s, i, ctx, subcommand)
		case "reset":
			return handleConfigReset(s, i, ctx, subcommand)
		default:
			bot.GetInteractionFailedResponse(s, i, "Subcommand does not exist")
			return fmt.Errorf("unknown embedfixer subcommand: %s", subcommand.Name)
		}
	},
}

func buildPlatformChoices() string {
	keys := make([]string, 0, len(supportedPlatforms))
	for _, platform := range supportedPlatforms {
		keys = append(keys, platform.Key)
	}

	return strings.Join(keys, ", ")
}
