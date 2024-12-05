package youtube

import (
	"bytes"
	"errors"
	"fmt"
	"os/exec"
	"strings"
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
	// Check if the URL is for a live stream
	if p.Live {
		return p.getLiveStreamURL()
	}

	return p.getVideoURL()
}

func (p YouTubePlayable) getLiveStreamURL() (string, error) {
	/* #nosec G204 get rekt */
	cmd := exec.Command("yt-dlp", "-g", "-f", "best", p.GetShareURL())
	var out bytes.Buffer
	cmd.Stdout = &out
	err := cmd.Run()
	if err != nil {
		return "", errors.New("failed to fetch live stream URL: " + err.Error())
	}

	return strings.TrimSpace(out.String()), nil
}

func (p YouTubePlayable) getVideoURL() (string, error) {
	/* #nosec G204 get rekt */
	cmd := exec.Command("yt-dlp", "-g", "-f", "251", p.GetShareURL()) // itag 251 = audio (webm)
	var out bytes.Buffer
	cmd.Stdout = &out
	err := cmd.Run()
	if err != nil {
		return "", errors.New("failed to fetch video URL: " + err.Error())
	}

	// Trim the output to remove any extra spaces or newline characters
	return strings.TrimSpace(out.String()), nil
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
