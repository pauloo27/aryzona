package radio

import (
	"github.com/Pauloo27/aryzona/logger"
	"github.com/Pauloo27/aryzona/providers/youtube"
)

var M3uRadio = &RadioType{
	Name:    "m3u8 playlist",
	IsOppus: false,
	GetDirectURL: func(url string) string {
		return url
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
}
