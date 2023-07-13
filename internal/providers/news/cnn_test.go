package news_test

import (
	"testing"

	"github.com/pauloo27/aryzona/internal/providers/news"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGetCNNTopStoriesFeed(t *testing.T) {
	feed, err := news.GetCNNTopStoriesFeed()
	require.Nil(t, err)
	assert.Equal(t, "CNN.com - RSS Channel - HP News", feed.Title)
	assert.Equal(t, "CNN.com delivers up-to-the-minute news and information on the latest top stories, weather, entertainment, politics and more.", feed.Description)
	assert.Equal(t, "Unknown", feed.Author) // =(
	assert.Equal(t, "https://www.cnn.com/homepage2/index.html", feed.URL)
	assert.Equal(t, "http://i2.cdn.turner.com/cnn/2015/images/09/24/cnn.digital.png", feed.ThumbnailURL)
	assert.NotEmpty(t, feed.Entries)
}

func TestGetCNNWorldFeed(t *testing.T) {
	feed, err := news.GetCNNWorldFeed()
	require.Nil(t, err)
	assert.Equal(t, "CNN.com - RSS Channel - World", feed.Title)
	assert.Equal(t, "CNN.com delivers up-to-the-minute news and information on the latest top stories, weather, entertainment, politics and more.", feed.Description)
	assert.Equal(t, "Unknown", feed.Author) // =(
	assert.Equal(t, "https://www.cnn.com/world/index.html", feed.URL)
	assert.Equal(t, "http://i2.cdn.turner.com/cnn/2015/images/09/24/cnn.digital.png", feed.ThumbnailURL)
	assert.NotEmpty(t, feed.Entries)
}

func TestGetCNNTechFeed(t *testing.T) {
	feed, err := news.GetCNNTechFeed()
	require.Nil(t, err)
	assert.Equal(t, "CNN.com - RSS Channel - App Tech Section", feed.Title)
	assert.Equal(t, "CNN.com delivers up-to-the-minute news and information on the latest top stories, weather, entertainment, politics and more.", feed.Description)
	assert.Equal(t, "Unknown", feed.Author) // =(
	assert.Equal(t, "https://www.cnn.com/app-tech-section/index.html", feed.URL)
	assert.Equal(t, "http://i2.cdn.turner.com/cnn/2015/images/09/24/cnn.digital.png", feed.ThumbnailURL)
	assert.NotEmpty(t, feed.Entries)
}
