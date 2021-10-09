package youtube

import (
	"time"

	"github.com/kkdai/youtube/v2"
)

type YouTubePlayable struct {
	Video *youtube.Video
	Live  bool
}

func (YouTubePlayable) CanPause() bool {
	return false
}

func (YouTubePlayable) GetName() string {
	return "YouTube video"
}

func (p YouTubePlayable) IsLive() bool {
	return p.Live
}

func (p YouTubePlayable) GetDuration() (time.Duration, error) {
	return p.Video.Duration, nil
}

func (YouTubePlayable) TogglePause() error {
	return nil
}

func (p YouTubePlayable) GetDirectURL() (string, error) {
	if p.Live {
		return getLiveURL(p.Video)
	}
	return defaultClient.GetStreamURL(p.Video, p.Video.Formats.FindByItag(140))
}

func (p YouTubePlayable) GetFullTitle() (title string, artist string) {
	return p.Video.Title, p.Video.Author
}

func (YouTubePlayable) IsLocal() bool {
	return false
}

func (YouTubePlayable) IsOppus() bool {
	return false
}

func AsPlayable(videoURL string) (YouTubePlayable, error) {
	vid, err := defaultClient.GetVideo(GetVideoID(videoURL))
	if err != nil {
		return YouTubePlayable{}, err
	}
	return YouTubePlayable{
		Video: vid,
		Live:  vid.Duration == 0,
	}, nil
}
