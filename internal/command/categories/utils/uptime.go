package utils

import (
	"fmt"
	"time"

	"github.com/Pauloo27/aryzona/internal/command"
	"github.com/Pauloo27/aryzona/internal/discord"
	"github.com/Pauloo27/aryzona/internal/providers/git"
	"github.com/Pauloo27/aryzona/internal/utils"
)

var UptimeCommand = command.Command{
	Name: "uptime", Description: "Tell how long the bot is running",
	Aliases: []string{"up"},
	Handler: func(ctx *command.CommandContext) {
		uptime := time.Since(*discord.Bot.StartedAt())
		ctx.SuccessEmbed(
			discord.NewEmbed().
				WithTitle("Bot uptime").
				WithField("Uptime", utils.FormatDuration(uptime)).
				WithField("Implementation", discord.Bot.Implementation()).
				WithField("Last commit", fmt.Sprintf("%s (%s)", git.CommitMessage, git.CommitHash)).
				WithField("Started at", discord.Bot.StartedAt().Format("2 Jan, 15:04")),
		)
	},
}
