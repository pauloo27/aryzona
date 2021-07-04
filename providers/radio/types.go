package radio

import (
	"github.com/Pauloo27/aryzona/logger"
	"github.com/Pauloo27/aryzona/providers/youtube"
	"github.com/Pauloo27/aryzona/utils"
	"github.com/buger/jsonparser"
)

var GenericRadio = &RadioType{
	Name:    "Generic radio with simple streaming URL",
	IsOppus: false,
	GetDirectURL: func(url string) string {
		return url
	},
	GetPlayingNow: func(url, directURL string) (title, artist string) {
		return "", ""
	},
}

var LainchanRadio = &RadioType{
	Name:    "Lainchan m3u8 playlist",
	IsOppus: false,
	GetDirectURL: func(url string) string {
		return url
	},
	GetPlayingNow: func(url, directURL string) (title, artist string) {
		data, err := utils.GetStreamMetadata(directURL)
		if err != nil {
			return "", ""
		}
		title, _ = jsonparser.GetString(data, "format", "tags", "StreamTitle")

		return
	},
}

var HunterFM = &RadioType{
	Name:    "HunterFM m3u8 playlist",
	IsOppus: false,
	GetDirectURL: func(url string) string {
		return url
	},
	GetPlayingNow: func(url, directURL string) (title, artist string) {
		data, err := utils.GetStreamMetadata(directURL)
		if err != nil {
			return "", ""
		}
		title, _ = jsonparser.GetString(data, "streams", "[0]", "tags", "title")
		artist, _ = jsonparser.GetString(data, "streams", "[0]", "tags", "artist")

		return
	},
}

var YTLive = &RadioType{
	Name:    "YouTube live",
	IsOppus: false,
	GetDirectURL: func(url string) string {
		url, err := youtube.GetLiveURL(url)

		if err != nil {
			logger.Errorf("%s", err)
		}

		return url
	},
	GetPlayingNow: func(url, directURL string) (title, artist string) {
		return "YouTube live", ""
	},
}
