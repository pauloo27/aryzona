package sysmon

import (
	"os"
	"os/exec"
	"strings"

	"github.com/Pauloo27/aryzona/internal/command"
	"github.com/Pauloo27/aryzona/internal/command/parameters"
	"github.com/Pauloo27/aryzona/internal/command/permissions"
	"github.com/Pauloo27/aryzona/internal/utils"
)

// should i remove it? probably...
var Bash = command.Command{
	Name:        "bash",
	Description: "Eval a bash command",
	Permission:  &permissions.MustBeOwner,
	Parameters: []*command.CommandParameter{
		{
			Name: "command", Description: "command to execute", Required: true,
			RequiredMessage: "Missing command", Type: parameters.ParameterText,
		},
	},
	Handler: func(ctx *command.CommandContext) {
		cmd := exec.Command(os.Getenv("DC_BOT_SHELL"), "-c", (ctx.Args[0].(string))) //#nosec G204
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
