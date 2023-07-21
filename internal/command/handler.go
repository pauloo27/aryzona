package command

import (
	"fmt"
	"strings"

	gonanoid "github.com/matoous/go-nanoid/v2"
	"github.com/pauloo27/aryzona/internal/data/services"
	"github.com/pauloo27/aryzona/internal/discord"
	"github.com/pauloo27/aryzona/internal/discord/model"
	"github.com/pauloo27/aryzona/internal/i18n"
	"github.com/pauloo27/logger"
	"go.opentelemetry.io/otel/codes"
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

	executionID := gonanoid.Must(8)

	trCtx, span := startCommandTrace(command, trigger, executionID)
	defer span.End()

	if trigger.PreferedLanguage == "" {
		_, getUserLangSpan := startChildSpan(trCtx, "GetUserLanguage")
		trigger.PreferedLanguage = services.User.GetLanguage(trigger.AuthorID, trigger.GuildID)
		getUserLangSpan.End()
		getUserLangSpan.SetStatus(codes.Ok, string(trigger.PreferedLanguage))
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
		executionID: executionID,
		trCtx:       trCtx,
		span:        span,
	}

	result := executeCommand(ctx, command)

	if result.Success {
		span.SetStatus(codes.Ok, "Success")
	} else {
		span.SetStatus(codes.Error, "Error")
	}

	_, replySpan := startChildSpan(trCtx, "Reply")

	err := trigger.Reply(ctx, result.Message)

	replySpan.End()

	if err != nil {
		replySpan.SetStatus(codes.Error, "Error")
		logger.Errorf("Cannot reply to command %s: %v", ctx.executionID, err)
	}
	replySpan.SetStatus(codes.Ok, "Success")
}

func executeCommand(
	ctx *Context,
	command *Command,
) Result {
	validaionsI18n := ctx.Lang.Validations.PreCommandValidation

	if command.Deferred && ctx.trigger.DeferResponse != nil {
		addEventToSpan(ctx.span, "DeferResponse")
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
		_, permissionSpan := startChildSpan(ctx.trCtx, "CheckPermission")
		if !command.Permission.Checker(ctx) {
			permissionSpan.End()
			permissionSpan.SetStatus(codes.Error, "PermissionDenied")
			return ctx.Error(
				validaionsI18n.PermissionRequired.Str(command.Permission.Name),
			)
		}
		permissionSpan.End()
		permissionSpan.SetStatus(codes.Ok, "PermissionGranted")
	}

	if command.Validations != nil {
		_, validationSpan := startChildSpan(ctx.trCtx, "RunValidations")

		for _, validation := range command.Validations {
			addEventToSpan(validationSpan, fmt.Sprintf("RunValidation: %s", validation.Name))
			ok, msg := RunValidation(ctx, validation)
			if !ok {
				validationSpan.End()
				validationSpan.SetStatus(codes.Error, fmt.Sprintf("ValidationFailed: %s", validation.Name))
				return ctx.Error(msg)
			}
		}
		validationSpan.End()
		validationSpan.SetStatus(codes.Ok, "ValidationsPassed")
	}

	if command.Parameters != nil {
		_, validationSpan := startChildSpan(ctx.trCtx, "ValidateParameters")

		values, syntaxError := command.ValidateParameters(ctx)
		if syntaxError != nil {
			msg := strings.SplitN(syntaxError.Error(), ":", 2)[1]
			validationSpan.End()
			validationSpan.SetStatus(codes.Error, "ParametersAreInvalid")
			return ctx.Error(msg)
		}
		ctx.Args = values
		validationSpan.End()
		validationSpan.SetStatus(codes.Ok, "ParametersAreValid")
	}

	defer func() {
		if err := recover(); err != nil {
			logger.Errorf("Panic catch while running command %s: %v", command.Name, err)
		}
	}()

	if command.SubCommands == nil || (len(ctx.RawArgs) == 0 && command.Handler != nil) {
		_, handlerSpan := startChildSpan(ctx.trCtx, "RunHandler")

		result := command.Handler(ctx)

		handlerSpan.End()
		handlerSpan.SetStatus(codes.Ok, "HandlerExecuted")

		return result
	}

	// TODO: trace
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
