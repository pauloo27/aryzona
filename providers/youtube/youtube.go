package youtube

import (
	"strings"

	"github.com/kkdai/youtube/v2"
)

var defaultClient = youtube.Client{}

func GetVideoID(videoURL string) string {
	// probably not the best way to do it tho
	return strings.Split(videoURL, "=")[1]
}
