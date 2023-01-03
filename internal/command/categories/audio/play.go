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
	"github.com/Pauloo27/aryzona/internal/discord/model"
	"github.com/Pauloo27/aryzona/internal/discord/voicer"
	"github.com/Pauloo27/aryzona/internal/discord/voicer/playable"
	"github.com/Pauloo27/aryzona/internal/providers/youtube"
	"github.com/Pauloo27/aryzona/internal/utils"
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

		if !vc.IsConnected() {
			if err := vc.Connect(); err != nil {
				ctx.Error("Cannot connect to your voice channel")
				return
			}
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

		// what should be added to the queue, since we support playlists...
		// it needs to be a list of playable
		var toPlay []playable.Playable
		// single result, it's info's used for the "Now playing" message
		var displayResult playable.Playable

		if len(results) > 1 {
			toPlay, displayResult = handleMultipleResults(ctx, results)
		} else {
			toPlay, displayResult = handleSingleResult(results[0])
		}

		if toPlay == nil {
			ctx.Error("Something went wrong =(")
			return
		}

		ctx.SuccessEmbed(
			buildPlayableInfoEmbed(displayResult, vc, ctx.AuthorID).WithTitle("Best result for: " + searchQuery),
		)

		utils.Go(func() {
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

func handleSingleResult(result *youtube.SearchResult) (toPlay []playable.Playable, displayResult playable.Playable) {
	toPlay = result.ToPlayable()

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

	return
}

func handleMultipleResults(ctx *command.CommandContext, results []*youtube.SearchResult) (toPlay []playable.Playable, displayResult playable.Playable) {
	selection := make(chan int)

	embed := model.NewEmbed().
		WithColor(command.SuccessEmbedColor).
		WithTitle("Multiple results found, please select one")

	sb := strings.Builder{}
	for i, result := range results {
		sb.WriteString(
			fmt.Sprintf(
				" - %s **%s** from %s (*%s*)\n",
				utils.Emojify(i+1),
				result.Title,
				result.Author,
				utils.ShortDuration(result.Duration),
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

	baseID, err := ctx.RegisterInteractionHandler(handleInteraction(ctx, selection))
	if err != nil {
		ctx.Error("Something went wrong")
		return
	}

	var components []model.MessageComponent

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

	// TODO: send proper response
	// TODO: add a timeout
	err = ctx.ReplyComplex(&model.ComplexMessage{
		Embeds:     []*model.Embed{embed},
		Components: components,
	})

	if err != nil {
		return nil, nil
	}

	result := results[<-selection]
	displayResult = result.ToPlayable()[0]
	toPlay = []playable.Playable{displayResult}
	return
}

func handleInteraction(ctx *command.CommandContext, selection chan<- int) func(string) {
	return func(id string) {
		indexStr := id[len(id)-1] - '0'
		index := int(indexStr)
		selection <- index - 1
	}
}
