package audio

import (
	"fmt"
	"time"

	"github.com/Pauloo27/aryzona/internal/command"
	"github.com/Pauloo27/aryzona/internal/core/f"
	"github.com/Pauloo27/aryzona/internal/discord"
	"github.com/Pauloo27/aryzona/internal/discord/model"
	"github.com/Pauloo27/aryzona/internal/discord/voicer"
	"github.com/Pauloo27/aryzona/internal/discord/voicer/playable"
	"github.com/Pauloo27/aryzona/internal/i18n"
)

func buildPlayableInfoEmbed(playable playable.Playable, vc *voicer.Voicer, requesterID string, t *i18n.PlayingInfo) *model.Embed {
	title, artist := playable.GetFullTitle()

	embed := model.NewEmbed().
		WithField(t.SongTitle.Str(), title)

	shareURL := playable.GetShareURL()
	if shareURL != "" {
		embed.WithURL(shareURL)
	}

	if artist != "" {
		embed.WithFieldInline(t.Artist.Str(), artist)
	}

	embed.WithFieldInline(t.Source.Str(), playable.GetName())

	thumbnailURL, err := playable.GetThumbnailURL()
	if err == nil && thumbnailURL != "" {
		embed.WithThumbnail(thumbnailURL)
	}

	eta := calcETA(playable, vc)

	if eta == -1 {
		embed.WithFieldInline(
			t.ETAKey.Str(),
			t.ETANever.Str(),
		)
	} else if eta != 0 {
		embed.WithFieldInline(
			t.ETAKey.Str(),
			t.ETAValue.Str(
				f.DurationAsDetailedDiffText(eta),
			),
		)
	}

	if playable.IsLive() {
		embed.WithFieldInline(t.DurationKey.Str(), t.DurationLive.Str(":red_circle:"))
	} else {
		position, posErr := vc.GetPosition()
		duration, durErr := playable.GetDuration()

		if vc.Playing() != nil && playable == vc.Playing().Playable && posErr == nil && durErr == nil {
			embed.WithField(t.DurationKey.Str(), fmt.Sprintf("%s/%s",
				f.ShortDuration(position),
				f.ShortDuration(duration),
			))
		} else if durErr == nil {
			embed.WithField(t.DurationKey.Str(), f.ShortDuration(duration))
		} else if posErr == nil {
			embed.WithField(t.Position.Str(), f.ShortDuration(position))
		}
	}

	if requesterID != "" {
		embed.WithFieldInline(t.RequestedBy.Str(), discord.AsMention(requesterID))
	}

	if vc != nil && vc.IsPaused() {
		embed.WithField(
			t.Warning.Str(),
			t.SongPausedWarning.Str(command.Prefix),
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
