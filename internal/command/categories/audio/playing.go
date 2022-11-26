package audio

import (
	"fmt"
	"strings"

	"github.com/Pauloo27/aryzona/internal/command"
	"github.com/Pauloo27/aryzona/internal/command/validations"
	"github.com/Pauloo27/aryzona/internal/discord/voicer"
	"github.com/Pauloo27/aryzona/internal/discord/voicer/playable"
	"github.com/Pauloo27/aryzona/internal/utils"
)

const (
	maxNextItems = 5
)

var PlayingCommand = command.Command{
	Name: "playing", Aliases: []string{"np", "nowplaying", "tocando"},
	Deferred:    true,
	Description: "Show what is playing now",
	Validations: []*command.CommandValidation{validations.MustBePlaying},
	Handler: func(ctx *command.CommandContext) {
		vc := ctx.Locals["vc"].(*voicer.Voicer)
		playing := ctx.Locals["playing"].(playable.Playable)

		embed := buildPlayableInfoEmbed(playing, vc).
			WithTitle("Now playing: " + playing.GetName())

		if vc.Queue.Size() > 1 {
			sb := strings.Builder{}
			next := vc.Queue.All()[1:]
			limit := len(next)
			if len(next) > maxNextItems {
				limit = maxNextItems
			}

			for _, item := range next[:limit] {
				var etaStr string

				eta := calcETA(item, vc)
				if eta == -1 {
					etaStr = "_Never_"
				} else {
					etaStr = utils.DurationAsDetailedDiffText(eta)
				}

				title, artist := item.GetFullTitle()
				if artist == "" {
					sb.WriteString(fmt.Sprintf("  -> %s (playing %s)\n", title, etaStr))
				} else {
					sb.WriteString(
						fmt.Sprintf("  -> %s - %s (playing %s)\n",
							artist,
							title,
							etaStr,
						),
					)
				}
			}
			if len(next) > maxNextItems {
				sb.WriteString("_... and more ..._")
			}
			embed.WithField(fmt.Sprintf("**Coming next (%d):**", len(next)), sb.String())
		}

		ctx.SuccessEmbed(embed)
	},
}
