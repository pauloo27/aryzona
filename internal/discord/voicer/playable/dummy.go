package playable

import "time"

/*
	DummyPlayable is used in tests
*/
type DummyPlayable struct {
	Name, Artist, Title string
	Duration            time.Duration
}

var _ Playable = DummyPlayable{}

func (DummyPlayable) CanPause() bool {
	return false
}

func (DummyPlayable) IsOppus() bool {
	return true
}

func (DummyPlayable) IsLocal() bool {
	return true
}

func (DummyPlayable) IsLive() bool {
	return false
}

func (d DummyPlayable) GetDuration() (time.Duration, error) {
	return d.Duration, nil
}

func (DummyPlayable) GetDirectURL() (string, error) {
	return "", nil
}

func (DummyPlayable) GetThumbnailURL() (string, error) {
	return "", nil
}

func (d DummyPlayable) GetName() string {
	return d.Name
}

func (d DummyPlayable) GetFullTitle() (string, string) {
	return d.Title, d.Artist
}
