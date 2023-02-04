package audio

import (
	"errors"
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/Pauloo27/aryzona/internal/audio/dca"
	"github.com/Pauloo27/aryzona/internal/command"
	"github.com/Pauloo27/aryzona/internal/command/parameters"
	"github.com/Pauloo27/aryzona/internal/command/validations"
	"github.com/Pauloo27/aryzona/internal/core/f"
	"github.com/Pauloo27/aryzona/internal/core/routine"
	"github.com/Pauloo27/aryzona/internal/discord/model"
	"github.com/Pauloo27/aryzona/internal/discord/voicer"
	"github.com/Pauloo27/aryzona/internal/discord/voicer/playable"
	"github.com/Pauloo27/aryzona/internal/i18n"
	"github.com/Pauloo27/aryzona/internal/providers/youtube"
	"github.com/Pauloo27/logger"
	k "github.com/Pauloo27/toolkit"
)

const (
	maxSearchResults       = 5
	firstResultTimeout     = 5 * time.Second
	multipleResultsTimeout = 30 * time.Second
)

var PlayCommand = command.Command{
	Name: "play", Aliases: []string{"p"},
	Description: "Play a video/song from youtube",
	Deferred:    true,
	Validations: []*command.CommandValidation{validations.MustBeOnAValidVoiceChannel},
	Parameters: []*command.CommandParameter{
		{Name: "song", Description: "Search query or URL", Required: true, Type: parameters.ParameterText},
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

		if len(results) > 1 {
			toPlay = handleMultipleResults(ctx, vc, searchQuery, results)
		} else {
			toPlay = handleSingleResult(ctx, vc, searchQuery, results[0])
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

func handleSingleResult(ctx *command.CommandContext, vc *voicer.Voicer, searchQuery string, result *youtube.SearchResult) (toPlay []playable.Playable) {
	t := ctx.T.(*i18n.CommandPlay)

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
		buildPlayableInfoEmbed(displayResult, vc, ctx.AuthorID).WithTitle(t.BestResult.Str(searchQuery)),
	)

	return
}

func handleMultipleResults(ctx *command.CommandContext, vc *voicer.Voicer, searchQuery string, results []*youtube.SearchResult) []playable.Playable {
	t := ctx.T.(*i18n.CommandPlay)

	selectionLock := sync.Mutex{}
	selectionCh := make(chan *youtube.SearchResult)
	selected := false

	selectResult := func(result *youtube.SearchResult) bool {
		selectionLock.Lock()
		defer selectionLock.Unlock()
		if selected {
			return false
		}
		selectionCh <- result
		if result != nil {
			selected = true
		}
		return true
	}
	firstResult := results[0]

	embed := model.NewEmbed().
		WithColor(command.PendingEmbedColor).
		WithTitle(t.MultipleResults.Str()).
		WithDescription(
			t.FirstResultWillPlay.Str(
				firstResult.Title,
				firstResult.ID,
				firstResult.Author,
				k.Is(firstResult.Duration == 0, t.Live.Str(":red_circle:"), f.ShortDuration(firstResult.Duration)),
				firstResultTimeout/time.Second,
			),
		)

	var components []model.MessageComponent

	baseID, err := ctx.RegisterInteractionHandler(
		func(fullID, baseID, userID string) (msg *model.ComplexMessage, done bool) {
			if userID != ctx.AuthorID {
				return nil, false
			}
			if fullID[len(baseID):] == "play-now" {
				ok := selectResult(firstResult)
				if !ok {
					return nil, false
				}
				return &model.ComplexMessage{
					Components: buildDisabledComponents(components, 0),
					Embeds: []*model.Embed{
						buildPlayableInfoEmbed(firstResult.ToPlayable()[0], vc, ctx.AuthorID).
							WithTitle(t.SelectedResult.Str(searchQuery)).
							WithColor(command.SuccessEmbedColor),
					},
				}, true
			}
			if !selectResult(nil) {
				return nil, false
			}
			return buildMultipleResultsMessage(ctx, vc, searchQuery, results, selectResult), true
		},
	)
	if err != nil {
		logger.Errorf("Cannot register interaction handler: %v", err)
		return nil
	}

	components = []model.MessageComponent{
		model.ButtonComponent{
			Label: t.ConfirmBtn.Str(),
			Style: model.SuccessButtonStyle,
			ID:    fmt.Sprintf("%splay-now", baseID),
		},
		model.ButtonComponent{
			Label: t.PlayOtherBtn.Str(),
			Style: model.PrimaryButtonStyle,
			ID:    fmt.Sprintf("%splay-other", baseID),
		},
	}

	err = ctx.ReplyComplex(&model.ComplexMessage{
		Embeds:     []*model.Embed{embed},
		Components: components,
	})

	if err != nil {
		logger.Errorf("Cannot send message: %v", err)
		return nil
	}

	select {
	case <-time.After(firstResultTimeout):
		embed := buildPlayableInfoEmbed(firstResult.ToPlayable()[0], vc, ctx.AuthorID).
			WithTitle(t.SelectedResult.Str(searchQuery)).
			WithColor(command.SuccessEmbedColor)

		err = ctx.EditComplex(
			&model.ComplexMessage{
				Embeds:     []*model.Embed{embed},
				Components: buildDisabledComponents(components, 0),
			},
		)
		if err != nil {
			logger.Error(err)
		}
		return firstResult.ToPlayable()
	case result := <-selectionCh:
		if result == nil {
			select {
			case <-time.After(multipleResultsTimeout):
				embed := buildPlayableInfoEmbed(firstResult.ToPlayable()[0], vc, ctx.AuthorID).
					WithTitle(t.SelectedResult.Str(searchQuery)).
					WithColor(command.SuccessEmbedColor)

				err = ctx.EditComplex(
					&model.ComplexMessage{
						Embeds:     []*model.Embed{embed},
						Components: buildDisabledComponents(components, 0),
					},
				)
				if err != nil {
					logger.Error(err)
				}
				return firstResult.ToPlayable()
			case result := <-selectionCh:
				return result.ToPlayable()
			}
		}
		return result.ToPlayable()
	}
}

func buildMultipleResultsMessage(ctx *command.CommandContext, vc *voicer.Voicer, searchQuery string, results []*youtube.SearchResult, selectResult func(*youtube.SearchResult) bool) *model.ComplexMessage {
	t := ctx.T.(*i18n.CommandPlay)

	embed := model.NewEmbed().
		WithColor(command.PendingEmbedColor).
		WithTitle(t.MultipleResultsSelectOne.Str())

	sb := strings.Builder{}
	for i, result := range results {
		sb.WriteString(
			fmt.Sprintf("%s\n",
				t.Entry.Str(
					f.Emojify(i+1),
					result.Title,
					result.Author,
					k.Is(result.Duration == 0, t.Live.Str(":red_circle:"), f.ShortDuration(result.Duration)),
				),
			),
		)
	}
	sb.WriteString(
		fmt.Sprintf(
			"\n\n%s",
			t.IfYouFailToSelect.Str(multipleResultsTimeout/time.Second),
		),
	)

	embed.WithDescription(sb.String())

	ctx.AddCommandDuration(embed)

	components := make([]model.MessageComponent, len(results))

	baseID, err := ctx.RegisterInteractionHandler(
		func(fullID, baseID, userID string) (msg *model.ComplexMessage, done bool) {
			if userID != ctx.AuthorID {
				return nil, false
			}
			indexStr := fullID[len(fullID)-1] - '0'
			index := int(indexStr) - 1
			result := results[index]
			selectResult(result)

			return &model.ComplexMessage{
				Components: buildDisabledComponents(components, index),
				Embeds: []*model.Embed{
					buildPlayableInfoEmbed(result.ToPlayable()[0], vc, ctx.AuthorID).
						WithTitle(t.SelectedResult.Str(searchQuery)).
						WithColor(command.SuccessEmbedColor),
				},
			}, true
		},
	)

	if err != nil {
		logger.Errorf("Error registering interaction handler: %v", err)
		return nil
	}

	for i := range results {
		components[i] = model.ButtonComponent{
			Label: fmt.Sprintf("%d", i+1),
			ID:    fmt.Sprintf("%s-play-%d", baseID, i+1),
			Style: model.PrimaryButtonStyle,
		}
	}

	return &model.ComplexMessage{
		Embeds:     []*model.Embed{embed},
		Components: components,
		Content:    "",
	}
}

func buildDisabledComponents(components []model.MessageComponent, selectedIndex int) []model.MessageComponent {
	disabledComponents := make([]model.MessageComponent, len(components))
	for i, component := range components {
		buttonComponent := component.(model.ButtonComponent)
		buttonComponent.Disabled = true
		if i != selectedIndex {
			buttonComponent.Style = model.SecondaryButtonStyle
		}
		disabledComponents[i] = buttonComponent
	}
	return disabledComponents
}
