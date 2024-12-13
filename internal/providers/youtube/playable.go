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

func (p YouTubePlayable) video() (*youtube.Video, error) {
	if p.vid != nil {
		return p.vid, nil
	}
	var err error
	p.vid, err = defaultClient.GetVideo(p.ID)
	return p.vid, err
}

func (p YouTubePlayable) CanPause() bool {
	return !p.Live
}

func (YouTubePlayable) GetName() string {
	return "YouTube video"
}

func (p YouTubePlayable) GetShareURL() string {
	vid, err := p.video()
	if err != nil {
		return ""
	}
	return fmt.Sprintf("https://youtu.be/%s", vid.ID)
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
	// for some reason, using the same client can result in some 403
	defaultClient = youtube.Client{}

	vid, err := p.video()
	if err != nil {
		return "", err
	}

	if p.Live {
		return getLiveURL(vid)
	}

	formats := vid.Formats.Itag(251)

	if len(formats) > 0 {
		format := formats[0]
		return defaultClient.GetStreamURL(vid, &format)
	}

	return defaultClient.GetStreamURL(vid, &vid.Formats[0])
}

func (p YouTubePlayable) GetFullTitle() (title string, artist string) {
	return p.Title, p.Author
}

func (YouTubePlayable) IsLocal() bool {
	return false
}

func (p YouTubePlayable) IsOpus() bool {
	vid, err := p.video()
	if err != nil {
		return false
	}
	return len(vid.Formats.Itag(251)) > 0
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
