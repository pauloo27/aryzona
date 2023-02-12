package audio

import (
	"errors"
	"time"

	"github.com/Pauloo27/aryzona/internal/audio/dca"
	"github.com/Pauloo27/aryzona/internal/command"
	"github.com/Pauloo27/aryzona/internal/command/parameters"
	"github.com/Pauloo27/aryzona/internal/command/validations"
	"github.com/Pauloo27/aryzona/internal/core/routine"
	"github.com/Pauloo27/aryzona/internal/discord/voicer"
	"github.com/Pauloo27/aryzona/internal/discord/voicer/playable"
	"github.com/Pauloo27/aryzona/internal/i18n"
	"github.com/Pauloo27/aryzona/internal/providers/youtube"
	"github.com/Pauloo27/logger"
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
		var connErrCh *chan error

		if !vc.IsConnected() {
			ch := make(chan error)
			connErrCh = &ch
			go func() {
				ch <- vc.Connect()
			}()
		} else {
			authorVoiceChannelID, found := ctx.Locals["authorVoiceChannelID"]
			if !found || *(vc.ChannelID) != authorVoiceChannelID.(string) {
				ctx.Error(t.NotInRightChannel.Str())
				return
			}
		}

		searchQuery := ctx.Args[0].(string)
		results, err := youtube.SearchFor(searchQuery, maxSearchResults)
		if err != nil || len(results) == 0 {
			ctx.Error(t.SomethingWentWrong.Str())
			logger.Warnf("Error searching for %s: %v", searchQuery, err)
			return
		}

		if connErrCh != nil {
			if err := <-(*connErrCh); err != nil {
				ctx.Error(t.CannotConnect.Str())
				return
			}
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
			ctx.Error(t.SomethingWentWrong.Str())
			return
		}

		routine.Go(func() {
			if err = vc.AppendManyToQueue(ctx.AuthorID, toPlay...); err != nil {
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
			},
		).
			WithTitle(
				t.BestResult.Str(ctx.SearchQuery),
			),
	)

	return
}
