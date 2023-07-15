package audio

import (
	"github.com/pauloo27/aryzona/internal/command"
	"github.com/pauloo27/aryzona/internal/command/parameters"
	"github.com/pauloo27/aryzona/internal/command/validations"
	"github.com/pauloo27/aryzona/internal/discord/voicer"
	"github.com/pauloo27/aryzona/internal/discord/voicer/playable"
	"github.com/pauloo27/aryzona/internal/i18n"
	"github.com/pauloo27/aryzona/internal/providers/youtube"
	"github.com/pauloo27/logger"
)

const maxSearchResults = 5

type SearchContext struct {
	*command.Context

	SearchQuery string
	Results     []*youtube.SearchResult
	Voicer      *voicer.Voicer

	SelectionCh chan Selection
}

var PlayCommand = command.Command{
	Name: "play", Aliases: []string{"p"},
	Deferred:    true,
	Validations: []*command.Validation{validations.MustBeOnAValidVoiceChannel},
	Parameters: []*command.Parameter{
		{Name: "song", Required: true, Type: parameters.ParameterText},
	},
	Handler: func(ctx *command.Context) {
		t := ctx.T.(*i18n.CommandPlay)

		searchQuery := ctx.Args[0].(string)

		results, err := youtube.SearchFor(searchQuery, maxSearchResults)
		if err != nil || len(results) == 0 {
			ctx.Error(t.SomethingWentWrong.Str())
			logger.Warnf("Error searching for %s: %v", searchQuery, err)
			return
		}

		vc := ctx.Locals["vc"].(*voicer.Voicer)
		if !vc.IsConnected() {
			if err := vc.Connect(); err != nil {
				ctx.Error(t.CannotConnect.Str())
				logger.Error(err)
				return
			}
		} else {
			ok, msg := validations.MustBeOnSameVoiceChannel.Checker(ctx)
			if !ok {
				logger.Error(msg)
				return
			}
		}

		// what should be added to the queue, since we support playlists...
		// it needs to be a list of playable
		var toPlay []playable.Playable

		searchCtx := &SearchContext{
			Context:     ctx,
			SearchQuery: searchQuery,
			Voicer:      vc,
			Results:     results,
		}

		if len(results) == 1 {
			toPlay = handleSingleResult(searchCtx)
		} else {
			toPlay = handleMultipleResults(searchCtx)
		}

		if toPlay == nil {
			return
		}

		vc.AppendManyToQueue(ctx.AuthorID, toPlay...)
	},
}

func handleSingleResult(ctx *SearchContext) (toPlay []playable.Playable) {
	t := ctx.T.(*i18n.CommandPlay)

	result := ctx.Results[0]

	toPlay = result.ToPlayable()

	var displayResult playable.Playable

	if result.IsPlaylist() {
		displayResult = playable.DummyPlayable{
			Name:     t.YouTubePlaylist.Str(),
			Artist:   result.Author,
			Title:    result.Title,
			Duration: result.Duration,
		}
	} else {
		displayResult = toPlay[0]
	}

	ctx.SuccessEmbed(
		buildPlayableInfoEmbed(
			PlayableInfo{
				Playable:    displayResult,
				Voicer:      ctx.Voicer,
				RequesterID: ctx.AuthorID,
				T:           t.PlayingInfo,
				Common:      t.Common,
			},
		).
			WithTitle(
				t.BestResult.Str(ctx.SearchQuery),
			),
	)
	return
}
