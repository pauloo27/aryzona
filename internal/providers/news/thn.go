package news

import (
	"github.com/mmcdole/gofeed"
)

func GetTHNFeed() (*NewsFeed, error) {
	fp := gofeed.NewParser()
	feed, err := fp.ParseURL("https://feeds.feedburner.com/TheHackersNews?format=xml")
	if err != nil {
		return nil, err
	}

	entries := make([]*NewsEntry, len(feed.Items))

	thumbnailURL := ""

	if feed.Image != nil {
		thumbnailURL = feed.Image.URL
	}

	for i, item := range feed.Items {
		entry := &NewsEntry{
			Title:        item.Title,
			Author:       item.Author.Name,
			URL:          item.Link,
			Description:  item.Description,
			ThumbnailURL: thumbnailURL,
			PostedAt:     item.PublishedParsed,
			EditedAt:     item.UpdatedParsed,
		}
		entries[i] = entry
	}

	return &NewsFeed{
		Title:        feed.Title,
		Description:  feed.Description,
		URL:          feed.Link,
		Author:       feed.Author.Name,
		ThumbnailURL: feed.Image.URL,
		Entries:      entries,
	}, nil
}
