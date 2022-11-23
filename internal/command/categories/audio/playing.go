package audio

import (
	"fmt"
	"strings"

	"github.com/Pauloo27/aryzona/internal/command"
	"github.com/Pauloo27/aryzona/internal/command/validations"
	"github.com/Pauloo27/aryzona/internal/discord"
	"github.com/Pauloo27/aryzona/internal/discord/voicer"
	"github.com/Pauloo27/aryzona/internal/discord/voicer/playable"
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
		requesterID := ctx.Locals["requesterID"].(string)

		embed := buildPlayableInfoEmbed(playing, vc, requesterID).
			WithTitle("Now playing: " + playing.GetName())

		if vc.Queue.Size() > 1 {
			sb := strings.Builder{}
			next := vc.Queue.All()[1:]
			limit := len(next)
			if len(next) > maxNextItems {
				limit = maxNextItems
			}
			for _, item := range next[:limit] {
				playable := item.Playable
				title, artist := playable.GetFullTitle()
				requester := discord.AsMention(item.Requester)
				if artist == "" {
					sb.WriteString(fmt.Sprintf("  -> %s _requested by %s_\n", title, requester))
				} else {
					sb.WriteString(fmt.Sprintf("  -> %s - %s _requested by %s_\n", artist, title, requester))
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
