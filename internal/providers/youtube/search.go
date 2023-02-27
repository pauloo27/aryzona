package youtube

import (
	"fmt"
	"net/url"
	"sync"
	"time"

	"github.com/Pauloo27/aryzona/internal/config"
	"github.com/Pauloo27/aryzona/internal/core/h"
	"github.com/Pauloo27/aryzona/internal/discord/voicer/playable"
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
	buf, err := h.Get(uri)
	if err != nil {
		return nil, err
	}

	results := gjson.GetBytes(buf, "items.#.id.videoId").Array()

	ids := make([]string, len(results))

	for i, id := range results {
		ids[i] = id.String()
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

	wg := sync.WaitGroup{}

	results := make([]*SearchResult, len(ids))
	for i, id := range ids {
		wg.Add(1)
		go func(i int, id string) {
			defer wg.Done()
			vid, err := defaultClient.GetVideo(id)
			if err != nil {
				logger.Warnf("error getting video %s from playlist %s: %s", id, searchQuery, err)
				return
			}
			results[i] = videoAsSearchResult(vid)
		}(i, id)
	}
	wg.Wait()
	return results, nil
}

func videoAsSearchResult(vid *yt.Video) *SearchResult {
	return &SearchResult{
		ID:           vid.ID,
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

	plVids := make([]*SearchResult, len(pl.Videos))
	for i, vid := range pl.Videos {
		plVids[i] = playlistVidAsSearchResult(vid)
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
		playables := make([]playable.Playable, len(r.plVids))
		for i, vid := range r.plVids {
			playables[i] = vid.ToPlayable()[0]
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
			vid:          r.vid,
		},
	}
}
