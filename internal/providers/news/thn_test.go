package news_test

import (
	"testing"

	"github.com/Pauloo27/aryzona/internal/providers/news"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGetTHNFeed(t *testing.T) {
	feed, err := news.GetTHNFeed()
	require.Nil(t, err)
	assert.Equal(t, "The Hacker News", feed.Title)
	assert.Equal(t, "Most trusted, widely-read independent cybersecurity news source for everyone; supported by hackers and IT professionals — Send TIPs to admin@thehackernews.com", feed.Description)
	assert.Equal(t, "Unknown", feed.Author)
	assert.Equal(t, "https://thehackernews.com", feed.URL)
	assert.Equal(t, "", feed.ThumbnailURL)
	assert.NotEmpty(t, feed.Entries)
}
