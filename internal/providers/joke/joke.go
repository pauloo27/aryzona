package joke

import (
	"encoding/json"

	"github.com/Pauloo27/aryzona/internal/core/h"
)

type Joke struct {
	ID        int    `json:"id"`
	Type      string `json:"type"`
	Setup     string `json:"setup"`
	Punchline string `json:"punchline"`
}

func GetRandomJoke() (*Joke, error) {
	data, err := h.Get("https://official-joke-api.appspot.com/random_joke")
	if err != nil {
		return nil, err
	}

	var joke Joke
	err = json.Unmarshal(data, &joke)

	return &joke, err
}
