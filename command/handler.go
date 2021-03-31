package command

import (
	"github.com/Pauloo27/aryzona/utils"
	"github.com/bwmarrin/discordgo"
)

func HandleCommand(commandName string, args []string, s *discordgo.Session, m *discordgo.MessageCreate) {
	command, ok := commandMap[commandName]
	if !ok {
		return
	}

	ctx := &CommandContext{m.Message, m, s, args}

	if command.Permission != nil {
		if !command.Permission.Checker(ctx) {
			ctx.Error(utils.Fmt("This command requires `%s`", command.Permission.Name))
			return
		}
	}

	command.Handler(ctx)
}
