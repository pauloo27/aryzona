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
	assert.Equal(t, "The Hacker News - Most Trusted Cyber Security and Computer Security Analysis", feed.Title)
	assert.Equal(t, "The Hacker News is the most trusted, widely-read, independent infosec source of the latest hacking news, cyber attacks, computer security, network security, and cybersecurity for ethical hackers, penetration testers, and information technology professionals.", feed.Description)
	assert.Equal(t, "Swati Khandelwal", feed.Author)
	assert.Equal(t, "https://thehackernews.com/", feed.URL)
	assert.Equal(t, "", feed.ThumbnailURL)
	assert.NotEmpty(t, feed.Entries)
}
