package subcommands

import (
	"fmt"
	"strings"

	"github.com/bachacode/gatoc/internal/bot"
	"github.com/bachacode/gatoc/internal/database"
	"github.com/bwmarrin/discordgo"
)

var ResponseMessageAdd bot.SlashSubcommand = bot.SlashSubcommand{
	Metadata: &discordgo.ApplicationCommandOption{
		Type:        discordgo.ApplicationCommandOptionSubCommand,
		Name:        "añadir",
		Description: "Añade una nueva respuesta",
		Options: []*discordgo.ApplicationCommandOption{
			{
				Type:        discordgo.ApplicationCommandOptionString,
				Name:        "message",
				Description: "Mensaje al que responder",
				Required:    true,
			},
			{
				Type:        discordgo.ApplicationCommandOptionString,
				Name:        "response",
				Description: "Mensaje de respuesta",
				Required:    true,
			},
			{
				Type:        discordgo.ApplicationCommandOptionUser,
				Name:        "user",
				Description: "Usuario al que responder especificamente (opcional)",
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
		var message string
		var responseMsg string
		var targetUser *discordgo.User
		var content string

		// Get the "user" option (optional)
		if userOption, ok := optionMap["user"]; ok {
			userID := userOption.Value.(string)
			if i.ApplicationCommandData().Resolved != nil && i.ApplicationCommandData().Resolved.Users != nil {
				targetUser = i.ApplicationCommandData().Resolved.Users[userID]
			}
		}

		// Get the "message" option (required)
		if messageOption, ok := optionMap["message"]; ok {
			message = messageOption.Value.(string)
		} else {
			content = "Ha ocurrido un error para obtener el mensaje."
			bot.EditDeferred(s, i, content)
			return fmt.Errorf("Error responding to interaction\n")
		}

		// Get the "response" option (required)
		if responseOption, ok := optionMap["response"]; ok {
			responseMsg = responseOption.Value.(string)
		} else {
			content = "Ha ocurrido un error para obtener el mensaje de respuesta."
			bot.EditDeferred(s, i, content)
			return fmt.Errorf("Error responding to interaction\n")
		}

		response := database.ResponseMessage{
			GuildID:  guildID,
			Message:  strings.ToLower(message),
			Response: responseMsg,
		}

		content = fmt.Sprintf("Se ha añadido un nuevo mensaje de respuesta para todos")
		if targetUser != nil {
			content = fmt.Sprintf("Se ha añadido un nuevo mensaje de respuesta para `%s`", targetUser.Username)
			response.UserID = &targetUser.ID
		}

		if result := db.Create(&response); result.Error != nil {
			return fmt.Errorf("Error creating new response message\n%v", result.Error)
		}

		bot.EditDeferred(s, i, content)

		return nil
	},
}
