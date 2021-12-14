package youtube

import (
	"github.com/kkdai/youtube/v2"
)

var (
	defaultClient = youtube.Client{}
)

func GetVideoID(videoURL string) string {
	matches := videoRegex.FindAllStringSubmatch(videoURL, -1)
	if len(matches) == 1 {
		match := matches[0]
		if len(match) == 2 {
			return match[1]
		}
	}
	return ""
}
