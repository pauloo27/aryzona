package play

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
	Handler: func(ctx *command.Context) command.Result {
		t := ctx.T.(*i18n.CommandPlay)

		input := ctx.Args[0].(string)

		if res := handleSpotifyLink(ctx, input, t); res != nil {
			return *res
		}

		return searchYoutube(ctx, input, t)
	},
}

func searchYoutube(ctx *command.Context, searchQuery string, t *i18n.CommandPlay) command.Result {
	results, err := youtube.SearchFor(searchQuery, maxSearchResults)
	if err != nil || len(results) == 0 {
		logger.Warnf("Error searching for %s: %v", searchQuery, err)
		return ctx.Error(t.SomethingWentWrong.Str())
	}

	vc := ctx.Locals["vc"].(*voicer.Voicer)
	if !vc.IsConnected() {
		if err := vc.Connect(); err != nil {
			logger.Error(err)
			return ctx.Error(t.CannotConnect.Str())
		}
	} else {
		ok, msg := validations.MustBeOnSameVoiceChannel.Checker(ctx)
		if !ok {
			return ctx.Error(msg)
		}
	}

	searchCtx := &SearchContext{
		Context:     ctx,
		SearchQuery: searchQuery,
		Voicer:      vc,
		Results:     results,
	}

	if len(results) == 1 {
		toPlay, cmdResult := handleSingleResult(searchCtx)
		vc.AppendManyToQueue(ctx.AuthorID, toPlay...)
		return cmdResult
	}

	resultCh := make(chan []playable.Playable)
	cmdResult := handleMultipleResults(searchCtx, resultCh)

	go func() {
		toPlay := <-resultCh
		vc.AppendManyToQueue(ctx.AuthorID, toPlay...)
	}()

	return cmdResult
}

func handleSingleResult(ctx *SearchContext) ([]playable.Playable, command.Result) {
	t := ctx.T.(*i18n.CommandPlay)

	result := ctx.Results[0]

	toPlay := result.ToPlayable()

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

	cmdResult := ctx.SuccessEmbed(
		BuildPlayableInfoEmbed(
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
	return toPlay, cmdResult
}
