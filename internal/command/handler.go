package command

import (
	"github.com/Pauloo27/aryzona/internal/discord"
	"github.com/Pauloo27/aryzona/internal/utils"
	"github.com/Pauloo27/logger"
)

func runValidation(ctx *CommandContext, validation *CommandValidation) (bool, string) {
	for _, depends := range validation.DependsOn {
		ok, msg := runValidation(ctx, depends)
		if !ok {
			return ok, msg
		}
	}
	return validation.Checker(ctx)
}

func HandleCommand(
	commandName string, args []string,
	event *Event, bot discord.BotAdapter,
) {
	command, ok := commandMap[commandName]
	if !ok {
		return
	}

	ctx := &CommandContext{
		Bot:        bot,
		RawArgs:    args,
		Reply:      event.Reply,
		ReplyEmbed: event.ReplyEmbed,
		AuthorID:   event.AuthorID,
		GuildID:    event.GuildID,
		Locals:     make(map[string]interface{}),
	}

	if command.Permission != nil {
		if !command.Permission.Checker(ctx) {
			ctx.Error(utils.Fmt("This command requires `%s`", command.Permission.Name))
			return
		}
	}

	for _, validation := range command.Validations {
		ok, msg := runValidation(ctx, validation)
		if !ok {
			ctx.Error(msg)
			return
		}
	}

	values, syntaxError := command.ValidateArguments(args)
	if syntaxError != nil {
		ctx.Error(syntaxError.Message)
		return
	}
	ctx.Args = values

	defer func() {
		if err := recover(); err != nil {
			logger.Errorf("Panic catch while running command %s: %v", command.Name, err)
		}
	}()

	command.Handler(ctx)
}
