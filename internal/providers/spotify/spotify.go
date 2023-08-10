package spotify

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

type Spotify struct {
	clientID, clientSecret string

	token *Token

	*http.Client
}

type AuthTransport struct {
	Proxied http.RoundTripper
	Spotify *Spotify
}

func NewSpotify(clientID, clientSecret string) *Spotify {
	sfy := &Spotify{
		clientID:     clientID,
		clientSecret: clientSecret,
	}

	sfy.Client = &http.Client{
		Transport: AuthTransport{
			http.DefaultTransport,
			sfy,
		},
	}

	return sfy
}

func (t AuthTransport) RoundTrip(req *http.Request) (res *http.Response, e error) {
	if req.URL.Host != "api.spotify.com" {
		return t.Proxied.RoundTrip(req)
	}

	if t.Spotify.token == nil || time.Now().After(t.Spotify.token.ExpiresAt) {
		err := t.Spotify.generateToken()
		if err != nil {
			return nil, fmt.Errorf("spotify: cannot generate token: %w", err)
		}
	}

	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", t.Spotify.token.AccessToken))
	res, e = t.Proxied.RoundTrip(req)
	return res, e
}

func parseBody[T any](res *http.Response) (*T, error) {
	t := new(T)

	resData, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	err = json.Unmarshal(resData, t)

	return t, err
}
