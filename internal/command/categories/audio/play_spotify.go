package audio

import (
	"regexp"

	"github.com/pauloo27/aryzona/internal/command"
	"github.com/pauloo27/aryzona/internal/config"
	"github.com/pauloo27/aryzona/internal/discord/voicer"
	"github.com/pauloo27/aryzona/internal/i18n"
	"github.com/pauloo27/aryzona/internal/providers/spotify"
	"github.com/pauloo27/aryzona/internal/providers/youtube"
	"github.com/pauloo27/logger"
)

var (
	SpotifyTrackRe = regexp.MustCompile(`^https?:\/\/open.spotify.com\/track\/([^/]+)$`)

	sfy *spotify.Spotify
)

func handleSpotifyLink(ctx *command.Context, link string, t *i18n.CommandPlay) *command.Result {
	matches := SpotifyTrackRe.FindStringSubmatch(link)

	if len(matches) != 2 {
		return nil
	}

	vc := ctx.Locals["vc"].(*voicer.Voicer)
	trackId := matches[1]

	var res command.Result

	if sfy == nil {
		sfy = spotify.NewSpotify(
			config.Config.SpotifyClientID,
			config.Config.SpotifyClientSecret,
		)
	}

	track, err := sfy.GetTrack(trackId)

	if err != nil {
		logger.Error(err)
		res = ctx.Error(t.SomethingWentWrong.Str())
		return &res
	}

	searchQuery := track.Name + " " + track.Artists[0].Name

	results, err := youtube.SearchFor(searchQuery, 1)
	if err != nil {
		logger.Error(err)
		res = ctx.Error(t.SomethingWentWrong.Str())
		return &res
	}

	if len(results) == 0 {
		logger.Error(err)
		res = ctx.Error(t.SomethingWentWrong.Str())
		return &res
	}

	vc.AppendManyToQueue(ctx.AuthorID, results[0].ToPlayable()...)

	res = ctx.SuccessEmbed(
		buildPlayableInfoEmbed(
			PlayableInfo{
				Playable:    results[0].ToPlayable()[0],
				Voicer:      vc,
				RequesterID: ctx.AuthorID,
				T:           t.PlayingInfo,
				Common:      t.Common,
			},
		).
			WithTitle(
				t.BestResult.Str(link),
			),
	)

	return &res
}
