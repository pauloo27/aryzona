package radio

import (
	"github.com/Pauloo27/aryzona/internal/providers/youtube"
)

type YouTubeRadio struct {
	BaseRadio
	playable      youtube.YouTubePlayable
	ID, Name, URL string
}

var _ RadioChannel = &YouTubeRadio{}

func newYouTubeRadio(id, name, url string) YouTubeRadio {
	playable, err := youtube.AsPlayable(url)
	if err != nil {
		// FIXME
		panic(err)
	}
	return YouTubeRadio{
		ID:        id,
		Name:      name,
		URL:       url,
		BaseRadio: BaseRadio{},
		playable:  playable,
	}
}

func (r YouTubeRadio) GetID() string {
	return r.ID
}

func (r YouTubeRadio) GetName() string {
	return r.Name
}

func (r YouTubeRadio) GetThumbnailURL() (string, error) {
	return "", nil
}

func (r YouTubeRadio) IsOpus() bool {
	return r.playable.IsOpus()
}

func (r YouTubeRadio) GetDirectURL() (string, error) {
	return r.playable.GetDirectURL()
}

func (r YouTubeRadio) GetFullTitle() (title, artist string) {
	return r.playable.GetFullTitle()
}
