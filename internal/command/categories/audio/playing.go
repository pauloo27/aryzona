package audio

import (
	"fmt"
	"strings"
	"time"

	"github.com/pauloo27/aryzona/internal/command"
	"github.com/pauloo27/aryzona/internal/command/validations"
	"github.com/pauloo27/aryzona/internal/core/f"
	"github.com/pauloo27/aryzona/internal/discord"
	"github.com/pauloo27/aryzona/internal/discord/voicer"
	"github.com/pauloo27/aryzona/internal/discord/voicer/playable"
	"github.com/pauloo27/aryzona/internal/i18n"
)

const (
	maxNextItems = 5
)

var PlayingCommand = command.Command{
	Name: "playing", Aliases: []string{"np"},
	Deferred:    true,
	Validations: []*command.Validation{validations.MustBePlaying},
	Handler: func(ctx *command.Context) command.Result {
		t := ctx.T.(*i18n.CommandPlaying)

		vc := ctx.Locals["vc"].(*voicer.Voicer)
		playing := ctx.Locals["playing"].(playable.Playable)
		requesterID := ctx.Locals["requesterID"].(string)

		embed := buildPlayableInfoEmbed(
			PlayableInfo{
				Playable:    playing,
				Voicer:      vc,
				RequesterID: requesterID,
				T:           t.PlayingInfo,
				Common:      t.Common,
			},
		).
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
					etaStr = t.PlayingInfo.ETANever.Str()
				} else {
					etaStr = f.DurationAsDetailedDiffText(eta, t.Common)
				}

				title, artist := playable.GetFullTitle()
				requester := discord.AsMention(item.Requester)

				var fullTitle string

				if artist == "" {
					fullTitle = title
				} else {
					fullTitle = fmt.Sprintf("%s - %s", artist, title)
				}

				sb.WriteString(t.Entry.Str(fullTitle, requester, etaStr) + "\n")
			}
			if len(next) > maxNextItems {
				sb.WriteString(t.AndMore.Str())
			}

			var queueDuration time.Duration

			for _, item := range next {
				duration, err := item.Playable.GetDuration()
				if err != nil {
					continue
				}
				queueDuration += duration
			}

			embed.WithField(t.ComingNext.Str(len(next), f.Pluralize(len(next), t.Song.Str(), t.Songs.Str()), f.ShortDuration(queueDuration)), sb.String())
		}

		return ctx.SuccessEmbed(embed)
	},
}
