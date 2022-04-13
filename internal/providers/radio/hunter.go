package radio

import (
	"fmt"
	"regexp"

	"github.com/Pauloo27/aryzona/internal/providers/ffmpeg"
	"github.com/buger/jsonparser"
)

type HunterRadio struct {
	BaseRadio
	ID, Name, URL string
}

var _ RadioChannel = &HunterRadio{}

var (
	hunterIDRe = regexp.MustCompile(`^https:\/\/hls\.hunter\.fm\/(\w+)\/\w+.m3u8$`)
)

func newHunterRadio(id, name, url string) HunterRadio {
	return HunterRadio{
		ID:        id,
		Name:      name,
		URL:       url,
		BaseRadio: BaseRadio{},
	}
}

func (r HunterRadio) GetID() string {
	return r.ID
}

func (r HunterRadio) GetName() string {
	return r.Name
}

func (r HunterRadio) GetShareURL() string {
	matches := hunterIDRe.FindStringSubmatch(r.URL)
	if len(matches) == 0 {
		return ""
	}
	hunterID := matches[1]
	return fmt.Sprintf("https://hunter.fm/%s/", hunterID)
}

func (r HunterRadio) GetThumbnailURL() (string, error) {
	return "", nil
}

func (r HunterRadio) IsOpus() bool {
	return false
}

func (r HunterRadio) GetDirectURL() (string, error) {
	return r.URL, nil
}

func (r HunterRadio) GetFullTitle() (title, artist string) {
	data, err := ffmpeg.GetStreamMetadata(r.URL)
	if err != nil {
		return
	}
	title, _ = jsonparser.GetString(data, "streams", "[0]", "tags", "title")
	artist, _ = jsonparser.GetString(data, "streams", "[0]", "tags", "artist")
	return
}
