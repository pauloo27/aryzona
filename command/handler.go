package command

import (
	"github.com/Pauloo27/aryzona/utils"
	"github.com/Pauloo27/logger"
	"github.com/bwmarrin/discordgo"
)

func HandleCommand(commandName string, args []string, s *discordgo.Session, eventCtx *Event) {
	command, ok := commandMap[commandName]
	if !ok {
		return
	}

	ctx := &CommandContext{
		Session: s,
		RawArgs: args,
		Message: eventCtx.Message,
	}

	if command.Permission != nil {
		if !command.Permission.Checker(ctx) {
			ctx.Error(utils.Fmt("This command requires `%s`", command.Permission.Name))
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
