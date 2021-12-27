package arkw

import (
	"fmt"

	"github.com/Pauloo27/aryzona/internal/command"
	"github.com/Pauloo27/aryzona/internal/command/parameters"
	"github.com/Pauloo27/aryzona/internal/discord"
	"github.com/Pauloo27/logger"

	"github.com/diamondburned/arikawa/v3/api"
	dc "github.com/diamondburned/arikawa/v3/discord"
	"github.com/diamondburned/arikawa/v3/gateway"
	"github.com/diamondburned/arikawa/v3/utils/json/option"
)

func mustGetIntegerChoises(arg *command.CommandParameter) (choises []dc.IntegerChoice) {
	for _, value := range arg.GetValidValues() {
		choises = append(choises, dc.IntegerChoice{
			Name:  fmt.Sprintf("%v", value),
			Value: value.(int),
		})
	}
	return
}

func mustGetStringChoises(arg *command.CommandParameter) (choises []dc.StringChoice) {
	for _, value := range arg.GetValidValues() {
		choises = append(choises, dc.StringChoice{
			Name:  fmt.Sprintf("%v", value),
			Value: value.(string),
		})
	}
	return
}

func mustGetOption(arg *command.CommandParameter) dc.CommandOption {
	switch arg.Type {
	case parameters.ParameterString:
		return &dc.StringOption{
			OptionName: arg.Name, Description: arg.Description, Required: arg.Required,
			Choices: mustGetStringChoises(arg),
		}
	case parameters.ParameterText:
		return &dc.StringOption{
			OptionName: arg.Name, Description: arg.Description, Required: arg.Required,
			Choices: mustGetStringChoises(arg),
		}
	case parameters.ParameterBool:
		return &dc.BooleanOption{
			OptionName: arg.Name, Description: arg.Description, Required: arg.Required,
		}
	case parameters.ParameterInt:
		return &dc.IntegerOption{
			OptionName: arg.Name, Description: arg.Description, Required: arg.Required,
			Choices: mustGetIntegerChoises(arg),
		}
	default:
		return &dc.UnknownCommandOption{}
	}
}

func registerCommands(bot ArkwBot) error {
	session := bot.d.s

	app, err := session.CurrentApplication()
	if err != nil {
		return err
	}

	var slashCommands []dc.Command
	for key, cmd := range command.GetCommandMap() {
		if key != cmd.Name {
			continue
		}

		slashCommand := dc.Command{
			Name: cmd.Name, Description: cmd.Description,
		}

		for _, arg := range cmd.Parameters {
			slashCommand.Options = append(
				slashCommand.Options, mustGetOption(arg),
			)
		}

		slashCommands = append(slashCommands, slashCommand)
	}

	if _, err = session.BulkOverwriteCommands(app.ID, slashCommands); err != nil {
		return err
	}

	session.AddHandler(func(i *gateway.InteractionCreateEvent) {
		if data, ok := i.Data.(*dc.CommandInteraction); ok {
			_, ok := command.GetCommandMap()[data.Name]
			if !ok {
				logger.Error("Invalid slash command interaction received:", data.Name)
				return
			}

			var args []string
			for _, option := range data.Options {
				args = append(args, option.String())
			}

			event := command.Event{
				AuthorID: i.Sender().ID.String(),
				GuildID:  i.GuildID.String(),
				Reply: func(message string) error {
					return session.RespondInteraction(i.ID, i.Token, api.InteractionResponse{
						Type: api.MessageInteractionWithSource,
						Data: &api.InteractionResponseData{
							Content: option.NewNullableString(message),
						},
					})
				},
				ReplyEmbed: func(e *discord.Embed) error {
					return session.RespondInteraction(i.ID, i.Token, api.InteractionResponse{
						Type: api.MessageInteractionWithSource,
						Data: &api.InteractionResponseData{
							Embeds: &[]dc.Embed{buildEmbed(e)},
						},
					})
				},
			}

			command.HandleCommand(data.Name, args, &event, bot)
		}
	})

	return nil
}
