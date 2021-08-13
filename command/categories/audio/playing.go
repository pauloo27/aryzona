package audio

import (
	"github.com/Pauloo27/aryzona/command"
	"github.com/Pauloo27/aryzona/discord"
	"github.com/Pauloo27/aryzona/discord/voicer"
	"github.com/Pauloo27/aryzona/utils"
	"github.com/Pauloo27/logger"
)

var PlayingCommand = command.Command{
	Name: "playing", Aliases: []string{"np", "nowplaying", "tocando"},
	Handler: func(ctx *command.CommandContext) {
		voicer := voicer.GetExistingVoicerForGuild(ctx.Message.GuildID)
		if voicer == nil {
			ctx.Error("Bot is not connect to a voice channel")
			return
		}
		playable := *(voicer.Playing)

		title, artist := playable.GetFullTitle()

		embedBuilder := utils.NewEmbedBuilder().
			Title("Now playing: "+playable.GetName()).
			Field("Title", title)

		if artist != "" {
			embedBuilder.Field("Artist", artist)
		}

		var memberName string
		member, err := discord.Session.State.Member(ctx.Message.GuildID, *(voicer.UserID))
		if err != nil {
			logger.Error(err)
			memberName = *(voicer.UserID)
		} else {
			memberName = member.Mention()
		}
		embedBuilder.Field("Requested by", memberName)

		ctx.SuccessEmbed(
			embedBuilder.Build(),
		)
	},
}
