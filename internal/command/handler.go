package command

import (
	"github.com/Pauloo27/aryzona/internal/discord"
	"github.com/Pauloo27/aryzona/internal/utils"
	"github.com/Pauloo27/logger"
)

func HandleCommand(
	commandName string, args []string,
	adapter *Adapter, bot discord.BotAdapter,
) {
	command, ok := commandMap[commandName]
	if !ok {
		return
	}

	ctx := &CommandContext{
		Bot:      bot,
		RawArgs:  args,
		AuthorID: adapter.AuthorID,
		GuildID:  adapter.GuildID,
		Locals:   make(map[string]interface{}),
		Command:  command,
	}

	// attach adapter
	ctx.Reply = func(msg string) error {
		return adapter.Reply(ctx, msg)
	}
	ctx.ReplyEmbed = func(embed *discord.Embed) error {
		return adapter.ReplyEmbed(ctx, embed)
	}

	if command.Permission != nil {
		if !command.Permission.Checker(ctx) {
			ctx.Error(utils.Fmt("This command requires `%s`", command.Permission.Name))
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

	values, syntaxError := command.ValidateParameters(args)
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

	if command.Deferred && adapter.DeferResponse != nil {
		err := adapter.DeferResponse()
		if err != nil {
			logger.Error("Cannot defer response:", err)
		}
	}

	command.Handler(ctx)
}
