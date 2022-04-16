package news

import "time"

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
