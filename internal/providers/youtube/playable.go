package youtube

import (
	"time"

	"github.com/Pauloo27/aryzona/internal/discord/voicer/playable"
	"github.com/Pauloo27/aryzona/internal/utils"
	"github.com/kkdai/youtube/v2"
)

type YouTubePlayablePlaylist struct {
	Videos        []playable.Playable
	Duration      time.Duration
	Title, Author string
}

type YouTubePlayable struct {
	ID, Title, Author, ThumbnailURL string
	Live                            bool
	Duration                        time.Duration
	opus                            bool
	video                           *youtube.Video
}

func (p YouTubePlayable) CanPause() bool {
	return !p.Live
}

func (YouTubePlayable) GetName() string {
	return "YouTube video"
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
	if p.video == nil {
		var err error
		p.video, err = defaultClient.GetVideo(p.ID)
		if err != nil {
			return "", err
		}
	}
	if p.Live {
		return getLiveURL(p.video)
	}
	format := p.video.Formats.FindByItag(251)
	if format != nil {
		p.opus = true
		return defaultClient.GetStreamURL(p.video, format)
	}
	return defaultClient.GetStreamURL(p.video, p.video.Formats.FindByItag(140))
}

func (p YouTubePlayable) GetFullTitle() (title string, artist string) {
	return p.Title, p.Author
}

func (YouTubePlayable) IsLocal() bool {
	return false
}

func (p YouTubePlayable) IsOppus() bool {
	return p.opus
}

func GetPlaylist(playlistURL string) (YouTubePlayablePlaylist, error) {
	playlist, err := defaultClient.GetPlaylist(playlistURL)
	if err != nil {
		return YouTubePlayablePlaylist{}, err
	}

	videos := make([]playable.Playable, len(playlist.Videos))

	duration := time.Duration(0)

	for i, vid := range playlist.Videos {
		duration += vid.Duration
		videos[i] = YouTubePlayable{
			ID:           vid.ID,
			Title:        vid.Title,
			Author:       vid.Author,
			Duration:     vid.Duration,
			Live:         vid.Duration == 0,
			ThumbnailURL: utils.Fmt("https://img.youtube.com/vi/%s/mqdefault.jpg", vid.ID),
			video:        nil, // will be lazy loaded
		}
	}

	return YouTubePlayablePlaylist{
		Title:    playlist.Title,
		Author:   playlist.Author,
		Videos:   videos,
		Duration: duration,
	}, nil
}

func AsPlayable(videoURL string) (YouTubePlayable, error) {
	vid, err := defaultClient.GetVideo(GetVideoID(videoURL))
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
		video:        vid,
	}, nil
}
