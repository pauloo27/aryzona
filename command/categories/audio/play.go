package audio

import (
	"github.com/Pauloo27/aryzona/audio/dca"
	"github.com/Pauloo27/aryzona/command"
	"github.com/Pauloo27/aryzona/discord/voicer"
	"github.com/Pauloo27/aryzona/providers/youtube"
	"github.com/Pauloo27/aryzona/utils"
	"github.com/Pauloo27/logger"
)

var PlayCommand = command.Command{
	Name: "play", Aliases: []string{"p", "tocar", "yt", "youtube"},
	Arguments: []*command.CommandArgument{
		{Name: "Search query", Required: true, Type: command.ArgumentText},
	},
	Handler: func(ctx *command.CommandContext) {
		vc, err := voicer.NewVoicerForUser(ctx.Message.Author.ID, ctx.Message.GuildID)
		if err != nil {
			ctx.Error("Cannot create voicer")
			return
		}
		if !vc.CanConnect() {
			ctx.Error("You are not in a voice channel")
			return
		}
		if vc.IsPlaying() {
			ctx.Error("Already playing something")
			return
		}
		if err = vc.Connect(); err != nil {
			ctx.Error("Cannot connect to your voice channel")
			return
		}
		searchQuery := ctx.Args[0].(string)
		result, err := youtube.GetBestResult(searchQuery)
		if err != nil {
			ctx.Error("Cannot find what you are looking for")
			return
		}
		playable, err := youtube.AsPlayable(result.URL)
		if err != nil {
			ctx.Error("Something went wrong when getting the video to play")
			return
		}
		ctx.SuccessEmbed(
			utils.NewEmbedBuilder().
				Title(utils.Fmt("Best result for %s:", searchQuery)).
				Thumbnail(result.Thumbnail).
				Field("Title", result.Title).
				Field("Uploader", result.Uploader).
				Build(),
		)
		go func() {
			if err = vc.Play(playable); err != nil {
				if is, vErr := utils.IsErrore(err); is {
					if vErr.ID == dca.ErrVoiceConnectionClosed.ID {
						return
					}
					ctx.Error(vErr.Message)
				} else {
					ctx.Error("Cannot play stuff")
					logger.Error(err)
				}
				return
			}
		}()
	},
}
