package news_test

import (
	"testing"

	"github.com/Pauloo27/aryzona/internal/providers/news"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGetCNNTopStoriesFeed(t *testing.T) {
	feed, err := news.GetCNNTopStoriesFeed()
	require.Nil(t, err)
	assert.Equal(t, feed.Title, "CNN.com - RSS Channel - HP Hero")
	assert.Equal(t, feed.Description, "CNN.com delivers up-to-the-minute news and information on the latest top stories, weather, entertainment, politics and more.")
	assert.Equal(t, feed.Author, "Unknown") // =(
	assert.Equal(t, feed.URL, "https://www.cnn.com/index.html")
	assert.Equal(t, feed.ThumbnailURL, "http://i2.cdn.turner.com/cnn/2015/images/09/24/cnn.digital.png")
	assert.NotEmpty(t, feed.Entries)
}

func TestGetCNNWorldFeed(t *testing.T) {
	feed, err := news.GetCNNWorldFeed()
	require.Nil(t, err)
	assert.Equal(t, feed.Title, "CNN.com - RSS Channel - World")
	assert.Equal(t, feed.Description, "CNN.com delivers up-to-the-minute news and information on the latest top stories, weather, entertainment, politics and more.")
	assert.Equal(t, feed.Author, "Unknown") // =(
	assert.Equal(t, feed.URL, "https://www.cnn.com/world/index.html")
	assert.Equal(t, feed.ThumbnailURL, "http://i2.cdn.turner.com/cnn/2015/images/09/24/cnn.digital.png")
	assert.NotEmpty(t, feed.Entries)
}

func TestGetCNNTechFeed(t *testing.T) {
	feed, err := news.GetCNNTechFeed()
	require.Nil(t, err)
	assert.Equal(t, feed.Title, "CNN.com - RSS Channel - App Tech Section")
	assert.Equal(t, feed.Description, "CNN.com delivers up-to-the-minute news and information on the latest top stories, weather, entertainment, politics and more.")
	assert.Equal(t, feed.Author, "Unknown") // =(
	assert.Equal(t, feed.URL, "https://www.cnn.com/app-tech-section/index.html")
	assert.Equal(t, feed.ThumbnailURL, "http://i2.cdn.turner.com/cnn/2015/images/09/24/cnn.digital.png")
	assert.NotEmpty(t, feed.Entries)
}
