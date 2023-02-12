package audio

import (
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/Pauloo27/aryzona/internal/command"
	"github.com/Pauloo27/aryzona/internal/core/f"
	"github.com/Pauloo27/aryzona/internal/discord/model"
	"github.com/Pauloo27/aryzona/internal/discord/voicer/playable"
	"github.com/Pauloo27/aryzona/internal/i18n"
	"github.com/Pauloo27/aryzona/internal/providers/youtube"
	"github.com/Pauloo27/logger"
	k "github.com/Pauloo27/toolkit"
)

type selectResultFn func(*youtube.SearchResult) bool

func handleMultipleResults(ctx *SearchContext) []playable.Playable {
	t := ctx.T.(*i18n.CommandPlay)

	selectionLock := sync.Mutex{}
	selectionCh := make(chan *youtube.SearchResult)
	selected := false

	var selectResult selectResultFn = func(result *youtube.SearchResult) bool {
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
	firstResult := ctx.Results[0]

	embed := model.NewEmbed().
		WithColor(command.PendingEmbedColor).
		WithTitle(t.MultipleResults.Str()).
		WithDescription(
			t.FirstResultWillPlay.Str(
				firstResult.Title,
				firstResult.Author,
				fmt.Sprintf(
					"https://youtu.be/%s",
					firstResult.ID,
				),
				k.Is(
					firstResult.Duration == 0,
					t.PlayingInfo.DurationLive.Str(":red_circle:"),
					f.ShortDuration(firstResult.Duration),
				),
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
						buildPlayableInfoEmbed(
							PlayableInfo{
								Playable:    firstResult.ToPlayable()[0],
								Voicer:      ctx.Voicer,
								RequesterID: ctx.AuthorID,
								T:           t.PlayingInfo,
							},
						).
							WithTitle(t.SelectedResult.Str(ctx.SearchQuery)).
							WithColor(command.SuccessEmbedColor),
					},
				}, true
			}
			if !selectResult(nil) {
				return nil, false
			}
			return buildMultipleResultsMessage(ctx, selectResult), true
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
		embed := buildPlayableInfoEmbed(
			PlayableInfo{
				Playable:    firstResult.ToPlayable()[0],
				Voicer:      ctx.Voicer,
				RequesterID: ctx.AuthorID,
				T:           t.PlayingInfo,
			},
		).
			WithTitle(t.SelectedResult.Str(ctx.SearchQuery)).
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
				embed := buildPlayableInfoEmbed(
PlayableInfo{
					Playable: firstResult.ToPlayable()[0],
					Voicer: ctx.Voicer,
					RequesterID: ctx.AuthorID,
					T: t.PlayingInfo,
				},
				).
					WithTitle(t.SelectedResult.Str(ctx.SearchQuery)).
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

func buildMultipleResultsMessage(ctx *SearchContext, selectResult selectResultFn) *model.ComplexMessage {
	t := ctx.T.(*i18n.CommandPlay)

	embed := model.NewEmbed().
		WithColor(command.PendingEmbedColor).
		WithTitle(t.MultipleResultsSelectOne.Str())

	sb := strings.Builder{}
	for i, result := range ctx.Results {
		sb.WriteString(
			fmt.Sprintf("%s\n",
				t.Entry.Str(
					f.Emojify(i+1),
					result.Title,
					result.Author,
					k.Is(
						result.Duration == 0,
						t.PlayingInfo.DurationLive.Str(":red_circle:"),
						f.ShortDuration(result.Duration),
					),
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

	components := make([]model.MessageComponent, len(ctx.Results))

	baseID, err := ctx.RegisterInteractionHandler(
		func(fullID, baseID, userID string) (msg *model.ComplexMessage, done bool) {
			if userID != ctx.AuthorID {
				return nil, false
			}
			indexStr := fullID[len(fullID)-1] - '0'
			index := int(indexStr) - 1
			result := ctx.Results[index]
			selectResult(result)

			return &model.ComplexMessage{
				Components: buildDisabledComponents(components, index),
				Embeds: []*model.Embed{
					buildPlayableInfoEmbed(
						PlayableInfo{
							Playable:    result.ToPlayable()[0],
							Voicer:      ctx.Voicer,
							RequesterID: ctx.AuthorID,
							T:           t.PlayingInfo,
						},
					).
						WithTitle(t.SelectedResult.Str(ctx.SearchQuery)).
						WithColor(command.SuccessEmbedColor),
				},
			}, true
		},
	)

	if err != nil {
		logger.Errorf("Error registering interaction handler: %v", err)
		return nil
	}

	for i := range ctx.Results {
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
