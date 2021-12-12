package audio

import (
	"time"

	"github.com/Pauloo27/aryzona/discord"
	"github.com/Pauloo27/aryzona/discord/voicer"
	"github.com/Pauloo27/aryzona/discord/voicer/playable"
)

func buildPlayableInfoEmbed(playable playable.Playable, vc *voicer.Voicer) *discord.Embed {
	title, artist := playable.GetFullTitle()

	embed := discord.NewEmbed().
		WithField("Title", title)

	if artist != "" {
		embed.WithFieldInline("Artist", artist)
	}

	thumbnailURL, err := playable.GetThumbnailURL()
	if err == nil && thumbnailURL != "" {
		embed.WithThumbnail(thumbnailURL)
	}

	if playable.IsLive() {
		embed.WithFieldInline("Duration", "**🔴 LIVE**")
	} else {
		if vc != nil {
			if position, err := vc.GetPosition(); err == nil {
				embed.WithFieldInline("Position", position.Truncate(time.Second).String())
			}
		}
		duration, err := playable.GetDuration()
		if err == nil {
			embed.WithFieldInline("Duration", duration.String())
		}
	}

	return embed
}
