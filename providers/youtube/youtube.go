package youtube

import (
	"strings"

	"github.com/kkdai/youtube/v2"
)

func GetMediaURL(url string) (string, error) {
	videoID := strings.Split(url, "=")[1]
	client := youtube.Client{}

	video, err := client.GetVideo(videoID)
	if err != nil {
		return "", nil
	}

	// 140 = audio
	return client.GetStreamURL(video, video.Formats.FindByItag(140))
}
