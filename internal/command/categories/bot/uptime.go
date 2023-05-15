package bot

import (
	"fmt"
	"os"
	"runtime"
	"time"

	"github.com/pauloo27/aryzona/internal/command"
	"github.com/pauloo27/aryzona/internal/core/f"
	"github.com/pauloo27/aryzona/internal/discord"
	"github.com/pauloo27/aryzona/internal/discord/model"
	"github.com/pauloo27/aryzona/internal/i18n"
	"github.com/pauloo27/aryzona/internal/providers/git"

	k "github.com/pauloo27/toolkit"
)

var UptimeCommand = command.Command{
	Name:    "uptime",
	Aliases: []string{"up"},
	Handler: func(ctx *command.CommandContext) {
		t := ctx.T.(*i18n.CommandUptime)

		uptime := time.Since(*discord.Bot.StartedAt())
		embed := model.NewEmbed().
			WithTitle(t.Title.Str()).
			WithField(t.Uptime.Str(":timer:"), f.DurationAsText(uptime, t.Common)).
			WithField(t.Implementation.Str(":gear:"), discord.Bot.Implementation()).
			WithField(t.Language.Str(":globe_with_meridians:"), t.Meta.DisplayName.Str()).
			WithField(
				t.HostInfoKey.Str(":computer:"),
				t.HostInfoValue.Str(runtime.Version(), runtime.Compiler, runtime.GOOS, runtime.GOARCH,
					k.Is(
						isDocker(),
						" (docker)",
						"",
					),
				),
			).
			WithField(
				t.StartedAt.Str(":star:"),
				t.FormatSimpleDateTime(*discord.Bot.StartedAt()),
			)

		if git.CommitHash != "" {
			embed.WithField(
				t.LastCommit.Str(":floppy_disk:"),
				fmt.Sprintf(
					"[%s (%s)](%s/commit/%s)",
					git.CommitMessage, git.CommitHash[:10], git.RemoteRepo, git.CommitHash,
				),
			)
		}
		ctx.SuccessEmbed(embed)
	},
}

func isDocker() bool {
	_, err := os.Stat("/.dockerenv")
	return !os.IsNotExist(err)
}
