package command

import (
	"fmt"
	"time"

	"github.com/Pauloo27/aryzona/internal/discord"
	"github.com/Pauloo27/logger"
)

func executeCommand(
	command *Command, ctx *CommandContext,
	adapter *Adapter, bot discord.BotAdapter,
) {
	if command.Deferred && adapter.DeferResponse != nil {
		err := adapter.DeferResponse()
		if err != nil {
			logger.Error("Cannot defer response:", err)
		}
	}

	if command.Permission != nil {
		if !command.Permission.Checker(ctx) {
			ctx.Error(fmt.Sprintf("This command requires `%s`", command.Permission.Name))
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
			ctx.Error(fmt.Sprintf("Missing sub command. Available sub commands: %v", subCommandNames))
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
		ctx.Error(fmt.Sprintf("Unknown sub command. Available sub commands: %v", subCommandNames))
	}
}

func HandleCommand(
	commandName string, args []string,
	adapter *Adapter, bot discord.BotAdapter,
	trigger CommandTrigger,
) {
	command, ok := commandMap[commandName]
	if !ok {
		return
	}

	ctx := &CommandContext{
		Bot:       bot,
		RawArgs:   args,
		AuthorID:  adapter.AuthorID,
		UsedName:  commandName,
		GuildID:   adapter.GuildID,
		Locals:    make(map[string]interface{}),
		Command:   command,
		startDate: time.Now(),
		Trigger:   trigger,
	}

	// attach adapter
	ctx.Reply = func(msg string) error {
		return adapter.Reply(ctx, msg)
	}
	ctx.ReplyEmbed = func(embed *discord.Embed) error {
		return adapter.ReplyEmbed(ctx, embed)
	}

	executeCommand(command, ctx, adapter, bot)
}
