package sysmon

import (
	"os/exec"
	"strings"

	"github.com/Pauloo27/aryzona/command"
	"github.com/Pauloo27/aryzona/command/permissions"
	"github.com/Pauloo27/aryzona/utils"
)

// should i remove it? probably...
/* #nosec G204 */
var Bash = command.Command{
	Name:        "bash",
	Description: "Eval a bash command",
	Permission:  &permissions.BeOwner,
	Handler: func(ctx *command.CommandContext) {
		if len(ctx.RawArgs) == 0 {
			ctx.Error("Missing bash command")
			return
		}
		cmd := exec.Command("bash", "-c", strings.Join(ctx.RawArgs, " "))
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
