package arkw

import (
	"fmt"

	"github.com/Pauloo27/aryzona/internal/command"
	"github.com/Pauloo27/aryzona/internal/command/parameters"
	"github.com/Pauloo27/aryzona/internal/discord/model"
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
	switch arg.Type.BaseType {
	case parameters.TypeString:
		return &dc.StringOption{
			OptionName: arg.Name, Description: arg.Description, Required: arg.Required,
			Choices: mustGetStringChoises(arg),
		}
	case parameters.TypeBool:
		return &dc.BooleanOption{
			OptionName: arg.Name, Description: arg.Description, Required: arg.Required,
		}
	case parameters.TypeInt:
		return &dc.IntegerOption{
			OptionName: arg.Name, Description: arg.Description, Required: arg.Required,
			Choices: mustGetIntegerChoises(arg),
		}
	default:
		logger.Fatalf("Cannot find discord type for %s", arg.Type.BaseType.Name)
	}
	return nil
}

func mustGetOptionValue(arg *command.CommandParameter) dc.CommandOptionValue {
	switch arg.Type.BaseType {
	case parameters.TypeString:
		return &dc.StringOption{
			OptionName: arg.Name, Description: arg.Description, Required: arg.Required,
			Choices: mustGetStringChoises(arg),
		}
	case parameters.TypeBool:
		return &dc.BooleanOption{
			OptionName: arg.Name, Description: arg.Description, Required: arg.Required,
		}
	case parameters.TypeInt:
		return &dc.IntegerOption{
			OptionName: arg.Name, Description: arg.Description, Required: arg.Required,
			Choices: mustGetIntegerChoises(arg),
		}
	default:
		logger.Fatalf("Cannot find discord type for %s", arg.Type.BaseType.Name)
	}
	return nil
}

func registerCommands(bot ArkwBot) error {
	s := bot.d.s

	app, err := s.CurrentApplication()
	if err != nil {
		return err
	}

	var slashCommands []api.CreateCommandData
	for key, cmd := range command.GetCommandMap() {
		if key != cmd.Name {
			continue
		}

		slashCommand := api.CreateCommandData{
			Name: cmd.Name, Description: cmd.Description,
		}

		for _, subCmd := range cmd.SubCommands {
			subCmdOptions := []dc.CommandOptionValue{}
			for _, subCmdParam := range subCmd.Parameters {
				subCmdOptions = append(subCmdOptions, mustGetOptionValue(subCmdParam))
			}
			slashCommand.Options = append(
				slashCommand.Options,
				&dc.SubcommandOption{
					OptionName:  subCmd.Name,
					Description: subCmd.Description,
					Options:     subCmdOptions,
				},
			)
		}

		for _, arg := range cmd.Parameters {
			slashCommand.Options = append(
				slashCommand.Options, mustGetOption(arg),
			)
		}

		slashCommands = append(slashCommands, slashCommand)
	}

	if _, err = s.BulkOverwriteCommands(app.ID, slashCommands); err != nil {
		return err
	}

	s.AddHandler(func(i *gateway.InteractionCreateEvent) {
		respond := func(message *model.ComplexMessage) error {
			var embeds []dc.Embed
			if len(message.Embeds) > 0 {
				embed := message.Embeds[0]
				embeds = append(embeds, buildEmbed(embed))
			}
			return s.RespondInteraction(i.ID, i.Token, api.InteractionResponse{
				Type: api.MessageInteractionWithSource,
				Data: &api.InteractionResponseData{
					Content: option.NewNullableString(message.Content),
					Embeds:  &embeds,
				},
			})
		}

		edit := func(message *model.ComplexMessage) error {
			var embeds []dc.Embed
			if len(message.Embeds) > 0 {
				embed := message.Embeds[0]
				embeds = append(embeds, buildEmbed(embed))
			}
			components := buildComponents(message.Components)
			row := dc.ActionRowComponent(components)

			_, err := s.EditInteractionResponse(i.AppID, i.Token, api.EditInteractionResponseData{
				Content:    option.NewNullableString(message.Content),
				Embeds:     &embeds,
				Components: dc.ComponentsPtr(&row),
			})
			return err
		}

		if data, ok := i.Data.(*dc.CommandInteraction); ok {
			cmd, ok := command.GetCommandMap()[data.Name]
			if !ok {
				logger.Error("Invalid slash command interaction received:", data.Name)
				return
			}

			var args []string
			for i, option := range data.Options {
				if i == 0 && cmd.SubCommands != nil {
					args = append(args, option.Name)
					for _, subCommandOption := range option.Options {
						args = append(args, fmt.Sprintf("%v", subCommandOption.Value))
					}
					break
				}

				args = append(args, option.String())
			}

			var member model.Member

			if !i.GuildID.IsNull() {
				m, err := bot.GetMember(i.GuildID.String(), i.SenderID().String())
				if err != nil {
					return
				}
				member = m
			}

			adapter := command.Adapter{
				Member:   member,
				AuthorID: i.Sender().ID.String(),
				GuildID:  i.GuildID.String(),
				DeferResponse: func() error {
					return s.RespondInteraction(
						i.ID,
						i.Token,
						api.InteractionResponse{
							Type: api.DeferredMessageInteractionWithSource,
						},
					)
				},
				ReplyComplex: func(ctx *command.CommandContext, message *model.ComplexMessage) error {
					if ctx.Command.Deferred {
						return edit(message)
					}
					return respond(message)
				},
				Reply: func(ctx *command.CommandContext, message string) error {
					if ctx.Command.Deferred {
						return edit(&model.ComplexMessage{Content: message})
					}
					return respond(&model.ComplexMessage{Content: message})
				},
				ReplyEmbed: func(ctx *command.CommandContext, embed *model.Embed) error {
					if ctx.Command.Deferred {
						return edit(&model.ComplexMessage{Embeds: []*model.Embed{embed}})
					}
					return respond(&model.ComplexMessage{Embeds: []*model.Embed{embed}})
				},
			}

			cType := model.ChannelTypeGuild
			if i.GuildID.String() == "" {
				cType = model.ChannelTypeDirect
			}

			command.HandleCommand(
				data.Name, args, &adapter, bot, command.CommandTriggerSlash,
				buildChannel(i.ChannelID.String(), buildGuild(i.GuildID.String()), cType),
			)
		}
	})

	return nil
}
