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
	assert.Equal(t, feed.Title, "The Hacker News")
	assert.Equal(t, feed.Description, "Most trusted, widely-read independent cybersecurity news source for everyone; supported by hackers and IT professionals — Send TIPs to admin@thehackernews.com")
	assert.Equal(t, feed.Author, "Swati Khandelwal")
	assert.Equal(t, feed.URL, "https://thehackernews.com/")
	assert.Equal(t, feed.ThumbnailURL, "http://creativecommons.org/images/public/somerights20.gif")
	assert.NotEmpty(t, feed.Entries)
}
