package sysmon

import (
	"os/exec"

	"github.com/Pauloo27/aryzona/command"
	"github.com/Pauloo27/aryzona/command/permissions"
	"github.com/Pauloo27/aryzona/utils"
)

var Eval = command.Command{
	Name:        "eval",
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
		if err != nil {
			ctx.Error(utils.Fmt("Something went wrong: %s", string(buffer)))
		} else {
			ctx.Success(utils.Fmt("Command ran successfully: %s", string(buffer)))
		}
	},
}
