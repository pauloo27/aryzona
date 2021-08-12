package youtube

import (
	"io"
	"net/http"
	"strings"

	"github.com/Pauloo27/aryzona/utils"
	"github.com/kkdai/youtube/v2"
)

/* #nosec GG107 */
func getFirstURL(manifestURL string) (string, error) {
	res, err := http.Get(manifestURL)
	if err != nil {
		return "", err
	}
	defer res.Body.Close()
	buffer, err := io.ReadAll(res.Body)
	if err != nil {
		return "", err
	}

	body := string(buffer)

	for _, line := range strings.Split(body, "\n") {
		if strings.HasPrefix(line, "https://") {
			return line, nil
		}
	}

	return "", utils.Errore{ID: "URL_NOT_FOUND", Message: "URL not found"}
}

func GetLiveURL(url string) (string, error) {
	video, err := defaultClient.GetVideo(GetVideoID(url))
	if err != nil {
		return "", err
	}
	return getLiveURL(video)
}

func getLiveURL(video *youtube.Video) (string, error) {
	manifest := video.HLSManifestURL
	if manifest == "" {
		return "", utils.Errore{
			ID:      "HLS_NOT_FOUND",
			Message: "HLS manifest not found",
		}
	}
	liveURL, err := getFirstURL(manifest)
	return liveURL, err
}
