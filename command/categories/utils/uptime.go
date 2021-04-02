package utils

import (
	"time"

	"github.com/Pauloo27/aryzona/command"
	"github.com/Pauloo27/aryzona/discord"
	"github.com/Pauloo27/aryzona/utils"
)

var UptimeCommand = command.Command{
	Name: "uptime", Description: "Get bot time up",
	Aliases: []string{"up"},
	Handler: func(ctx *command.CommandContext) {
		uptime := time.Now().Sub(discord.StartedAt)
		ctx.Success(
			utils.Fmt("Bot started at %v. Up %s.",
				discord.StartedAt.Format("2 Jan, 15:06"), utils.FormatDuration(uptime),
			),
		)
	},
}
