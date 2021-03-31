package sysmon

import (
	"os/exec"

	"github.com/Pauloo27/aryzona/command"
	"github.com/Pauloo27/aryzona/command/permissions"
	"github.com/Pauloo27/aryzona/utils"
)

var Bash = command.Command{
	Name:        "bash",
	Description: "Eval a bash command",
	Permission:  &permissions.BeOwner,
	Handler: func(ctx *command.CommandContext) {
		name := ctx.Args[0]
		if name == "" {
			ctx.Error("Missing the command name")
			return
		}
		var args []string
		if len(ctx.Args) >= 1 {
			args = ctx.Args[1:]
		}
		cmd := exec.Command(name, args...)
		buffer, err := cmd.CombinedOutput()
		// TODO: escape ```
		if err != nil {
			ctx.Error(utils.Fmt("Something went wrong:\n```\n%s\n```", string(buffer)))
		} else {
			ctx.Success(utils.Fmt("Command ran successfully:\n```\n%s\n```", string(buffer)))
		}
	},
}
