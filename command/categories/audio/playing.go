package audio

import (
	"strings"
	"time"

	"github.com/Pauloo27/aryzona/command"
	"github.com/Pauloo27/aryzona/discord"
	"github.com/Pauloo27/aryzona/discord/voicer"
	"github.com/Pauloo27/aryzona/utils"
)

var PlayingCommand = command.Command{
	Name: "playing", Aliases: []string{"np", "nowplaying", "tocando"},
	Description: "Show what is playing now",
	Handler: func(ctx *command.CommandContext) {
		vc := voicer.GetExistingVoicerForGuild(ctx.GuildID)
		if vc == nil {
			ctx.Error("Bot is not connect to a voice channel")
			return
		}
		playable := vc.Playing()

		if playable == nil {
			ctx.Error("Nothing playing...")
			return
		}

		title, artist := playable.GetFullTitle()

		embed := discord.NewEmbed().
			WithTitle("Now playing: "+playable.GetName()).
			WithField("Title", title)

		if artist != "" {
			embed.WithField("Artist", artist)
		}

		embed.WithField("Requested by", discord.AsMention(*vc.UserID))
		if playable.IsLive() {
			embed.WithField("Duration", "**🔴 LIVE**")
		} else {
			position, err := vc.GetPosition()
			if err == nil {
				embed.WithField("Position", position.Truncate(time.Second).String())
			}
			duration, err := playable.GetDuration()
			if err == nil {
				embed.WithField("Duration", duration.String())
			}
		}

		if vc.Queue.Size() > 1 {
			sb := strings.Builder{}
			next := vc.Queue.All()[1:]
			limit := len(next)
			if len(next) > 10 {
				limit = 10
			}
			for _, item := range next[:limit] {
				title, artist := item.GetFullTitle()
				if artist == "" {
					sb.WriteString(utils.Fmt("  -> %s\n", title))
				} else {
					sb.WriteString(utils.Fmt("  -> %s - %s\n", artist, title))
				}
			}
			if len(next) > 10 {
				sb.WriteString("_... and more ..._")
			}
			embed.WithField("**Coming next:**", sb.String())
		}

		ctx.SuccessEmbed(embed)
	},
}
