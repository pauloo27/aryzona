package youtube

import (
	"github.com/Pauloo27/aryzona/audio"
	"github.com/kkdai/youtube/v2"
)

type YouTubePlayable struct {
	Video *youtube.Video
}

func (YouTubePlayable) CanPause() bool {
	return false
}

func (YouTubePlayable) Pause() error {
	return nil
}

func (YouTubePlayable) Unpause() error {
	return nil
}

func (YouTubePlayable) TogglePause() error {
	return nil
}

func (p YouTubePlayable) GetDirectURL() (string, error) {
	return defaultClient.GetStreamURL(p.Video, p.Video.Formats.FindByItag(140))
}

func (p YouTubePlayable) GetFullTitle() (artist string, title string) {
	return p.Video.Author, p.Video.Title
}

func (YouTubePlayable) IsLocal() bool {
	return false
}

func (YouTubePlayable) IsOppus() bool {
	return false
}

func AsPlayable(videoURL string) (audio.Playable, error) {
	vid, err := defaultClient.GetVideo(GetVideoID(videoURL))
	if err != nil {
		return nil, err
	}
	return YouTubePlayable{
		Video: vid,
	}, nil
}
