package command

import (
	"errors"
	"strings"

	gonanoid "github.com/matoous/go-nanoid/v2"
	"github.com/pauloo27/aryzona/internal/db"
	"github.com/pauloo27/aryzona/internal/db/entity"
	"github.com/pauloo27/aryzona/internal/discord"
	"github.com/pauloo27/aryzona/internal/discord/model"
	"github.com/pauloo27/aryzona/internal/i18n"
	"github.com/pauloo27/logger"
	"xorm.io/xorm"
)

func GetCommand(name string) *Command {
	return commandMap[name]
}

func HandleCommand(
	usedCommandName string,
	args []string,
	command *Command,
	bot discord.BotAdapter,
	trigger *TriggerEvent,
) {
	if trigger.PreferedLanguage == "" {
		trigger.PreferedLanguage = getUserLanguage(trigger.AuthorID, trigger.GuildID)
	}

	lang := i18n.MustGetLanguage(trigger.PreferedLanguage)

	t := i18n.GetCommand(lang, command.Name)

	ctx := &Context{
		Bot:         bot,
		T:           t,
		Lang:        lang,
		RawArgs:     args,
		MessageID:   trigger.MessageID,
		AuthorID:    trigger.AuthorID,
		UsedName:    usedCommandName,
		GuildID:     trigger.GuildID,
		Locals:      make(map[string]any),
		Command:     command,
		startTime:   trigger.EventTime,
		TriggerType: trigger.Type,
		Channel:     trigger.Channel,
		trigger:     trigger,
		executionID: gonanoid.Must(5),
	}

	result := executeCommand(ctx, command)
	trigger.Reply(ctx, result.Message)
}

func executeCommand(
	ctx *Context,
	command *Command,
) Result {
	validaionsI18n := ctx.Lang.Validations.PreCommandValidation

	if command.Deferred && ctx.trigger.DeferResponse != nil {
		err := ctx.trigger.DeferResponse()
		if err != nil {
			logger.Error("Cannot defer response:", err)
		}
	}

	if command.Ephemeral && ctx.TriggerType != CommandTriggerSlash {
		return ctx.Error(
			validaionsI18n.MustBeExecutedAsSlashCommand.Str(command.Name),
		)
	}

	if command.Permission != nil {
		if !command.Permission.Checker(ctx) {
			return ctx.Error(
				validaionsI18n.PermissionRequired.Str(command.Permission.Name),
			)
		}
	}

	for _, validation := range command.Validations {
		ok, msg := RunValidation(ctx, validation)
		if !ok {
			return ctx.Error(msg)
		}
	}

	values, syntaxError := command.ValidateParameters(ctx)
	if syntaxError != nil {
		msg := strings.SplitN(syntaxError.Error(), ":", 2)[1]
		return ctx.Error(msg)
	}
	ctx.Args = values

	defer func() {
		if err := recover(); err != nil {
			logger.Errorf("Panic catch while running command %s: %v", command.Name, err)
		}
	}()

	if command.SubCommands == nil || (len(ctx.RawArgs) == 0 && command.Handler != nil) {
		return command.Handler(ctx)
	} else {
		var subCommandNames []string
		for _, subCommand := range command.SubCommands {
			subCommandNames = append(subCommandNames, subCommand.Name)
		}
		if len(ctx.RawArgs) == 0 {
			return ctx.Error(validaionsI18n.MissingSubCommand.Str(subCommandNames))
		}
		subCommandName := ctx.RawArgs[0]
		for _, subCommand := range command.SubCommands {
			subCommand.parent = command
			if subCommand.Name == subCommandName {
				ctx.RawArgs = ctx.RawArgs[1:]
				return executeCommand(ctx, subCommand)
			}
			for _, alias := range subCommand.Aliases {
				if alias == subCommandName {
					ctx.RawArgs = ctx.RawArgs[1:]
					return executeCommand(ctx, command)
				}
			}
		}
		return ctx.Error(validaionsI18n.InvalidSubCommand.Str(subCommandNames))
	}
}

func HandleInteraction(fullID, userID string) *model.ComplexMessage {
	splitted := strings.Split(fullID, ":")
	if len(splitted) != 2 {
		logger.Error("Invalid interaction id", fullID)
		return nil
	}
	baseID, id := splitted[0], splitted[1]
	ctx, ok := commandInteractionMap[baseID]
	if !ok {
		logger.Error("Cannot find interaction adapter for id", baseID)
		return nil
	}
	newMessage, done := ctx.interactionHandler(id, userID, baseID)
	if done {
		delete(commandInteractionMap, baseID)
	}
	return newMessage
}

func getUserLanguage(userID, guildID string) i18n.LanguageName {
	var user = entity.User{ID: userID}

	hasUser, err := db.DB.Get(&user)
	if err != nil && !errors.Is(err, xorm.ErrNotExist) {
		logger.Error(err)
	}

	if hasUser {
		if user.PreferredLocale != "" {
			return user.PreferredLocale
		}

		if user.LastSlashCommandLocale != "" {
			return user.LastSlashCommandLocale
		}
	}

	if guildID == "" {
		return i18n.DefaultLanguageName
	}

	var guild = entity.Guild{ID: guildID}

	hasGuild, err := db.DB.Get(&guild)

	if err != nil && !errors.Is(err, xorm.ErrNotExist) {
		logger.Error(err)
	}

	if hasGuild {
		return guild.PreferredLocale
	}

	return i18n.DefaultLanguageName
}
