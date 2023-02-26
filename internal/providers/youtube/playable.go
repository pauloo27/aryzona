package youtube

import (
	"fmt"
	"time"

	"github.com/kkdai/youtube/v2"
)

type YouTubePlayable struct {
	ID, Title, Author, ThumbnailURL string
	Duration                        time.Duration
	vid                             *youtube.Video
	Live                            bool
}

func (p YouTubePlayable) video() *youtube.Video {
	if p.vid != nil {
		return p.vid
	}
	p.vid, _ = defaultClient.GetVideo(p.ID)
	return p.vid
}

func (p YouTubePlayable) CanPause() bool {
	return !p.Live
}

func (YouTubePlayable) GetName() string {
	return "YouTube video"
}

func (p YouTubePlayable) GetShareURL() string {
	return fmt.Sprintf("https://youtu.be/%s", p.video().ID)
}

func (p YouTubePlayable) IsLive() bool {
	return p.Live
}

func (p YouTubePlayable) GetThumbnailURL() (string, error) {
	return p.ThumbnailURL, nil
}

func (p YouTubePlayable) GetDuration() (time.Duration, error) {
	return p.Duration, nil
}

func (YouTubePlayable) TogglePause() error {
	return nil
}

func (p YouTubePlayable) GetDirectURL() (string, error) {
	if p.Live {
		return getLiveURL(p.video())
	}
	if format := p.video().Formats.FindByItag(251); format != nil {
		return defaultClient.GetStreamURL(p.video(), format)
	}
	return defaultClient.GetStreamURL(p.video(), p.video().Formats.FindByItag(140))
}

func (p YouTubePlayable) GetFullTitle() (title string, artist string) {
	return p.Title, p.Author
}

func (YouTubePlayable) IsLocal() bool {
	return false
}

func (p YouTubePlayable) IsOpus() bool {
	return p.video().Formats.FindByItag(251) != nil
}

func GetVideo(videoURL string) (YouTubePlayable, error) {
	vid, err := defaultClient.GetVideo(videoURL)
	if err != nil {
		return YouTubePlayable{}, err
	}
	return YouTubePlayable{
		ID:           vid.ID,
		Title:        vid.Title,
		Author:       vid.Author,
		ThumbnailURL: vid.Thumbnails[0].URL,
		Duration:     vid.Duration,
		Live:         vid.Duration == 0,
		vid:          vid,
	}, nil
}
