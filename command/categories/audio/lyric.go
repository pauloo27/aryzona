package audio

import (
	"github.com/Pauloo27/aryzona/command"
	"github.com/Pauloo27/aryzona/discord"
	"github.com/Pauloo27/aryzona/discord/voicer"
	"github.com/Pauloo27/aryzona/utils"
	"github.com/Pauloo27/logger"
	"github.com/Pauloo27/lyric"
)

var LyricCommand = command.Command{
	Name: "lyric", Aliases: []string{"ly", "letra", "letras", "lyrics"},
	Description: "Show lyric from a song",
	Arguments: []*command.CommandArgument{
		{Name: "song", Description: "Search terms", Type: command.ArgumentText, Required: false},
	},
	Handler: func(ctx *command.CommandContext) {
		var searchTerms string

		if len(ctx.Args) != 0 {
			searchTerms = ctx.Args[0].(string)
		} else {
			vc := voicer.GetExistingVoicerForGuild(ctx.GuildID)
			if vc == nil {
				ctx.Error("Bot is not connect to a voice channel, you can pass the song title")
				return
			}
			playable := vc.Playing()
			if playable == nil {
				ctx.Error("Nothing playing...")
				return
			}

			title, artist := playable.GetFullTitle()
			searchTerms = utils.Fmt("%s %s", title, artist)
		}

		result, err := lyric.SearchDDG(searchTerms)
		if err != nil {
			ctx.ErrorEmbed(
				discord.NewEmbed().
					WithDescription("No results found for " + searchTerms),
			)
			logger.Error(err)
			return
		}
		ctx.Success(result)
	},
}
