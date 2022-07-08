package audio

import (
	"errors"
	"fmt"

	"github.com/Pauloo27/aryzona/internal/audio/dca"
	"github.com/Pauloo27/aryzona/internal/command"
	"github.com/Pauloo27/aryzona/internal/command/parameters"
	"github.com/Pauloo27/aryzona/internal/command/validations"
	"github.com/Pauloo27/aryzona/internal/discord/voicer"
	"github.com/Pauloo27/aryzona/internal/discord/voicer/playable"
	"github.com/Pauloo27/aryzona/internal/providers/youtube"
	"github.com/Pauloo27/aryzona/internal/utils"
)

var PlayCommand = command.Command{
	Name: "play", Aliases: []string{"p", "tocar", "yt", "youtube"},
	Description: "Play a video/song from u2b",
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
		}

		searchQuery := ctx.Args[0].(string)
		resultURL, isPlaylist, err := youtube.GetBestResult(searchQuery)
		if err != nil {
			ctx.Error("Cannot find what you are looking for")
			return
		}

		var playlist youtube.YouTubePlayablePlaylist
		var result playable.Playable

		if isPlaylist {
			playlist, err = youtube.GetPlaylist(resultURL)
			if err != nil {
				ctx.Error("Cannot find what you are looking for")
				return
			}
			result = playable.DummyPlayable{
				Name:     "YouTube Playlist",
				Artist:   playlist.Author,
				Title:    playlist.Title,
				Duration: playlist.Duration,
			}
		} else {
			result, err = youtube.AsPlayable(resultURL)
			if err != nil {
				ctx.Error(fmt.Sprintf("Something went wrong when getting the video to play: %v", err))
				return
			}
		}

		embed := buildPlayableInfoEmbed(result, vc).WithTitle("Best result for: " + searchQuery)
		ctx.SuccessEmbed(embed)

		var vidsToAppend []playable.Playable
		if isPlaylist {
			vidsToAppend = playlist.Videos
		} else {
			vidsToAppend = []playable.Playable{result}
		}

		utils.Go(func() {
			if err = vc.AppendManyToQueue(vidsToAppend...); err != nil {
				if errors.Is(err, dca.ErrVoiceConnectionClosed) {
					return
				}
				ctx.Error(fmt.Sprintf("Cannot play stuff: %v", err))
				return
			}
		})
	},
}
