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
		uptime := time.Since(discord.StartedAt)
		ctx.SuccessEmbed(
			utils.NewEmbedBuilder().
				Title("Bot uptime").
				Field("Uptime", utils.FormatDuration(uptime)).
				Field("Started at", discord.StartedAt.Format("2 Jan, 15:04")).
				Build(),
		)
	},
}
