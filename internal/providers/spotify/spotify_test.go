package spotify_test

import (
	"os"
	"testing"

	"github.com/pauloo27/aryzona/internal/providers/spotify"
	"github.com/stretchr/testify/assert"
)

var (
	clientId, clientSecret string

	playlistId = "1D6l3qeCbryB9COT1CGalw"
	trackId    = "6K4t31amVTZDgR3sKmwUJJ"

	sfy *spotify.Spotify
)

func TestMain(m *testing.M) {
	clientId = os.Getenv("SPOTIFY_CLIENT_ID")
	clientSecret = os.Getenv("SPOTIFY_CLIENT_SECRET")

	os.Exit(m.Run())
}

func TestNewSpotify(t *testing.T) {
	if clientId == "" {
		t.Skip("spotify: SPOTIFY_CLIENT_ID is not set")
	}
	if clientSecret == "" {
		t.Skip("spotify: SPOTIFY_CLIENT_SECRET is not set")
	}

	sfy = spotify.NewSpotify(clientId, clientSecret)

	assert.NotNil(t, sfy)
}

func TestGetPlaylist(t *testing.T) {
	if sfy == nil {
		t.Skip("spotify: spotify instance is nil")
	}

	playlist, err := sfy.GetPlaylist(playlistId)
	assert.NoError(t, err)

	assert.NotNil(t, playlist)
	assert.Equal(t, "forbidden", playlist.Name)
	assert.Equal(t, "mauriciofsnts", playlist.Owner.DisplayName)

	assert.Len(t, playlist.Images, 3)
	assert.NotEmpty(t, playlist.Images[0].Url)
	assert.NotEmpty(t, playlist.Images[1].Url)
	assert.NotEmpty(t, playlist.Images[2].Url)

	assert.NotNil(t, playlist.Tracks)
	assert.Greater(t, playlist.Tracks.Total, 25)
}

func TestGetPlaylistItems(t *testing.T) {
	if sfy == nil {
		t.Skip("spotify: spotify instance is nil")
	}

	playlist, err := sfy.GetPlaylistItems(playlistId, 10, 0)
	assert.NoError(t, err)

	assert.NotNil(t, playlist)
	assert.Len(t, playlist.Items, 10)

	assert.NotNil(t, playlist.Items)
	assert.Greater(t, playlist.Total, 25)
}

func TestGetTrack(t *testing.T) {
	if sfy == nil {
		t.Skip("spotify: spotify instance is nil")
	}

	track, err := sfy.GetTrack(trackId)
	assert.NoError(t, err)

	assert.NotNil(t, track)
	assert.Equal(t, "The Less I Know The Better", track.Name)
	assert.Len(t, track.Artists, 1)
	assert.Equal(t, "Tame Impala", track.Artists[0].Name)
}
