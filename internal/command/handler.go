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

var (
	CommandLogLevel = logger.Level{
		Name:  "COMMAND",
		Color: "\033[38;5;5m",
	}
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

	ctx := &CommandContext{
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

	executeCommand(ctx, command)
}

func executeCommand(
	ctx *CommandContext,
	command *Command,
) {
	logger.Logf(
		CommandLogLevel,
		"[i %s] (%s) <u %s> <g %s><c %s> executed: %s %s",
		ctx.executionID, ctx.TriggerType,
		ctx.AuthorID, ctx.GuildID, ctx.Channel.ID(), ctx.UsedName, ctx.RawArgs,
	)

	validaionsI18n := ctx.Lang.Validations.PreCommandValidation

	if command.Deferred && ctx.trigger.DeferResponse != nil {
		err := ctx.trigger.DeferResponse()
		if err != nil {
			logger.Error("Cannot defer response:", err)
		}
	}

	if command.Ephemeral && ctx.TriggerType != CommandTriggerSlash {
		ctx.Error(
			validaionsI18n.MustBeExecutedAsSlashCommand.Str(command.Name),
		)
		return
	}

	if command.Permission != nil {
		if !command.Permission.Checker(ctx) {
			ctx.Error(
				validaionsI18n.PermissionRequired.Str(command.Permission.Name),
			)
			return
		}
	}

	for _, validation := range command.Validations {
		ok, msg := RunValidation(ctx, validation)
		if !ok {
			ctx.Error(msg)
			return
		}
	}

	values, syntaxError := command.ValidateParameters(ctx)
	if syntaxError != nil {
		msg := strings.SplitN(syntaxError.Error(), ":", 2)[1]
		ctx.Error(msg)
		return
	}
	ctx.Args = values

	defer func() {
		if err := recover(); err != nil {
			logger.Errorf("Panic catch while running command %s: %v", command.Name, err)
		}
	}()

	if command.SubCommands == nil || (len(ctx.RawArgs) == 0 && command.Handler != nil) {
		command.Handler(ctx)
	} else {
		var subCommandNames []string
		for _, subCommand := range command.SubCommands {
			subCommandNames = append(subCommandNames, subCommand.Name)
		}
		if len(ctx.RawArgs) == 0 {
			ctx.Error(validaionsI18n.MissingSubCommand.Str(subCommandNames))
			return
		}
		subCommandName := ctx.RawArgs[0]
		for _, subCommand := range command.SubCommands {
			subCommand.parent = command
			if subCommand.Name == subCommandName {
				ctx.RawArgs = ctx.RawArgs[1:]
				executeCommand(ctx, subCommand)
				return
			}
			for _, alias := range subCommand.Aliases {
				if alias == subCommandName {
					ctx.RawArgs = ctx.RawArgs[1:]
					executeCommand(ctx, command)
					return
				}
			}
		}
		ctx.Error(validaionsI18n.InvalidSubCommand.Str(subCommandNames))
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
