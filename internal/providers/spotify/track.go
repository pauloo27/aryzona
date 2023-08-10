package spotify

import (
	"fmt"
)

type Track struct {
	Name    string `json:"name"`
	Artists []struct {
		Name string `json:"name"`
	}
}

func (s *Spotify) GetTrack(id string) (*Track, error) {
	uri := fmt.Sprintf(
		"https://api.spotify.com/v1/tracks/%s",
		id,
	)

	res, err := s.Get(uri)
	if err != nil {
		return nil, err
	}

	if res.StatusCode != 200 {
		return nil, fmt.Errorf("spotify: cannot get track items, got status code %d", res.StatusCode)
	}

	return parseBody[Track](res)
}
