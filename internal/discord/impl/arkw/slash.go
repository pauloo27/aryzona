package arkw

import (
	"fmt"
	"strings"
	"time"

	"github.com/pauloo27/aryzona/internal/command"
	"github.com/pauloo27/aryzona/internal/command/parameters"
	"github.com/pauloo27/aryzona/internal/data/services"
	"github.com/pauloo27/aryzona/internal/discord/model"
	"github.com/pauloo27/aryzona/internal/i18n"
	"github.com/pauloo27/logger"

	"github.com/diamondburned/arikawa/v3/api"
	dc "github.com/diamondburned/arikawa/v3/discord"
	"github.com/diamondburned/arikawa/v3/gateway"
	"github.com/diamondburned/arikawa/v3/utils/json/option"
)

type languageContext struct {
	defaultLangCmd *i18n.CommandDefinition
	otherLangs     []*i18n.Language
	otherLangsCmd  []*i18n.CommandDefinition
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

		defaultCmdLang := i18n.MustGetCommandDefinition(defaultLang, cmd.Name)
		if defaultCmdLang == nil {
			logger.Fatalf("Command %s not found in default language", cmd.Name)
			break
		}

		otherCmdLangs := make([]*i18n.CommandDefinition, len(otherLangs))
		for i, lang := range otherLangs {
			otherCmdLangs[i] = i18n.MustGetCommandDefinition(lang, cmd.Name)
			if otherCmdLangs[i] == nil {
				logger.Fatalf("Command %s not found in language %s", cmd.Name, lang.Name)
				return nil
			}
		}

		slashCommand := api.CreateCommandData{
			Name:                     defaultCmdLang.Name.Str(),
			Description:              defaultCmdLang.Description.Str(),
			NameLocalizations:        make(dc.StringLocales),
			DescriptionLocalizations: make(dc.StringLocales),
		}

		for i, lang := range otherLangs {
			cmdLang := otherCmdLangs[i]
			dcName := dc.Language(lang.Name.DiscordName())
			slashCommand.NameLocalizations[dcName] = cmdLang.Name.Str()
			slashCommand.DescriptionLocalizations[dcName] = cmdLang.Description.Str()
		}

		for i, arg := range cmd.Parameters {
			slashCommand.Options = append(
				slashCommand.Options, mustGetOption(
					arg, i, languageContext{defaultCmdLang, otherLangs, otherCmdLangs},
				),
			)
		}

		for i, subCmd := range cmd.SubCommands {
			subCmdOptions := make([]dc.CommandOptionValue, len(subCmd.Parameters))
			for j, subCmdParam := range subCmd.Parameters {
				subCmdOptions[j] = mustGetOptionValue(
					subCmdParam, i, j, languageContext{defaultCmdLang, otherLangs, otherCmdLangs},
				)
			}

			if len(defaultCmdLang.SubCommands) <= i {
				logger.Fatalf("Subcommand %s not found in default language", subCmd.Name)
				return nil
			}

			defaultSubCmdName := defaultCmdLang.SubCommands[i].Name.Str()
			defaultSubCmdDescription := defaultCmdLang.SubCommands[i].Description.Str()

			localizedSubCmdNames := make(dc.StringLocales)
			localizedSubCmdDescriptions := make(dc.StringLocales)

			for j, l := range otherCmdLangs {
				langName := dc.Language(otherLangs[j].Name.DiscordName())

				if len(l.SubCommands) <= i {
					logger.Fatalf(
						"Cannot find sub command %d for command %s in the language %s",
						i, l.Name.Str(), langName,
					)
					return nil
				}
				localizedSubCmdNames[langName] = l.SubCommands[i].Name.Str()
				localizedSubCmdDescriptions[langName] = l.SubCommands[i].Description.Str()
			}

			slashCommand.Options = append(
				slashCommand.Options,
				&dc.SubcommandOption{
					OptionName:               defaultSubCmdName,
					Description:              defaultSubCmdDescription,
					OptionNameLocalizations:  localizedSubCmdNames,
					DescriptionLocalizations: localizedSubCmdDescriptions,
					Options:                  subCmdOptions,
				},
			)
		}

		slashCommands = append(slashCommands, slashCommand)
	}

	if _, err = bot.s.BulkOverwriteCommands(app.ID, slashCommands); err != nil {
		return err
	}

	bot.s.AddHandler(func(i *gateway.InteractionCreateEvent) {
		eventTime := time.Now()

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

			rows := buildRows(message.ComponentRows)

			_, err := bot.s.EditInteractionResponse(i.AppID, i.Token,
				api.EditInteractionResponseData{
					Content:    option.NewNullableString(message.Content),
					Embeds:     &embeds,
					Components: &rows,
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

			rows := buildRows(newMessage.ComponentRows)

			var content option.NullableString
			if newMessage.Content != "" {
				content = option.NewNullableString(newMessage.Content)
			}

			err := bot.s.RespondInteraction(i.ID, i.Token, api.InteractionResponse{
				Type: api.UpdateMessage,
				Data: &api.InteractionResponseData{
					Content:    content,
					Embeds:     embedsPtr,
					Components: &rows,
				},
			})
			if err != nil {
				logger.Error(err)
			}
		case *dc.CommandInteraction:
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

			langName := i18n.FindLanguageName(
				strings.Replace(string(i.Locale), "-", "_", 1),
			)

			go func() {
				err = services.User.SetLastSlashCommandLocale(
					i.Sender().ID.String(), langName,
				)
				if err != nil {
					logger.Error(err)
				}
			}()

			channel := buildChannel(i.ChannelID.String(), buildGuild(i.GuildID.String()), cType)

			trigger := command.TriggerEvent{
				EventTime:        eventTime,
				Channel:          channel,
				Type:             command.CommandTriggerSlash,
				PreferedLanguage: langName,
				AuthorID:         i.Sender().ID.String(),
				MessageID:        i.ID.String(),
				GuildID:          i.GuildID.String(),
				DeferResponse: func() error {
					return bot.s.RespondInteraction(
						i.ID,
						i.Token,
						api.InteractionResponse{
							Type: api.DeferredMessageInteractionWithSource,
						},
					)
				},
				Reply: func(ctx *command.Context, message *model.ComplexMessage) error {
					if ctx.Command.Deferred {
						return edit(message, flags)
					}
					return respond(message, flags)
				},
				Edit: func(ctx *command.Context, message *model.ComplexMessage) error {
					return edit(message, flags)
				},
			}

			command.HandleCommand(data.Name, args, cmd, bot, &trigger)
		}
	})

	return nil
}

func mustGetIntegerChoises(arg *command.Parameter) (choises []dc.IntegerChoice) {
	for _, value := range arg.GetValidValues() {
		choises = append(choises, dc.IntegerChoice{
			Name:  fmt.Sprintf("%v", value),
			Value: value.(int),
		})
	}
	return
}

func mustGetStringChoises(arg *command.Parameter) (choises []dc.StringChoice) {
	for _, value := range arg.GetValidValues() {
		choises = append(choises, dc.StringChoice{
			Name:  fmt.Sprintf("%v", value),
			Value: value.(string),
		})
	}
	return
}

func mustGetOption(param *command.Parameter, i int, lang languageContext) dc.CommandOption {
	if len(lang.defaultLangCmd.Parameters) <= i {
		logger.Fatalf("Cannot find parameter %d for command %s in the default language", i, lang.defaultLangCmd.Name)
		return nil
	}

	defaultParamName := lang.defaultLangCmd.Parameters[i].Name.Str()
	defaultParamDescription := lang.defaultLangCmd.Parameters[i].Description.Str()

	localizedParamName := make(dc.StringLocales)
	localizedParamDescription := make(dc.StringLocales)

	for j, l := range lang.otherLangsCmd {
		langName := dc.Language(lang.otherLangs[j].Name.DiscordName())

		if len(l.Parameters) <= i {
			logger.Fatalf("Cannot find parameter %d for command %s in the language %s", i, l.Name.Str(), langName)
			return nil
		}
		localizedParamName[langName] = l.Parameters[i].Name.Str()
		localizedParamDescription[langName] = l.Parameters[i].Description.Str()
	}

	switch param.Type.BaseType {
	case parameters.TypeString:
		return &dc.StringOption{
			OptionName:               defaultParamName,
			Description:              defaultParamDescription,
			OptionNameLocalizations:  localizedParamName,
			DescriptionLocalizations: localizedParamDescription,
			Required:                 param.Required,
			Choices:                  mustGetStringChoises(param),
		}
	case parameters.TypeBool:
		return &dc.BooleanOption{
			OptionName:               defaultParamName,
			Description:              defaultParamDescription,
			OptionNameLocalizations:  localizedParamName,
			DescriptionLocalizations: localizedParamDescription,
			Required:                 param.Required,
		}
	case parameters.TypeInt:
		return &dc.IntegerOption{
			OptionName:               defaultParamName,
			Description:              defaultParamDescription,
			OptionNameLocalizations:  localizedParamName,
			DescriptionLocalizations: localizedParamDescription,
			Required:                 param.Required,
			Choices:                  mustGetIntegerChoises(param),
		}
	default:
		logger.Fatalf("Cannot find discord type for %s", param.Type.BaseType.Name)
	}
	return nil
}

func mustGetOptionValue(param *command.Parameter, subCmdIdx, paramIdx int, lang languageContext) dc.CommandOptionValue {
	if len(lang.defaultLangCmd.SubCommands[subCmdIdx].Parameters) <= paramIdx {
		logger.Fatalf(
			"Cannot find parameter %d for command %s in the default language",
			paramIdx,
			lang.defaultLangCmd.Name,
		)
		return nil
	}

	defaultParamName := lang.defaultLangCmd.SubCommands[subCmdIdx].Parameters[paramIdx].Name.Str()
	defaultParamDescription := lang.defaultLangCmd.SubCommands[subCmdIdx].Parameters[paramIdx].Description.Str()

	localizedParamName := make(dc.StringLocales)
	localizedParamDescription := make(dc.StringLocales)

	for j, l := range lang.otherLangsCmd {
		langName := dc.Language(lang.otherLangs[j].Name.DiscordName())

		if len(l.SubCommands[subCmdIdx].Parameters) <= paramIdx {
			logger.Fatalf(
				"Cannot find parameter %d for command %s in the language %s",
				paramIdx,
				l.Name.Str(),
				langName,
			)
			return nil
		}
		localizedParamName[langName] = l.SubCommands[subCmdIdx].Parameters[paramIdx].Name.Str()
		localizedParamDescription[langName] = l.SubCommands[subCmdIdx].Parameters[paramIdx].Description.Str()
	}

	switch param.Type.BaseType {
	case parameters.TypeString:
		return &dc.StringOption{
			OptionName:               defaultParamName,
			Description:              defaultParamDescription,
			OptionNameLocalizations:  localizedParamName,
			DescriptionLocalizations: localizedParamDescription,
			Required:                 param.Required,
			Choices:                  mustGetStringChoises(param),
		}
	case parameters.TypeBool:
		return &dc.BooleanOption{
			OptionName:               defaultParamName,
			Description:              defaultParamDescription,
			OptionNameLocalizations:  localizedParamName,
			DescriptionLocalizations: localizedParamDescription,
			Required:                 param.Required,
		}
	case parameters.TypeInt:
		return &dc.IntegerOption{
			OptionName:               defaultParamName,
			Description:              defaultParamDescription,
			OptionNameLocalizations:  localizedParamName,
			DescriptionLocalizations: localizedParamDescription,
			Required:                 param.Required,
			Choices:                  mustGetIntegerChoises(param),
		}
	default:
		logger.Fatalf("Cannot find discord type for %s", param.Type.BaseType.Name)
	}
	return nil
}
