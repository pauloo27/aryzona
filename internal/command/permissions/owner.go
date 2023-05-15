package permissions

import (
	"github.com/pauloo27/aryzona/internal/command"
	"github.com/pauloo27/aryzona/internal/config"
)

var MustBeOwner = &command.CommandPermission{
	Name: "be the bot owner",
	Checker: func(ctx *command.CommandContext) bool {
		return ctx.AuthorID == config.Config.OwnerID
	},
}
