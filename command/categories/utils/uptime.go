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
		uptime := time.Since(*discord.Bot.StartedAt())
		ctx.SuccessEmbed(
			discord.NewEmbed().
				WithTitle("Bot uptime").
				WithField("Uptime", utils.FormatDuration(uptime)).
				WithField("Started at", discord.Bot.StartedAt().Format("2 Jan, 15:04")),
		)
	},
}
