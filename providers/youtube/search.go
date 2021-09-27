package youtube

import (
	"errors"
	"regexp"

	"github.com/Pauloo27/searchtube"
)

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
