package youtube

import (
	"github.com/Pauloo27/aryzona/audio"
	"github.com/kkdai/youtube/v2"
)

type YouTubePlayable struct {
	Video  *youtube.Video
	IsLive bool
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
	if p.IsLive {
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

func AsPlayable(videoURL string) (audio.Playable, error) {
	vid, err := defaultClient.GetVideo(GetVideoID(videoURL))
	if err != nil {
		return nil, err
	}
	return YouTubePlayable{
		Video: vid,
	}, nil
}

func AsPlayableLive(liveURL string) (audio.Playable, error) {
	vid, err := defaultClient.GetVideo(GetVideoID(liveURL))
	if err != nil {
		return nil, err
	}
	return YouTubePlayable{
		IsLive: true,
		Video:  vid,
	}, nil
}
