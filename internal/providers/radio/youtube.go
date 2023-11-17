package radio

import (
	"log/slog"

	"github.com/lmittmann/tint"
	"github.com/pauloo27/aryzona/internal/discord/voicer/playable"
	"github.com/pauloo27/aryzona/internal/providers/youtube"
)

type YouTubeRadio struct {
	BaseRadio
	ID, Name, URL string
}

var _ RadioChannel = &YouTubeRadio{}

func newYouTubeRadio(id, name, url string) RadioChannel {
	return YouTubeRadio{
		ID:        id,
		Name:      name,
		URL:       url,
		BaseRadio: BaseRadio{},
	}
}

func (r YouTubeRadio) GetID() string {
	return r.ID
}

func (r YouTubeRadio) GetName() string {
	return r.Name
}

func (r YouTubeRadio) GetShareURL() string {
	return r.GetPlayable().GetShareURL()
}

func (r YouTubeRadio) GetThumbnailURL() (string, error) {
	return "", nil
}

func (r YouTubeRadio) IsOpus() bool {
	return r.GetPlayable().IsOpus()
}

func (r YouTubeRadio) GetDirectURL() (string, error) {
	return r.GetPlayable().GetDirectURL()
}

func (r YouTubeRadio) GetFullTitle() (title, artist string) {
	return r.GetPlayable().GetFullTitle()
}

func (r YouTubeRadio) GetPlayable() playable.Playable {
	vid, err := youtube.GetVideo(r.URL)
	if err != nil {
		slog.Error("Error getting video", tint.Err(err))
	}
	return vid
}
