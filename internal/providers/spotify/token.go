package spotify

import (
	"fmt"
	"net/url"
	"time"
)

type Token struct {
	AccessToken string `json:"access_token"`
	TokenType   string `json:"token_type"`
	ExpiresIn   int64  `json:"expires_in"`
	ExpiresAt   time.Time
}

func (s *Spotify) generateToken() error {
	reqData := url.Values{}
	reqData.Set("grant_type", "client_credentials")
	reqData.Set("client_id", s.clientID)
	reqData.Set("client_secret", s.clientSecret)

	requestedAt := time.Now()

	res, err := s.PostForm("https://accounts.spotify.com/api/token", reqData)
	if err != nil {
		return err
	}

	if res.StatusCode != 200 {
		return fmt.Errorf("spotify: cannot generate token, got status code %d", res.StatusCode)
	}

	token, err := parseBody[Token](res)
	if err != nil {
		return err
	}

	token.ExpiresAt = requestedAt.Add(time.Duration(token.ExpiresIn) * time.Second)

	s.token = token

	return nil
}
