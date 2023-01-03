package youtube

import (
	"fmt"
	"net/url"
	"time"

	"github.com/Pauloo27/aryzona/internal/config"
	"github.com/Pauloo27/aryzona/internal/discord/voicer/playable"
	"github.com/Pauloo27/aryzona/internal/utils"
	"github.com/Pauloo27/logger"
	yt "github.com/kkdai/youtube/v2"
	"github.com/tidwall/gjson"
)

type SearchResult struct {
	vid           *yt.Video
	pl            *yt.Playlist
	plVids        []*SearchResult // will be lazy loaded
	ID            string
	Title, Author string
	Duration      time.Duration
	ThumbnailURL  string
}

func searchWithAPI(searchQuery string, limit int) ([]string, error) {
	apiKey := config.Config.YoutubeAPIKey
	uri := fmt.Sprintf(
		"https://www.googleapis.com/youtube/v3/search?q=%s&key=%s&maxResults=%d&type=video",
		url.QueryEscape(searchQuery), url.QueryEscape(apiKey),
		limit,
	)
	buf, err := utils.Get(uri)
	if err != nil {
		return nil, err
	}
	var ids []string
	for _, id := range gjson.GetBytes(buf, "items.#.id.videoId").Array() {
		ids = append(ids, id.String())
	}
	return ids, nil
}

func SearchFor(searchQuery string, limit int) ([]*SearchResult, error) {
	if pl, err := defaultClient.GetPlaylist(searchQuery); err == nil && len(pl.Videos) != 0 {
		return []*SearchResult{playlistAsSearchResult(pl)}, nil
	}

	if vid, err := defaultClient.GetVideo(searchQuery); err == nil {
		return []*SearchResult{videoAsSearchResult(vid)}, nil
	}

	ids, err := searchWithAPI(searchQuery, limit)
	if err != nil {
		return nil, err
	}

	var results []*SearchResult
	for _, id := range ids {
		vid, err := defaultClient.GetVideo(id)
		if err != nil {
			logger.Warnf("error getting video %s from playlist %s: %s", id, searchQuery, err)
			continue
		}
		results = append(results, videoAsSearchResult(vid))
	}
	return results, nil
}

func videoAsSearchResult(vid *yt.Video) *SearchResult {
	return &SearchResult{
		vid:          vid,
		Title:        vid.Title,
		Author:       vid.Author,
		Duration:     vid.Duration,
		ThumbnailURL: vid.Thumbnails[0].URL,
	}
}

func playlistAsSearchResult(pl *yt.Playlist) *SearchResult {
	duration := time.Duration(0)
	for _, vid := range pl.Videos {
		duration += vid.Duration
	}

	var plVids []*SearchResult
	for _, vid := range pl.Videos {
		plVids = append(plVids, playlistVidAsSearchResult(vid))
	}

	return &SearchResult{
		pl:           pl,
		plVids:       plVids,
		ID:           pl.ID,
		Title:        pl.Title,
		Author:       pl.Author,
		Duration:     duration,
		ThumbnailURL: pl.Videos[0].Thumbnails[0].URL,
	}
}

func playlistVidAsSearchResult(vid *yt.PlaylistEntry) *SearchResult {
	return &SearchResult{
		ID:           vid.ID,
		Title:        vid.Title,
		Author:       vid.Author,
		Duration:     vid.Duration,
		ThumbnailURL: fmt.Sprintf("https://img.youtube.com/vi/%s/mqdefault.jpg", vid.ID),
	}
}

func (r *SearchResult) IsPlaylist() bool {
	return r.pl != nil
}

func (r *SearchResult) IsLive() bool {
	return r.vid != nil && r.Duration == 0
}

func (r *SearchResult) ToPlayable() []playable.Playable {
	if r.IsPlaylist() {
		var playables []playable.Playable
		for _, vid := range r.plVids {
			playables = append(playables, vid.ToPlayable()[0])
		}
		return playables
	}
	return []playable.Playable{
		YouTubePlayable{
			ID:           r.ID,
			Title:        r.Title,
			Author:       r.Author,
			ThumbnailURL: r.ThumbnailURL,
			Duration:     r.Duration,
			Live:         r.IsLive(),
			video:        r.vid,
		},
	}
}
