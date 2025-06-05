package spotify_test

import (
	"os"
	"testing"

	"github.com/pauloo27/aryzona/internal/providers/spotify"
	"github.com/stretchr/testify/assert"
)

var (
	clientID, clientSecret string

	playlistID = "1D6l3qeCbryB9COT1CGalw"
	trackID    = "6K4t31amVTZDgR3sKmwUJJ"

	sfy *spotify.Spotify
)

func TestMain(m *testing.M) {
	clientID = os.Getenv("SPOTIFY_CLIENT_ID")
	clientSecret = os.Getenv("SPOTIFY_CLIENT_SECRET")

	os.Exit(m.Run())
}

func TestNewSpotify(t *testing.T) {
	if clientID == "" {
		t.Skip("spotify: SPOTIFY_CLIENT_ID is not set")
	}
	if clientSecret == "" {
		t.Skip("spotify: SPOTIFY_CLIENT_SECRET is not set")
	}

	sfy = spotify.NewSpotify(clientID, clientSecret)

	assert.NotNil(t, sfy)
}

func TestGetPlaylist(t *testing.T) {
	if sfy == nil {
		t.Skip("spotify: spotify instance is nil")
	}

	playlist, err := sfy.GetPlaylist(playlistID)
	assert.NoError(t, err)

	assert.NotNil(t, playlist)
	assert.Equal(t, "forbidden", playlist.Name)
	assert.Equal(t, "mauriciofsnts", playlist.Owner.DisplayName)

	assert.Len(t, playlist.Images, 3)
	assert.NotEmpty(t, playlist.Images[0].URL)
	assert.NotEmpty(t, playlist.Images[1].URL)
	assert.NotEmpty(t, playlist.Images[2].URL)

	assert.NotNil(t, playlist.Tracks)
	assert.Greater(t, playlist.Tracks.Total, 25)
}

func TestGetPlaylistItems(t *testing.T) {
	if sfy == nil {
		t.Skip("spotify: spotify instance is nil")
	}

	playlist, err := sfy.GetPlaylistItems(playlistID, 10, 0)
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

	track, err := sfy.GetTrack(trackID)
	assert.NoError(t, err)

	assert.NotNil(t, track)
	assert.Equal(t, "The Less I Know The Better", track.Name)
	assert.Len(t, track.Artists, 1)
	assert.Equal(t, "Tame Impala", track.Artists[0].Name)
}
