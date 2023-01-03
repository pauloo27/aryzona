package audio

import (
	"fmt"
	"time"

	"github.com/Pauloo27/aryzona/internal/command"
	"github.com/Pauloo27/aryzona/internal/discord"
	"github.com/Pauloo27/aryzona/internal/discord/model"
	"github.com/Pauloo27/aryzona/internal/discord/voicer"
	"github.com/Pauloo27/aryzona/internal/discord/voicer/playable"
	"github.com/Pauloo27/aryzona/internal/utils"
)

func buildPlayableInfoEmbed(playable playable.Playable, vc *voicer.Voicer, requesterID string) *model.Embed {
	title, artist := playable.GetFullTitle()

	embed := model.NewEmbed().
		WithField("Title", title)

	shareURL := playable.GetShareURL()
	if shareURL != "" {
		embed.WithURL(shareURL)
	}

	if artist != "" {
		embed.WithFieldInline("Artist", artist)
	}

	embed.WithFieldInline("Source", playable.GetName())

	thumbnailURL, err := playable.GetThumbnailURL()
	if err == nil && thumbnailURL != "" {
		embed.WithThumbnail(thumbnailURL)
	}

	eta := calcETA(playable, vc)

	if eta == -1 {
		embed.WithFieldInline(
			"Will play",
			"_Never_",
		)
	} else if eta != 0 {
		embed.WithFieldInline(
			"Will play",
			fmt.Sprintf(
				"in %s",
				utils.DurationAsDetailedDiffText(eta),
			),
		)
	}

	if playable.IsLive() {
		embed.WithFieldInline("Duration", "**🔴 LIVE**")
	} else {
		position, posErr := vc.GetPosition()
		duration, durErr := playable.GetDuration()

		if vc.Playing() != nil && playable == vc.Playing().Playable && posErr == nil && durErr == nil {
			embed.WithField("Duration", fmt.Sprintf("%s/%s",
				utils.ShortDuration(position),
				utils.ShortDuration(duration),
			))
		} else if durErr == nil {
			embed.WithField("Duration", utils.ShortDuration(duration))
		} else if posErr == nil {
			embed.WithField("Position", utils.ShortDuration(position))
		}
	}

	if requesterID != "" {
		embed.WithFieldInline("Requested by", discord.AsMention(requesterID))
	}

	if vc != nil && vc.IsPaused() {
		embed.WithField(
			"Warning",
			fmt.Sprintf("Song is **paused**, use **%sresume**", command.Prefix),
		)
	}

	return embed
}

func calcETA(playable playable.Playable, vc *voicer.Voicer) time.Duration {
	if vc == nil {
		return 0
	}

	var eta time.Duration
	for i, entry := range vc.Queue.All() {
		if entry.Playable == playable {
			break
		}
		if entry.Playable.IsLive() {
			return -1
		}

		duration, _ := entry.Playable.GetDuration()
		if i == 0 {
			position, _ := vc.GetPosition()
			duration -= position
		}
		eta += duration
	}
	return eta
}
