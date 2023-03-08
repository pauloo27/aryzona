package command

import (
	"time"

	"github.com/Pauloo27/aryzona/internal/discord"
	"github.com/Pauloo27/aryzona/internal/discord/model"
	"github.com/Pauloo27/aryzona/internal/i18n"
	"github.com/Pauloo27/logger"
	"github.com/matoous/go-nanoid/v2"
)

var (
	CommandLogLevel = logger.Level{
		Name:  "COMMAND",
		Color: "\033[38;5;5m",
	}
)

func executeCommand(
	command *Command, ctx *CommandContext,
	adapter *Adapter, bot discord.BotAdapter,
) {
	logger.Logf(
		CommandLogLevel,
		"[i %s] (%s) <u %s> <g %s><c %s> executed: %s %s",
		ctx.executionID, ctx.Trigger,
		ctx.AuthorID, ctx.GuildID, ctx.Channel.ID(), ctx.UsedName, ctx.RawArgs,
	)

	if command.Deferred && adapter.DeferResponse != nil {
		err := adapter.DeferResponse()
		if err != nil {
			logger.Error("Cannot defer response:", err)
		}
	}

	if command.Ephemeral && ctx.Trigger != CommandTriggerSlash {
		ctx.Errorf("That command must be executed in a slash command. **Use `/%s` instead**", command.Name)
		return
	}

	if command.Permission != nil {
		if !command.Permission.Checker(ctx) {
			ctx.Errorf("This command requires `%s`", command.Permission.Name)
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

	values, syntaxError := command.ValidateParameters(ctx.RawArgs)
	if syntaxError != nil {
		msg := syntaxError.Error()
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
			ctx.Errorf("Missing sub command. Available sub commands: %v", subCommandNames)
			return
		}
		subCommandName := ctx.RawArgs[0]
		for _, subCommand := range command.SubCommands {
			if subCommand.Name == subCommandName {
				ctx.RawArgs = ctx.RawArgs[1:]
				executeCommand(subCommand, ctx, adapter, bot)
				return
			}
			for _, alias := range subCommand.Aliases {
				if alias == subCommandName {
					ctx.RawArgs = ctx.RawArgs[1:]
					executeCommand(subCommand, ctx, adapter, bot)
					return
				}
			}
		}
		ctx.Errorf("Unknown sub command. Available sub commands: %v", subCommandNames)
	}
}

func GetCommandLang(commandName string) i18n.LanguageName {
	lang, found := commandLangMap[commandName]
	if found {
		return lang
	}
	return i18n.DefaultLanguageName
}

func HandleCommand(
	commandName string, args []string,
	langName i18n.LanguageName,
	startTime time.Time,
	adapter *Adapter, bot discord.BotAdapter,
	trigger CommandTrigger, channel model.TextChannel,
) {
	command, ok := commandMap[commandName]
	if !ok {
		return
	}

	lang := i18n.MustGetLanguage(langName)

	t := i18n.GetCommand(lang, command.Name)

	ctx := &CommandContext{
		Bot:       bot,
		T:         t,
		Lang:      lang,
		RawArgs:   args,
		AuthorID:  adapter.AuthorID,
		UsedName:  commandName,
		GuildID:   adapter.GuildID,
		Locals:    make(map[string]any),
		Command:   command,
		startTime: startTime,
		Trigger:   trigger,
		Channel:   channel,
	}

	ctx.executionID = gonanoid.Must(5)

	logResponse := func() {
		logger.Logf(
			CommandLogLevel,
			"[i %s] got response %s, took %s",
			ctx.executionID,
			// there's no need to log the response, also, not microsoft here
			"<omitted>",
			ctx.processTime,
		)
	}

	// attach adapter
	ctx.Reply = func(msg string) error {
		logResponse()
		return adapter.Reply(ctx, msg)
	}
	ctx.ReplyEmbed = func(embed *model.Embed) error {
		logResponse()
		return adapter.ReplyEmbed(ctx, embed)
	}
	ctx.Edit = func(msg string) error {
		logResponse()
		return adapter.Edit(ctx, msg)
	}
	ctx.EditEmbed = func(embed *model.Embed) error {
		logResponse()
		return adapter.EditEmbed(ctx, embed)
	}
	ctx.ReplyComplex = func(data *model.ComplexMessage) error {
		logResponse()
		return adapter.ReplyComplex(ctx, data)
	}
	ctx.EditComplex = func(data *model.ComplexMessage) error {
		logResponse()
		return adapter.EditComplex(ctx, data)
	}

	executeCommand(command, ctx, adapter, bot)
}

func HandleInteraction(id, userID string) *model.ComplexMessage {
	baseID := id[:10]
	ctx, ok := commandInteractionMap[baseID]
	if !ok {
		logger.Error("Cannot find interaction adapter for id", baseID)
		return nil
	}
	newMessage, done := ctx.interactionHandler(id, baseID, userID)
	if done {
		delete(commandInteractionMap, baseID)
	}
	return newMessage
}
