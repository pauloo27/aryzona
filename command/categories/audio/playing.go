package audio

import (
	"github.com/Pauloo27/aryzona/command"
	"github.com/Pauloo27/aryzona/discord/voicer"
	"github.com/Pauloo27/aryzona/utils"
)

var PlayingCommand = command.Command{
	Name: "playing", Aliases: []string{"np", "nowplaying", "tocando"},
	Description: "Show what is playing now",
	Handler: func(ctx *command.CommandContext) {
		vc := voicer.GetExistingVoicerForGuild(ctx.Message.GuildID)
		if vc == nil {
			ctx.Error("Bot is not connect to a voice channel")
			return
		}
		playable := *(vc.Playing)

		title, artist := playable.GetFullTitle()

		embedBuilder := utils.NewEmbedBuilder().
			Title("Now playing: "+playable.GetName()).
			Field("Title", title)

		if artist != "" {
			embedBuilder.Field("Artist", artist)
		}

		embedBuilder.Field("Requested by", utils.AsMention(*vc.UserID))

		ctx.SuccessEmbed(
			embedBuilder.Build(),
		)
	},
}
