package audio

import (
	"fmt"

	"github.com/pauloo27/aryzona/internal/command"
	"github.com/pauloo27/aryzona/internal/command/parameters"
	"github.com/pauloo27/aryzona/internal/discord/voicer"
	"github.com/pauloo27/aryzona/internal/i18n"
	"github.com/pauloo27/lyric"
)

var LyricCommand = command.Command{
	Name: "lyric", Aliases: []string{"ly", "lyrics"},
	Parameters: []*command.Parameter{
		{Name: "song", Type: parameters.ParameterText, Required: false},
	},
	Handler: func(ctx *command.Context) command.Result {
		t := ctx.T.(*i18n.CommandLyric)

		var searchTerms string

		if len(ctx.Args) != 0 {
			searchTerms = ctx.Args[0].(string)
		} else {
			vc := voicer.GetExistingVoicerForGuild(ctx.GuildID)
			if vc == nil {
				return ctx.Error(t.NotConnected.Str())
			}
			entry := vc.Playing()
			if entry == nil {
				return ctx.Error(t.NothingPlaying.Str())
			}
			playable := entry.Playable

			title, artist := playable.GetFullTitle()
			searchTerms = fmt.Sprintf("%s %s", title, artist)
		}

		result, err := lyric.SearchDDG(searchTerms)
		if err != nil {
			return ctx.Error(t.NoResults.Str(searchTerms))
		}
		return ctx.ReplyRaw(result)
	},
}
