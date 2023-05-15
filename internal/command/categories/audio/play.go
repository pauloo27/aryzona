package audio

import (
	"errors"
	"sync"
	"time"

	"github.com/pauloo27/aryzona/internal/audio/dca"
	"github.com/pauloo27/aryzona/internal/command"
	"github.com/pauloo27/aryzona/internal/command/parameters"
	"github.com/pauloo27/aryzona/internal/command/validations"
	"github.com/pauloo27/aryzona/internal/core/routine"
	"github.com/pauloo27/aryzona/internal/discord/voicer"
	"github.com/pauloo27/aryzona/internal/discord/voicer/playable"
	"github.com/pauloo27/aryzona/internal/i18n"
	"github.com/pauloo27/aryzona/internal/providers/youtube"
	"github.com/pauloo27/logger"
)

const (
	maxSearchResults       = 5
	firstResultTimeout     = 5 * time.Second
	multipleResultsTimeout = 30 * time.Second
)

type SearchContext struct {
	*command.CommandContext

	SearchQuery string
	Results     []*youtube.SearchResult
	Voicer      *voicer.Voicer
}

var PlayCommand = command.Command{
	Name: "play", Aliases: []string{"p"},
	Deferred:    true,
	Validations: []*command.CommandValidation{validations.MustBeOnAValidVoiceChannel},
	Parameters: []*command.CommandParameter{
		{Name: "song", Required: true, Type: parameters.ParameterText},
	},
	Handler: func(ctx *command.CommandContext) {
		t := ctx.T.(*i18n.CommandPlay)

		vc := ctx.Locals["vc"].(*voicer.Voicer)
		searchQuery := ctx.Args[0].(string)
		hasResultsAndIsConnected := true

		var results []*youtube.SearchResult
		var wg sync.WaitGroup
		var err error

		if !vc.IsConnected() {
			wg.Add(1)
			go func() {
				defer wg.Done()
				if err := vc.Connect(); err != nil {
					hasResultsAndIsConnected = false
					ctx.Error(t.CannotConnect.Str())
					logger.Error(err)
				}
			}()
		} else {
			ok, msg := validations.MustBeOnSameVoiceChannel.Checker(ctx)
			if !ok {
				logger.Error(msg)
				return
			}
		}

		wg.Add(1)
		go func() {
			defer wg.Done()
			results, err = youtube.SearchFor(searchQuery, maxSearchResults)
			if err != nil || len(results) == 0 {
				hasResultsAndIsConnected = false
				ctx.Error(t.SomethingWentWrong.Str())
				logger.Warnf("Error searching for %s: %v", searchQuery, err)
			}
		}()

		wg.Wait()
		if !hasResultsAndIsConnected {
			return
		}

		// what should be added to the queue, since we support playlists...
		// it needs to be a list of playable
		var toPlay []playable.Playable

		searchCtx := &SearchContext{
			CommandContext: ctx,
			SearchQuery:    searchQuery,
			Voicer:         vc,
			Results:        results,
		}

		if len(results) > 1 {
			toPlay = handleMultipleResults(searchCtx)
		} else {
			toPlay = handleSingleResult(searchCtx)
		}

		if toPlay == nil {
			return
		}

		routine.Go(func() {
			if err := vc.AppendManyToQueue(ctx.AuthorID, toPlay...); err != nil {
				if errors.Is(err, dca.ErrVoiceConnectionClosed) {
					return
				}
				ctx.Error(t.SomethingWentWrong.Str())
				logger.Error(err)
				return
			}
		})
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
