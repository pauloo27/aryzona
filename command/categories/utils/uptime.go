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
		ctx.ReplyWithEmbed(
			utils.NewEmbedBuilder().
			Title("Bot uptime").
			Color(0xC0FFEE).
			Field("Uptime", utils.FormatDuration(uptime)).
			Field("Started at", discord.StartedAt.Format("2 Jan, 15:04")).
			Build(),
		)
	},
}
