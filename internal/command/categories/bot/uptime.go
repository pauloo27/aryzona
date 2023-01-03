package bot

import (
	"fmt"
	"runtime"
	"time"

	"github.com/Pauloo27/aryzona/internal/command"
	"github.com/Pauloo27/aryzona/internal/discord"
	"github.com/Pauloo27/aryzona/internal/discord/model"
	"github.com/Pauloo27/aryzona/internal/providers/git"
	"github.com/Pauloo27/aryzona/internal/utils"
)

var UptimeCommand = command.Command{
	Name: "uptime", Description: "Tell how long the bot is running",
	Aliases: []string{"up"},
	Handler: func(ctx *command.CommandContext) {
		uptime := time.Since(*discord.Bot.StartedAt())
		embed := model.NewEmbed().
			WithTitle("Bot uptime").
			WithField(":timer: Uptime", utils.DurationAsText(uptime)).
			WithField(":gear: Implementation", discord.Bot.Implementation()).
			WithField(
				":computer: Host info",
				fmt.Sprintf("Compiled with **%s (%s)**, running on a **%s %s%s**",
					runtime.Version(), runtime.Compiler, runtime.GOOS, runtime.GOARCH,
					utils.ConditionalString(
						utils.IsDocker(),
						" (docker)",
						"",
					),
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
