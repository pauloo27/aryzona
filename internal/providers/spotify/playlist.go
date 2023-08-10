package spotify

import (
	"fmt"
)

type Playlist struct {
	Name   string `json:"name"`
	Images []struct {
		Url    string `json:"url"`
		Height int    `json:"height"`
		Width  int    `json:"width"`
	} `json:"images"`
	Owner struct {
		DisplayName string `json:"display_name"`
	}
	Tracks PlaylistItems `json:"tracks"`
}

type PlaylistItems struct {
	Items []struct {
		Track Track `json:"track"`
	} `json:"items"`
	Limit int `json:"limit"`
	Total int `json:"total"`
}

const (
	playlistItemsFields = "name,total,limit,items(track(name,artists(name)))"
)

func (s *Spotify) GetPlaylist(id string) (*Playlist, error) {
	uri := fmt.Sprintf(
		"https://api.spotify.com/v1/playlists/%s",
		id,
	)

	res, err := s.Get(uri)
	if err != nil {
		return nil, err
	}

	if res.StatusCode != 200 {
		return nil, fmt.Errorf("spotify: cannot get playlist, got status code %d", res.StatusCode)
	}

	return parseBody[Playlist](res)
}

func (s *Spotify) GetPlaylistItems(id string, limit int, offset int) (*PlaylistItems, error) {
	uri := fmt.Sprintf(
		"https://api.spotify.com/v1/playlists/%s/tracks?limit=%d&offset=%d&fields=%s",
		id,
		limit,
		offset,
		playlistItemsFields,
	)

	res, err := s.Get(uri)
	if err != nil {
		return nil, err
	}

	if res.StatusCode != 200 {
		return nil, fmt.Errorf("spotify: cannot get playlist items, got status code %d", res.StatusCode)
	}

	return parseBody[PlaylistItems](res)
}
