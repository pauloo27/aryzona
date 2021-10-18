package youtube

import (
	"errors"
	"regexp"

	"github.com/Pauloo27/searchtube"
)

// FIXME: this regex allows things like `youtu.be.com`, which is not valid...
// with a "right" string payload (and the hard part, the proper domain), the
// regex might say that the link is from  youtube when it's not, leading the
// bot to connect in a "invalid" server which can lead to "IP leaking" or
// whatever...
var videoRegex = regexp.MustCompile(`(?:https?:\/\/)?(?:www\.)?youtu\.?be(?:\.com)?\/?.*(?:watch|embed)?(?:.*v=|v\/|\/)([\w\-_]+)\&?`)

func GetBestResult(searchQuery string) (string, error) {
	if videoRegex.MatchString(searchQuery) {
		return searchQuery, nil
	}
	results, err := searchtube.Search(searchQuery, 1)
	if err != nil {
		return "", err
	}
	if len(results) == 0 {
		return "", errors.New("no results found")
	}
	return results[0].URL, nil
}
