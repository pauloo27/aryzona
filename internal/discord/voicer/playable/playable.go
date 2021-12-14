package playable

import "time"

type Playable interface {
	CanPause() bool
	IsLive() bool
	GetDuration() (time.Duration, error)
	GetDirectURL() (string, error)
	GetName() string
	IsOppus() bool
	IsLocal() bool
	GetFullTitle() (title, artist string)
	GetThumbnailURL() (string, error)
}
