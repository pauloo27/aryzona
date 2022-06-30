package news

import (
	"time"

	"github.com/mmcdole/gofeed"
)

var (
	unixEpoch = time.Unix(0, 0)
)

type NewsEntry struct {
	Title        string
	ThumbnailURL string
	Description  string
	URL          string
	Author       string
	PostedAt     *time.Time
	EditedAt     *time.Time
}

type NewsFeed struct {
	Title        string
	Description  string
	ThumbnailURL string
	URL          string
	Author       string
	Entries      []*NewsEntry
}

func ParseFeed(feedURL string) (*NewsFeed, error) {
	fp := gofeed.NewParser()
	feed, err := fp.ParseURL(feedURL)
	if err != nil {
		return nil, err
	}

	entries := make([]*NewsEntry, len(feed.Items))

	feedAuthor, thumbnailURL := "Unknown", ""

	if len(feed.Authors) > 0 {
		feedAuthor = feed.Authors[0].Name
	}

	if feed.Image != nil {
		thumbnailURL = feed.Image.URL
	}

	for i, item := range feed.Items {
		entry := &NewsEntry{
			Title:        item.Title,
			URL:          item.Link,
			Description:  item.Description,
			ThumbnailURL: thumbnailURL,
			PostedAt:     item.PublishedParsed,
			EditedAt:     item.UpdatedParsed,
		}

		if len(item.Authors) > 0 {
			entry.Author = item.Authors[0].Name
		} else {
			entry.Author = "Unknown"
		}
		if entry.PostedAt == nil {
			entry.PostedAt = &unixEpoch
		}

		entries[i] = entry
	}

	return &NewsFeed{
		Title:        feed.Title,
		Description:  feed.Description,
		URL:          feed.Link,
		Author:       feedAuthor,
		ThumbnailURL: feed.Image.URL,
		Entries:      entries,
	}, nil
}
