package arkw

import (
	"fmt"
	"strings"
	"time"

	"github.com/Pauloo27/aryzona/internal/command"
	"github.com/Pauloo27/aryzona/internal/command/parameters"
	"github.com/Pauloo27/aryzona/internal/discord/model"
	"github.com/Pauloo27/aryzona/internal/i18n"
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
	app, err := bot.s.CurrentApplication()
	if err != nil {
		return err
	}

	defaultLang := i18n.MustGetLanguage(i18n.DefaultLanguageName)
	otherLangs := make([]*i18n.Language, len(i18n.LanguagesName)-1)

	langCounter := 0

	for _, langName := range i18n.LanguagesName {
		if langName == i18n.DefaultLanguageName {
			continue
		}
		otherLangs[langCounter] = i18n.MustGetLanguage(langName)
		langCounter++
	}

	var slashCommands []api.CreateCommandData
	for key, cmd := range command.GetCommandMap() {
		if key != cmd.Name {
			continue
		}

		defaultCmdLang := i18n.GetCommandDefinition(defaultLang, cmd.Name)
		if defaultCmdLang == nil {
			logger.Fatalf("Command %s not found in default language", cmd.Name)
			break
		}

		slashCommand := api.CreateCommandData{
			Name: defaultCmdLang.Name.Str(), Description: defaultCmdLang.Description.Str(),
			NameLocalizations:        make(dc.StringLocales),
			DescriptionLocalizations: make(dc.StringLocales),
		}

		for _, lang := range otherLangs {
			cmdLang := i18n.GetCommandDefinition(lang, cmd.Name)
			if cmdLang == nil {
				logger.Fatalf("Command %s not found in language %s", cmd.Name, lang.Name)
				break
			}
			dcName := dc.Language(lang.Name.DiscordName())
			slashCommand.NameLocalizations[dcName] = cmdLang.Name.Str()
			slashCommand.DescriptionLocalizations[dcName] = cmdLang.Description.Str()
		}

		// TODO: i18n
		for _, subCmd := range cmd.SubCommands {
			subCmdOptions := []dc.CommandOptionValue{}
			for _, subCmdParam := range subCmd.Parameters {
				subCmdOptions = append(subCmdOptions, mustGetOptionValue(subCmdParam))
			}
			slashCommand.Options = append(
				slashCommand.Options,
				&dc.SubcommandOption{
					OptionName: subCmd.Name,

					Description: subCmd.Description,
					Options:     subCmdOptions,
				},
			)
		}

		// TODO: i18n
		for _, arg := range cmd.Parameters {
			slashCommand.Options = append(
				slashCommand.Options, mustGetOption(arg),
			)
		}

		slashCommands = append(slashCommands, slashCommand)
	}

	if _, err = bot.s.BulkOverwriteCommands(app.ID, slashCommands); err != nil {
		return err
	}

	bot.s.AddHandler(func(i *gateway.InteractionCreateEvent) {
		respond := func(message *model.ComplexMessage, flags dc.MessageFlags) error {
			var embeds []dc.Embed
			if len(message.Embeds) > 0 {
				embed := message.Embeds[0]
				embeds = append(embeds, buildEmbed(embed))
			}
			return bot.s.RespondInteraction(i.ID, i.Token, api.InteractionResponse{
				Type: api.MessageInteractionWithSource,
				Data: &api.InteractionResponseData{
					Content: option.NewNullableString(message.Content),
					Embeds:  &embeds,
					Flags:   flags,
				},
			})
		}

		edit := func(message *model.ComplexMessage, flags dc.MessageFlags) error {
			var embeds []dc.Embed
			if len(message.Embeds) > 0 {
				embed := message.Embeds[0]
				embeds = append(embeds, buildEmbed(embed))
			}

			var components *dc.ContainerComponents
			if len(message.Components) > 0 {
				rawComponents := buildComponents(message.Components)
				row := dc.ActionRowComponent(rawComponents)
				components = dc.ComponentsPtr(&row)
			}

			_, err := bot.s.EditInteractionResponse(i.AppID, i.Token, api.EditInteractionResponseData{
				Content:    option.NewNullableString(message.Content),
				Embeds:     &embeds,
				Components: components,
			})
			return err
		}

		switch data := i.Data.(type) {
		case dc.ComponentInteraction:
			newMessage := command.HandleInteraction(string(data.ID()), i.Sender().ID.String())
			if newMessage == nil {
				err := bot.s.RespondInteraction(i.ID, i.Token, api.InteractionResponse{
					Type: api.DeferredMessageUpdate,
				})
				if err != nil {
					logger.Error(err)
				}
				return
			}
			var embeds []dc.Embed

			for _, embed := range newMessage.Embeds {
				embeds = append(embeds, buildEmbed(embed))
			}

			var embedsPtr *[]dc.Embed
			if len(embeds) > 0 {
				embedsPtr = &embeds
			}

			components := buildComponents(newMessage.Components)
			row := dc.ActionRowComponent(components)

			var componentsPtr *dc.ContainerComponents
			if len(components) > 0 {
				componentsPtr = &dc.ContainerComponents{&row}
			}

			var content option.NullableString
			if newMessage.Content != "" {
				content = option.NewNullableString(newMessage.Content)
			}

			err := bot.s.RespondInteraction(i.ID, i.Token, api.InteractionResponse{
				Type: api.UpdateMessage,
				Data: &api.InteractionResponseData{
					Content:    content,
					Embeds:     embedsPtr,
					Components: componentsPtr,
				},
			})
			if err != nil {
				logger.Error(err)
			}
		case *dc.CommandInteraction:
			startTime := time.Now()
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

			cType := model.ChannelTypeGuild
			if i.GuildID.String() == "" {
				cType = model.ChannelTypeDirect
			}

			var flags dc.MessageFlags
			if cmd.Ephemeral {
				flags = 64
			}

			adapter := command.Adapter{
				AuthorID: i.Sender().ID.String(),
				GuildID:  i.GuildID.String(),
				DeferResponse: func() error {
					return bot.s.RespondInteraction(
						i.ID,
						i.Token,
						api.InteractionResponse{
							Type: api.DeferredMessageInteractionWithSource,
						},
					)
				},
				ReplyComplex: func(ctx *command.CommandContext, message *model.ComplexMessage) error {
					if ctx.Command.Deferred {
						return edit(message, flags)
					}
					return respond(message, flags)
				},
				Reply: func(ctx *command.CommandContext, message string) error {
					if ctx.Command.Deferred {
						return edit(&model.ComplexMessage{Content: message}, flags)
					}
					return respond(&model.ComplexMessage{Content: message}, flags)
				},
				ReplyEmbed: func(ctx *command.CommandContext, embed *model.Embed) error {
					if ctx.Command.Deferred {
						return edit(&model.ComplexMessage{Embeds: []*model.Embed{embed}}, flags)
					}
					return respond(&model.ComplexMessage{Embeds: []*model.Embed{embed}}, flags)
				},
				EditComplex: func(ctx *command.CommandContext, message *model.ComplexMessage) error {
					return edit(message, flags)
				},
				Edit: func(ctx *command.CommandContext, message string) error {
					return edit(&model.ComplexMessage{Content: message}, flags)
				},
				EditEmbed: func(ctx *command.CommandContext, embed *model.Embed) error {
					return edit(&model.ComplexMessage{Embeds: []*model.Embed{embed}}, flags)
				},
			}

			langName := i18n.FindLanguageName(strings.Replace(string(i.Locale), "-", "_", 1))

			command.HandleCommand(
				data.Name, args, langName, startTime, &adapter, bot, command.CommandTriggerSlash,
				buildChannel(i.ChannelID.String(), buildGuild(i.GuildID.String()), cType),
			)
		}
	})

	return nil
}
