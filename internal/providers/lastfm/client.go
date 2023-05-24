package lastfm

import (
	"fmt"
	"io"
	"net/http"
	"net/url"

	"github.com/tidwall/gjson"
)

type LastFmClient struct {
	apiKey     string
	httpClient *http.Client
}

var (
	DefaultClient *LastFmClient
)

func NewLastFmClient(apiKey string) *LastFmClient {
	return &LastFmClient{
		apiKey:     apiKey,
		httpClient: &http.Client{},
	}
}

func (c *LastFmClient) GetSimilarTracks(artist, track string, limit int) ([]Track, error) {
	path :=
		fmt.Sprintf(
			"https://ws.audioscrobbler.com/2.0/?method=track.getsimilar&track=%s&artist=%s&autocorrect=1&api_key=%s&format=json",
			url.QueryEscape(track),
			url.QueryEscape(artist),
			c.apiKey,
		)

	if limit > 0 {
		path += fmt.Sprintf("&limit=%d", limit)
	}

	req, err := http.NewRequest("GET", path, nil)

	if err != nil {
		return nil, err
	}
	res, err := c.httpClient.Do(req)

	if err != nil {
		return nil, err
	}

	if res.StatusCode != 200 {
		return nil, fmt.Errorf("status code was %d but 200 was expected", res.StatusCode)
	}

	b, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	found := 0
	tracks := make([]Track, limit)

	data := gjson.ParseBytes(b)
	data.Get("similartracks.track").ForEach(func(key, value gjson.Result) bool {
		track := Track{
			Name:   value.Get("name").String(),
			Artist: value.Get("artist.name").String(),
			URL:    value.Get("url").String(),
		}
		tracks[found] = track
		found++

		return true
	})

	return tracks[:found], nil
}

func (c *LastFmClient) SearchTrack(query string, limit int) ([]Track, error) {
	path :=
		fmt.Sprintf(
			"https://ws.audioscrobbler.com/2.0/?method=track.search&track=%s&autocorrect=1&api_key=%s&format=json",
			url.QueryEscape(query),
			c.apiKey,
		)

	if limit > 0 {
		path += fmt.Sprintf("&limit=%d", limit)
	}

	req, err := http.NewRequest("GET", path, nil)

	if err != nil {
		return nil, err
	}
	res, err := c.httpClient.Do(req)

	if err != nil {
		return nil, err
	}

	if res.StatusCode != 200 {
		return nil, fmt.Errorf("status code was %d but 200 was expected", res.StatusCode)
	}

	b, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	found := 0
	tracks := make([]Track, limit)

	data := gjson.ParseBytes(b)
	data.Get("results.trackmatches.track").ForEach(func(key, value gjson.Result) bool {
		track := Track{
			Name:   value.Get("name").String(),
			Artist: value.Get("artist").String(),
			URL:    value.Get("url").String(),
		}
		tracks[found] = track
		found++

		return true
	})

	return tracks[:found], nil
}

func (c *LastFmClient) GetTopTracks(limit int) ([]Track, error) {
	path :=
		fmt.Sprintf(
			"https://ws.audioscrobbler.com/2.0/?method=chart.gettoptracks&api_key=%s&format=json",
			c.apiKey,
		)

	if limit > 0 {
		path += fmt.Sprintf("&limit=%d", limit)
	}

	req, err := http.NewRequest("GET", path, nil)

	if err != nil {
		return nil, err
	}
	res, err := c.httpClient.Do(req)

	if err != nil {
		return nil, err
	}

	if res.StatusCode != 200 {
		return nil, fmt.Errorf("status code was %d but 200 was expected", res.StatusCode)
	}

	b, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	found := 0
	tracks := make([]Track, limit)

	data := gjson.ParseBytes(b)
	data.Get("tracks.track").ForEach(func(key, value gjson.Result) bool {
		track := Track{
			Name:   value.Get("name").String(),
			Artist: value.Get("artist.name").String(),
			URL:    value.Get("url").String(),
		}
		tracks[found] = track
		found++

		return true
	})

	return tracks[:found], nil
}
