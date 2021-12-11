package youtube

import (
	"errors"
	"fmt"
	"net/url"
	"os"
	"regexp"

	"github.com/Pauloo27/aryzona/utils"
	"github.com/Pauloo27/searchtube"
	"github.com/buger/jsonparser"
)

var (
	// FIXME: this regex allows things like `youtu.be.com`, which is not valid...
	// with a "right" string payload (and the hard part, the proper domain), the
	// regex might say that the link is from  youtube when it's not, leading the
	// bot to connect in a "invalid" server which can lead to "IP leaking" or
	// whatever...
	videoRegex = regexp.MustCompile(`^(?:https?:\/\/)?(?:www\.)?youtu\.?be(?:\.com)?\/?.*(?:watch|embed)?(?:.*v=|v\/|\/)([\w\-_]+)\&?$`)

	playlistRegex = regexp.MustCompile(`^https?:\/\/(www.youtube.com|youtube.com)\/playlist(.*)$`)
)

func getBestUsingFromAPI(searchQuery string) (id string, err error) {
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
	if videoRegex.MatchString(searchQuery) {
		return searchQuery, false, nil
	}
	if playlistRegex.MatchString(searchQuery) {
		return searchQuery, true, nil
	}

	id, err := getBestUsingFromAPI(searchQuery)
	if err == nil {
		return fmt.Sprintf("https://youtu.be/%s", id), false, err
	}

	results, err := searchtube.Search(searchQuery, 1)
	if err != nil {
		return "", false, err
	}
	if len(results) == 0 {
		return "", false, errors.New("no results found")
	}
	return results[0].URL, false, nil
}
