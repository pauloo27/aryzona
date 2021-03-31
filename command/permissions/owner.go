package permissions

import (
	"os"

	"github.com/Pauloo27/aryzona/command"
)

var BeOwner = command.CommandPermission{
	Name: "be the bot owner",
	Checker: func(ctx *command.CommandContext) bool {
		return ctx.Message.Author.ID == os.Getenv("DC_BOT_OWNER_ID")
	},
}
