package utils

import (
	"strings"

	"github.com/Pauloo27/aryzona/internal/command"
	"github.com/Pauloo27/aryzona/internal/utils"
)

var HelpCommand = command.Command{
	Name: "help", Description: "List all commands",
	Aliases: []string{"h"},
	Handler: func(ctx *command.CommandContext) {
		sb := strings.Builder{}
		sb.WriteString("List of commands:\n")
		lastCategory := ""
		for _, cmd := range command.GetCommandList() {
			if lastCategory != cmd.GetCategory().Name {
				sb.WriteString(utils.Fmt("\n**%s %s**:\n", cmd.GetCategory().Emoji, cmd.GetCategory().Name))
			}
			var permission string
			if cmd.Permission != nil {
				permission = utils.Fmt("(_requires you to... %s_)", cmd.Permission.Name)
			}
			var aliases string
			if len(cmd.Aliases) > 0 {
				aliases = utils.Fmt("(aka %s)", strings.Join(cmd.Aliases, ", "))
			}
			sb.WriteString(utils.Fmt(
				" - `%s%s` %s: **%s** %s\n",
				command.Prefix, cmd.Name, aliases, cmd.Description, permission,
			))
			lastCategory = cmd.GetCategory().Name
		}
		ctx.Success(sb.String())
	},
}
