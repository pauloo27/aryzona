package radio

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/Pauloo27/aryzona/internal/core/h"
	"github.com/tidwall/gjson"
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
	matches := hunterIDRe.FindStringSubmatch(r.URL)
	if len(matches) == 0 {
		return
	}
	hunterID := matches[1]

	data, err := h.Get("https://api.hunter.fm/stations/")
	if err != nil {
		return
	}

	result := gjson.ParseBytes(data)
	result.ForEach(func(key, value gjson.Result) bool {
		if value.Get("url").String() == hunterID {
			title = value.Get("live.now.name").String()

			artists := value.Get("live.now.singers").Array()
			feats := value.Get("live.now.feats").Array()

			sb := strings.Builder{}
			for i, a := range artists {
				if i != 0 {
					sb.WriteString(", ")
				}
				sb.WriteString(a.String())
			}

			if len(feats) > 0 {
				sb.WriteString(" feat. ")
				for i, f := range feats {
					if i != 0 {
						sb.WriteString(", ")
					}
					sb.WriteString(f.String())
				}
			}

			artist = sb.String()
			return false
		}
		return true
	})

	return
}
