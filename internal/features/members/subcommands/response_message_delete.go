package subcommands

import (
	"fmt"

	"github.com/bachacode/gatoc/internal/bot"
	"github.com/bachacode/gatoc/internal/database"
	"github.com/bwmarrin/discordgo"
)

var ResponseMessageDelete bot.SlashSubcommand = bot.SlashSubcommand{
	Metadata: &discordgo.ApplicationCommandOption{
		Type:        discordgo.ApplicationCommandOptionSubCommand,
		Name:        "eliminar",
		Description: "Elimina una respuesta",
		Options: []*discordgo.ApplicationCommandOption{
			{
				Type:        discordgo.ApplicationCommandOptionString,
				Name:        "id",
				Description: "Respuesta a eliminar",
				Required:    true,
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
		var responseID string
		var content string
		if idOption, ok := optionMap["id"]; ok {
			responseID = idOption.Value.(string)
		} else {
			content = "Ha ocurrido un error para obtener la respuesta."
			bot.EditDeferred(s, i, content)
			return fmt.Errorf("Error responding to interaction\n")
		}

		if result := db.Delete(&database.ResponseMessage{}, responseID); result.Error != nil {
			content = "Ha ocurrido un error al eliminar el mensaje de respuesta"
			bot.EditDeferred(s, i, content)
			return fmt.Errorf("Error deleting welcome role: %s\n%v", responseID, result.Error)
		} else if result.RowsAffected < 1 {
			content = "El ID introducido no existe en los mensajes de respuesta"
			bot.EditDeferred(s, i, content)
			return nil
		}

		content = fmt.Sprintf("El mensaje de ID `%s` ha sido eliminado de los mensajes de respuesta", responseID)
		bot.EditDeferred(s, i, content)

		return nil
	},
}
