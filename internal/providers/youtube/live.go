package youtube

import (
	"errors"
	"io"
	"net/http"
	"strings"

	"github.com/kkdai/youtube/v2"
)

/* #nosec GG107 */
func getFirstURL(manifestURL string) (string, error) {
	res, err := http.Get(manifestURL)
	if err != nil {
		return "", err
	}
	/* #nosec G307 */
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

	return "", errors.New("URL not found in manifest")
}

func getLiveURL(video *youtube.Video) (string, error) {
	manifest := video.HLSManifestURL
	if manifest == "" {
		return "", errors.New(
			"HLS manifest not found",
		)
	}
	liveURL, err := getFirstURL(manifest)
	return liveURL, err
}
