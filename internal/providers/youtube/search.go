package youtube

import (
	"errors"
	"fmt"
	"net/url"
	"os"

	"github.com/Pauloo27/aryzona/internal/utils"
	"github.com/Pauloo27/searchtube"
	"github.com/buger/jsonparser"
)

func searchWithAPI(searchQuery string) (id string, err error) {
	apiKey := os.Getenv("DC_BOT_YOUTUBE_API_KEY")
	uri := fmt.Sprintf(
		"https://www.googleapis.com/youtube/v3/search?q=%s&key=%s&maxResults=1&type=video",
		url.QueryEscape(searchQuery), url.QueryEscape(apiKey),
	)
	buf, err := utils.Get(uri)
	if err != nil {
		return "", err
	}
	return jsonparser.GetString(buf, "items", "[0]", "id", "videoId")
}

func GetBestResult(searchQuery string) (url string, isPlaylist bool, err error) {
	pl, err := defaultClient.GetPlaylist(searchQuery)

	// if the search query is a personal mix, it returns no errors, but an
	// empty playlist, so we checking the size of it
	if err == nil && len(pl.Videos) != 0 {
		return searchQuery, true, nil
	}

	_, err = defaultClient.GetVideo(searchQuery)
	if err == nil {
		return searchQuery, false, nil
	}

	id, err := searchWithAPI(searchQuery)
	if err == nil {
		return fmt.Sprintf("https://youtu.be/%s", id), false, err
	}

	results, err := searchtube.Search(searchQuery, 1)
	if err == nil || results == nil || len(results) == 0 {
		return "", false, errors.New("no results found")
	}
	return results[0].URL, false, nil
}
