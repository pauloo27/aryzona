package audio

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/pauloo27/aryzona/internal/command"
	"github.com/pauloo27/aryzona/internal/core/f"
	"github.com/pauloo27/aryzona/internal/core/routine"
	"github.com/pauloo27/aryzona/internal/discord"
	"github.com/pauloo27/aryzona/internal/discord/model"
	"github.com/pauloo27/aryzona/internal/discord/voicer/playable"
	"github.com/pauloo27/aryzona/internal/i18n"
	"github.com/pauloo27/logger"
	k "github.com/pauloo27/toolkit"
)

const (
	firstResultTimeout     = 5 * time.Second
	multipleResultsTimeout = 30 * time.Second
)

const (
	SelectionTypeCancel SelectionType = "cancel"
	SelectionTypeSelect SelectionType = "select"
)

const (
	ActionCancel    = "cancel"
	ActionPlayNow   = "play_now"
	ActionPlayOther = "play_other"
	ActionSelect    = "select"
)

type (
	SelectionType string

	Selection struct {
		Type  SelectionType
		Index int
	}
)

func handleMultipleResults(ctx *SearchContext) []playable.Playable {
	selectionCh := make(chan Selection)
	ctx.SelectionCh = selectionCh

	promptFirstResultConfirmation(ctx)

	selection := <-selectionCh

	if selection.Type == SelectionTypeCancel {
		return nil
	}

	return ctx.Results[selection.Index].ToPlayable()
}

func promptFirstResultConfirmation(ctx *SearchContext) {
	t := ctx.T.(*i18n.CommandPlay)

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
				int64(firstResultTimeout/time.Second),
			),
		)

	ctx.AddCommandDuration(embed)

	baseID := ctx.MessageID

	components := []model.MessageComponent{
		model.ButtonComponent{
			Label:  t.ConfirmBtn.Str(),
			Style:  model.SuccessButtonStyle,
			ID:     ActionPlayNow,
			BaseID: baseID,
		},
		model.ButtonComponent{
			Label:  t.PlayOtherBtn.Str(),
			Style:  model.PrimaryButtonStyle,
			ID:     ActionPlayOther,
			BaseID: baseID,
		},
		model.ButtonComponent{
			Label:  t.CancelBtn.Str(),
			Style:  model.DangerButtonStyle,
			ID:     ActionCancel,
			BaseID: baseID,
		},
	}

	msg := &model.ComplexMessage{
		Embeds: []*model.Embed{embed},
		ComponentRows: []model.MessageComponentRow{
			{
				Components: components,
			},
		},
	}

	err := ctx.ReplyWithInteraction(
		baseID,
		msg,
		handleConfirmationInteraction(ctx, msg),
	)

	if err != nil {
		ctx.Error(t.SomethingWentWrong.Str())
	}
}

func handleConfirmationInteraction(ctx *SearchContext, msg *model.ComplexMessage) command.InteractionHandler {
	actionCh := make(chan string)

	components := msg.ComponentRows[0].Components

	routine.GoAndRecover(func() {
		select {
		case action := <-actionCh:
			switch action {
			case ActionPlayNow:
				ctx.SelectionCh <- Selection{
					Type:  SelectionTypeSelect,
					Index: 0,
				}
			case ActionCancel:
				ctx.SelectionCh <- Selection{
					Type: SelectionTypeCancel,
				}
			}
		case <-time.After(firstResultTimeout):
			ctx.SelectionCh <- Selection{
				Type:  SelectionTypeSelect,
				Index: 0,
			}
			err := ctx.EditComplex(
				&model.ComplexMessage{
					Embeds: []*model.Embed{buildSelectedResultEmbed(ctx, 0)},
					ComponentRows: []model.MessageComponentRow{
						{
							Components: discord.DisableButtons(components, 0),
						},
					},
				},
			)
			if err != nil {
				logger.Error(err)
			}
			command.RemoveInteractionHandler(ctx.MessageID)
		}
	})

	return func(id, userID, baseID string) (newMessage *model.ComplexMessage, done bool) {
		if userID != ctx.AuthorID {
			return nil, false
		}

		actionCh <- id

		if id == ActionPlayOther {
			return handlePlayOther(ctx), false
		}

		selectedIndex := -1

		for i, component := range components {
			if component.(model.ButtonComponent).ID == id {
				selectedIndex = i
				break
			}
		}

		var embeds []*model.Embed

		if id == ActionPlayNow {
			embeds = []*model.Embed{
				buildSelectedResultEmbed(ctx, 0),
			}
		}

		return &model.ComplexMessage{
			Embeds: embeds,
			ComponentRows: []model.MessageComponentRow{
				{
					Components: discord.DisableButtons(components, selectedIndex),
				},
			},
		}, true
	}
}

func handlePlayOther(ctx *SearchContext) *model.ComplexMessage {
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
			t.IfYouFailToSelect.Str(int64(multipleResultsTimeout/time.Second)),
		),
	)

	embed.WithDescription(sb.String())

	baseID := ctx.MessageID

	cancelButton := model.ButtonComponent{
		Label:  t.CancelBtn.Str(),
		Style:  model.DangerButtonStyle,
		ID:     ActionCancel,
		BaseID: baseID,
	}

	components := make([]model.MessageComponent, len(ctx.Results))
	for i := range ctx.Results {
		components[i] = model.ButtonComponent{
			Label:  fmt.Sprintf("%d", i+1),
			ID:     fmt.Sprintf("%s%d", ActionSelect, i),
			BaseID: baseID,
			Style:  model.PrimaryButtonStyle,
		}
	}

	msg := &model.ComplexMessage{
		Embeds: []*model.Embed{embed},
		ComponentRows: []model.MessageComponentRow{
			{
				Components: components,
			},
			{
				Components: []model.MessageComponent{cancelButton},
			},
		},
	}

	// overwrite original message handler
	ctx.RegisterInteractionHandler(
		baseID,
		handlePlayOtherInteraction(ctx, msg),
	)

	return msg
}

func handlePlayOtherInteraction(ctx *SearchContext, msg *model.ComplexMessage) command.InteractionHandler {
	selectedIndexCh := make(chan int)

	buttonsRow := msg.ComponentRows[0].Components
	cancelRow := msg.ComponentRows[1].Components

	routine.GoAndRecover(func() {
		select {
		case index := <-selectedIndexCh:
			if index >= 0 && index < len(ctx.Results) {
				ctx.SelectionCh <- Selection{
					Type:  SelectionTypeSelect,
					Index: index,
				}
			} else {
				ctx.SelectionCh <- Selection{
					Type: SelectionTypeCancel,
				}
			}
		case <-time.After(multipleResultsTimeout):
			ctx.SelectionCh <- Selection{
				Type:  SelectionTypeSelect,
				Index: 0,
			}
			err := ctx.EditComplex(
				&model.ComplexMessage{
					Embeds: []*model.Embed{buildSelectedResultEmbed(ctx, 0)},
					ComponentRows: []model.MessageComponentRow{
						{
							Components: discord.DisableButtons(buttonsRow, 0),
						},
						{
							Components: discord.DisableButtons(cancelRow, -1),
						},
					},
				},
			)
			if err != nil {
				logger.Error(err)
			}
			command.RemoveInteractionHandler(ctx.MessageID)
		}
	})

	return func(id, userID, baseID string) (newMessage *model.ComplexMessage, done bool) {
		if userID != ctx.AuthorID {
			return nil, false
		}

		var selectedIndex int
		var err error

		if id == ActionCancel {
			selectedIndex = -1
		} else {
			selectedIndex, err = strconv.Atoi(strings.TrimPrefix(id, ActionSelect))
			if err != nil {
				logger.Error(err)
				selectedIndex = -1
			}
		}

		selectedIndexCh <- selectedIndex

		var embeds []*model.Embed

		if selectedIndex != -1 {
			embeds = []*model.Embed{
				buildSelectedResultEmbed(ctx, selectedIndex),
			}
		}

		return &model.ComplexMessage{
			Embeds: embeds,
			ComponentRows: []model.MessageComponentRow{
				{
					Components: discord.DisableButtons(buttonsRow, selectedIndex),
				},
				{
					Components: discord.DisableButtons(cancelRow, k.Is(selectedIndex == -1, 0, -1)),
				},
			},
		}, true
	}
}

func buildSelectedResultEmbed(ctx *SearchContext, index int) *model.Embed {
	t := ctx.T.(*i18n.CommandPlay)

	result := ctx.Results[index].ToPlayable()[0]

	return buildPlayableInfoEmbed(
		PlayableInfo{
			Playable:    result,
			Voicer:      ctx.Voicer,
			RequesterID: ctx.AuthorID,
			T:           t.PlayingInfo,
			Common:      t.Common,
		},
	).
		WithTitle(t.SelectedResult.Str(ctx.SearchQuery)).
		WithColor(command.SuccessEmbedColor)

}
