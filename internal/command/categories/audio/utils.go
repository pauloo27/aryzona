package audio

import (
	"fmt"
	"time"

	"github.com/pauloo27/aryzona/internal/command"
	"github.com/pauloo27/aryzona/internal/core/f"
	"github.com/pauloo27/aryzona/internal/discord"
	"github.com/pauloo27/aryzona/internal/discord/model"
	"github.com/pauloo27/aryzona/internal/discord/voicer"
	"github.com/pauloo27/aryzona/internal/discord/voicer/playable"
	"github.com/pauloo27/aryzona/internal/i18n"
)

type PlayableInfo struct {
	Playable    playable.Playable
	Voicer      *voicer.Voicer
	RequesterID string
	T           *i18n.PlayingInfo
	Common      *i18n.Common
}

func buildPlayableInfoEmbed(info PlayableInfo) *model.Embed {
	playable := info.Playable
	vc := info.Voicer
	requesterID := info.RequesterID
	t := info.T

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
				f.DurationAsDetailedDiffText(eta, info.Common),
			),
		)
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
