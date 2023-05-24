package lastfm_test

import (
	"os"
	"testing"

	"github.com/pauloo27/aryzona/internal/providers/lastfm"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var (
	client *lastfm.LastFmClient
	apiKey string
)

func TestMain(m *testing.M) {
	// TODO: add to pipeline
	apiKey = os.Getenv("LAST_FM_API_KEY")
	client = lastfm.NewLastFmClient(apiKey)
	os.Exit(m.Run())
}

func TestGetSimilarTracks(t *testing.T) {
	if apiKey == "" {
		t.Skip()
		return
	}

	tracks, err := client.GetSimilarTracks("drake", "god's plan", 2)
	assert.Nil(t, err)
	assert.NotEmpty(t, tracks)
	require.Len(t, tracks, 2)

	assert.NotEmpty(t, tracks[0].Name)
	assert.NotEmpty(t, tracks[0].Artist)
	assert.NotEmpty(t, tracks[0].URL)

	assert.NotEmpty(t, tracks[1].Name)
	assert.NotEmpty(t, tracks[1].Artist)
	assert.NotEmpty(t, tracks[1].URL)
}

func TestSearch(t *testing.T) {
	if apiKey == "" {
		t.Skip()
		return
	}

	tracks, err := client.SearchTrack("drake god plans", 1)
	assert.Nil(t, err)
	assert.NotEmpty(t, tracks)
	require.Len(t, tracks, 1)

	assert.NotEmpty(t, tracks[0].Name)
	assert.NotEmpty(t, tracks[0].Artist)
	assert.NotEmpty(t, tracks[0].URL)
}

func TestGetTopTracks(t *testing.T) {
	if apiKey == "" {
		t.Skip()
		return
	}

	tracks, err := client.GetTopTracks(5)
	assert.Nil(t, err)
	assert.NotEmpty(t, tracks)
	require.Len(t, tracks, 5)

	for _, track := range tracks {
		assert.NotEmpty(t, track.Name)
		assert.NotEmpty(t, track.Artist)
		assert.NotEmpty(t, track.URL)
	}
}
