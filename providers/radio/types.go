package radio

import (
	"io"
	"net/http"
	"regexp"

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
		return
	},
}

var RadioCidade = &RadioType{
	Name:    "Radio Cidade (a custom type so I can get the now playing)",
	IsOppus: false,
	GetDirectURL: func(url string) string {
		return url
	},
	GetPlayingNow: func(url, directURL string) (title, artist string) {
		res, err := http.Get("https://np.tritondigital.com/public/nowplaying?mountName=RADIOCIDADEAAC&numberToFetch=1&eventType=track")
		if err != nil {
			return
		}
		defer res.Body.Close()
		body, err := io.ReadAll(res.Body)
		if err != nil {
			return
		}
		// nobody deserves to deal with XML... lets just pretend that we got a string
		// as response and a regex is the way to parse it...
		bodyStr := string(body)
		parseRegex := regexp.MustCompile(`CDATA\[([^\]]+)\]`)
		matches := parseRegex.FindAllStringSubmatch(bodyStr, -1)
		if len(matches) < 4 {
			return
		}
		title = matches[2][1]
		artist = matches[3][1]
		return
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
			return
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
			return
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
