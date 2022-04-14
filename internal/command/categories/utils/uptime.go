package utils

import (
	"fmt"
	"runtime"
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
		embed := discord.NewEmbed().
			WithTitle("Bot uptime").
			WithField(":timer: Uptime", utils.FormatDuration(uptime)).
			WithField(":gear: Implementation", discord.Bot.Implementation()).
			WithField(
				":computer: System info",
				fmt.Sprintf("%s %s %s",
					runtime.GOOS, runtime.GOARCH, runtime.Version(),
				)).
			WithField(":star: Started at", discord.Bot.StartedAt().Format("2 Jan, 15:04"))

		if git.CommitHash != "" {
			embed.WithField(
				":floppy_disk: Last commit", fmt.Sprintf("[%s (%s)](%s/commit/%s)",
					git.CommitMessage, git.CommitHash[:10], git.RemoteRepo, git.CommitHash),
			)
		}
		ctx.SuccessEmbed(embed)
	},
}
