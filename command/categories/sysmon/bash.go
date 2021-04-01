package sysmon

import (
	"os/exec"
	"strings"

	"github.com/Pauloo27/aryzona/command"
	"github.com/Pauloo27/aryzona/command/permissions"
	"github.com/Pauloo27/aryzona/utils"
)

var Bash = command.Command{
	Name:        "bash",
	Description: "Eval a bash command",
	Permission:  &permissions.BeOwner,
	Handler: func(ctx *command.CommandContext) {
		cmd := exec.Command("bash", "-c", strings.Join(ctx.Args, " "))
		buffer, err := cmd.CombinedOutput()
		output := string(buffer)
		output = strings.ReplaceAll(output, "`", "\\`")
		if err != nil {
			ctx.Error(utils.Fmt("Something went wrong:\n```\n%s\n```", output))
		} else {
			ctx.Success(utils.Fmt("Command ran successfully:\n```\n%s\n```", output))
		}
	},
}
