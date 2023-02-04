package audio

import (
	"fmt"
	"strings"

	"github.com/Pauloo27/aryzona/internal/command"
	"github.com/Pauloo27/aryzona/internal/command/validations"
	"github.com/Pauloo27/aryzona/internal/core/f"
	"github.com/Pauloo27/aryzona/internal/discord"
	"github.com/Pauloo27/aryzona/internal/discord/voicer"
	"github.com/Pauloo27/aryzona/internal/discord/voicer/playable"
	"github.com/Pauloo27/aryzona/internal/i18n"
)

const (
	maxNextItems = 5
)

var PlayingCommand = command.Command{
	Name: "playing", Aliases: []string{"np"},
	Deferred:    true,
	Description: "Show what is playing now",
	Validations: []*command.CommandValidation{validations.MustBePlaying},
	Handler: func(ctx *command.CommandContext) {
		t := ctx.T.(*i18n.CommandPlaying)

		vc := ctx.Locals["vc"].(*voicer.Voicer)
		playing := ctx.Locals["playing"].(playable.Playable)
		requesterID := ctx.Locals["requesterID"].(string)

		embed := buildPlayableInfoEmbed(playing, vc, requesterID).
			WithTitle(t.Title.Str(playing.GetName()))

		if vc.Queue.Size() > 1 {
			sb := strings.Builder{}
			next := vc.Queue.All()[1:]
			limit := len(next)
			if len(next) > maxNextItems {
				limit = maxNextItems
			}

			for _, item := range next[:limit] {
				var etaStr string
				playable := item.Playable

				eta := calcETA(playable, vc)
				if eta == -1 {
					etaStr = t.Never.Str()
				} else {
					etaStr = f.DurationAsDetailedDiffText(eta)
				}

				title, artist := playable.GetFullTitle()
				requester := discord.AsMention(item.Requester)

				var fullTitle string

				if artist == "" {
					fullTitle = title
				} else {
					fullTitle = fmt.Sprintf("%s - %s", artist, title)
				}

				sb.WriteString(t.Entry.Str(fullTitle, requester, etaStr))
			}
			if len(next) > maxNextItems {
				sb.WriteString(t.AndMore.Str())
			}
			embed.WithField(t.ComingNext.Str(len(next)), sb.String())
		}

		ctx.SuccessEmbed(embed)
	},
}
