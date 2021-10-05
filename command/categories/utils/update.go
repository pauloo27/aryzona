package utils

import (
	"github.com/Pauloo27/aryzona/command"
	"github.com/Pauloo27/aryzona/command/permissions"
	"github.com/Pauloo27/aryzona/command/slash"
)

var UpdateCommand = command.Command{
	Name:        "update",
	Description: "Update slash commands",
	Permission:  &permissions.BeOwner,
	Handler: func(ctx *command.CommandContext) {
		ctx.Success("Wait a little...")
		err := slash.RegisterCommands(true)
		if err != nil {
			ctx.Error(err.Error())
			return
		}
		ctx.Success("All done!")
	},
}
