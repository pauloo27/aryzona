package audio

import (
	"errors"
	"fmt"
	"strings"
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
	"github.com/Pauloo27/aryzona/internal/providers/youtube"
	"github.com/Pauloo27/logger"
)

const (
	maxSearchResults = 5
	selectionTimeout = 10 * time.Second
)

var PlayCommand = command.Command{
	Name: "play", Aliases: []string{"p", "tocar", "yt", "youtube"},
	Description: "Play a video/song from youtube",
	Deferred:    true,
	Validations: []*command.CommandValidation{validations.MustBeOnAValidVoiceChannel},
	Parameters: []*command.CommandParameter{
		{Name: "song", Description: "Search query or URL", Required: true, Type: parameters.ParameterText},
	},
	Handler: func(ctx *command.CommandContext) {
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
				ctx.Error("You are not in the right voice channel")
				return
			}
		}

		searchQuery := ctx.Args[0].(string)
		results, err := youtube.SearchFor(searchQuery, maxSearchResults)
		if err != nil || len(results) == 0 {
			ctx.Error("Cannot search for this song")
			logger.Warnf("Error searching for %s: %v", searchQuery, err)
			return
		}

		if connErrCh != nil {
			if err := <-(*connErrCh); err != nil {
				ctx.Error("Cannot connect to your voice channel")
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
			ctx.Error("Something went wrong =(")
			return
		}

		routine.Go(func() {
			if err = vc.AppendManyToQueue(ctx.AuthorID, toPlay...); err != nil {
				if errors.Is(err, dca.ErrVoiceConnectionClosed) {
					return
				}
				ctx.Errorf("Cannot play stuff: %v", err)
				return
			}
		})
	},
}

func handleSingleResult(ctx *command.CommandContext, vc *voicer.Voicer, searchQuery string, result *youtube.SearchResult) (toPlay []playable.Playable) {
	toPlay = result.ToPlayable()

	var displayResult playable.Playable

	if result.IsPlaylist() {
		displayResult = playable.DummyPlayable{
			Name:     "YouTube Playlist",
			Artist:   result.Author,
			Title:    result.Title,
			Duration: result.Duration,
		}
	} else {
		displayResult = toPlay[0]
	}

	ctx.SuccessEmbed(
		buildPlayableInfoEmbed(displayResult, vc, ctx.AuthorID).WithTitle("Best result for: " + searchQuery),
	)

	return
}

func handleMultipleResults(ctx *command.CommandContext, vc *voicer.Voicer, searchQuery string, results []*youtube.SearchResult) []playable.Playable {
	selection := make(chan playable.Playable)
	var components []model.MessageComponent

	embed := model.NewEmbed().
		WithColor(command.PendingEmbedColor).
		WithTitle("Multiple results found, please select one")

	sb := strings.Builder{}
	for i, result := range results {
		sb.WriteString(
			fmt.Sprintf(
				" - %s **%s** from %s (*%s*)\n",
				f.Emojify(i+1),
				result.Title,
				result.Author,
				f.ShortDuration(result.Duration),
			),
		)
	}
	sb.WriteString(
		fmt.Sprintf(
			"\n\n**If you fail to select one in %d seconds, the first result will be selected**",
			selectionTimeout/time.Second,
		),
	)

	embed.WithDescription(sb.String())

	ctx.AddCommandDuration(embed)

	baseID, err := ctx.RegisterInteractionHandler(
		func(id string) *model.ComplexMessage {
			indexStr := id[len(id)-1] - '0'
			index := int(indexStr) - 1
			result := results[index].ToPlayable()[0]
			selection <- result

			return &model.ComplexMessage{
				Components: buildDisabledComponents(components, index),
				Embeds: []*model.Embed{
					buildPlayableInfoEmbed(result, vc, ctx.AuthorID).
						WithTitle("Selected result for: " + searchQuery).
						WithColor(command.SuccessEmbedColor),
				},
			}
		},
	)

	if err != nil {
		ctx.Error("Something went wrong")
		return nil
	}

	for i := range results {
		components = append(
			components,
			model.ButtonComponent{
				Label: fmt.Sprintf("%d", i+1),
				ID:    fmt.Sprintf("%s-play-%d", baseID, i+1),
				Style: model.PrimaryButtonStyle,
			},
		)
	}

	err = ctx.ReplyComplex(&model.ComplexMessage{
		Embeds:     []*model.Embed{embed},
		Components: components,
	})

	if err != nil {
		return nil
	}

	select {
	case result := <-selection:
		return []playable.Playable{result}
	case <-time.After(selectionTimeout):
		result := results[0].ToPlayable()[0]
		err = ctx.EditComplex(&model.ComplexMessage{
			Components: buildDisabledComponents(components, 0),
			Embeds: []*model.Embed{
				buildPlayableInfoEmbed(result, vc, ctx.AuthorID).
					WithTitle("Selected result for: " + searchQuery).
					WithColor(command.SuccessEmbedColor),
			},
		})
		if err != nil {
			logger.Errorf("Error editing message: %v", err)
		}
		return []playable.Playable{result}
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
