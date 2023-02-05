package audio

import (
	"fmt"

	"github.com/Pauloo27/aryzona/internal/command"
	"github.com/Pauloo27/aryzona/internal/command/parameters"
	"github.com/Pauloo27/aryzona/internal/discord/voicer"
	"github.com/Pauloo27/aryzona/internal/i18n"
	"github.com/Pauloo27/lyric"
)

var LyricCommand = command.Command{
	Name: "lyric", Aliases: []string{"ly", "lyrics"},
	Parameters: []*command.CommandParameter{
		{Name: "song", Type: parameters.ParameterText, Required: false},
	},
	Handler: func(ctx *command.CommandContext) {
		t := ctx.T.(*i18n.CommandLyric)

		var searchTerms string

		if len(ctx.Args) != 0 {
			searchTerms = ctx.Args[0].(string)
		} else {
			vc := voicer.GetExistingVoicerForGuild(ctx.GuildID)
			if vc == nil {
				ctx.Error(t.NotConnected.Str())
				return
			}
			entry := vc.Playing()
			if entry == nil {
				ctx.Error(t.NothingPlaying.Str())
				return
			}
			playable := entry.Playable

			title, artist := playable.GetFullTitle()
			searchTerms = fmt.Sprintf("%s %s", title, artist)
		}

		result, err := lyric.SearchDDG(searchTerms)
		if err != nil {
			ctx.Error(t.NoResults.Str(searchTerms))
			return
		}
		_ = ctx.Reply(result)
	},
}
