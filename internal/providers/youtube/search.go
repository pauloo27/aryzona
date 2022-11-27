package youtube

import (
	"errors"
	"fmt"
	"net/url"

	"github.com/Pauloo27/aryzona/internal/config"
	"github.com/Pauloo27/aryzona/internal/utils"
	"github.com/Pauloo27/logger"
	"github.com/Pauloo27/searchtube"
	"github.com/tidwall/gjson"
)

func searchWithAPI(searchQuery string) (id string, err error) {
	apiKey := config.Config.YoutubeAPIKey
	uri := fmt.Sprintf(
		"https://www.googleapis.com/youtube/v3/search?q=%s&key=%s&maxResults=1&type=video",
		url.QueryEscape(searchQuery), url.QueryEscape(apiKey),
	)
	buf, err := utils.Get(uri)
	if err != nil {
		return "", err
	}
	return gjson.GetBytes(buf, "items.0.id.videoId").String(), nil
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
	logger.Debug("Err", err)

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
