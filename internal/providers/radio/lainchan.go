package radio

import (
	"html"

	"github.com/Pauloo27/aryzona/internal/utils"
	"github.com/tidwall/gjson"
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

func (r LainchanRadio) GetShareURL() string {
	return "https://lainchan.org/radio"
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
	data, err := utils.Get("https://lainon.life/radio/status-json.xsl")
	if err != nil {
		return
	}

	result := gjson.GetBytes(data, "icestats.source")
	result.ForEach(func(key, value gjson.Result) bool {
		if value.Get("listenurl").String() == r.URL {
			title = html.UnescapeString(
				value.Get("title").String(),
			)
			artist = html.UnescapeString(
				value.Get("artist").String(),
			)
			return false
		}
		return true
	})
	return
}
