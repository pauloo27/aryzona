package youtube

import (
	"errors"
	"regexp"

	"github.com/Pauloo27/searchtube"
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

func GetBestResult(searchQuery string) (url string, isPlaylist bool, err error) {
	if videoRegex.MatchString(searchQuery) {
		return searchQuery, false, nil
	}
	if playlistRegex.MatchString(searchQuery) {
		return searchQuery, true, nil
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
