package radio

import (
	"html"

	"github.com/Pauloo27/aryzona/internal/providers/ffmpeg"
	"github.com/buger/jsonparser"
)

type LainchanRadio struct {
	BaseRadio
	ID, Name, URL string
}

var _ RadioChannel = &LainchanRadio{}

func newLainchanRadio(id, name, url string) LainchanRadio {
	return LainchanRadio{
		ID:        id,
		Name:      name,
		URL:       url,
		BaseRadio: BaseRadio{},
	}
}

func (r LainchanRadio) GetID() string {
	return r.ID
}

func (r LainchanRadio) GetName() string {
	return r.Name
}

func (r LainchanRadio) GetThumbnailURL() (string, error) {
	return "", nil
}

func (r LainchanRadio) IsOpus() bool {
	return false
}

func (r LainchanRadio) GetDirectURL() (string, error) {
	return r.URL, nil
}

func (r LainchanRadio) GetFullTitle() (title, artist string) {
	data, err := ffmpeg.GetStreamMetadata(r.URL)
	if err != nil {
		return
	}
	title, _ = jsonparser.GetString(data, "format", "tags", "StreamTitle")
	title = html.UnescapeString(title)
	return
}
