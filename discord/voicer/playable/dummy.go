package playable

import "time"

/*
	DummyPlayable is used in tests
*/
type DummyPlayable struct {
	Name, Artist, Title string
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

func (DummyPlayable) GetDuration() (time.Duration, error) {
	return 0, nil
}

func (DummyPlayable) GetDirectURL() (string, error) {
	return "", nil
}

func (DummyPlayable) GetThumbnailURL() (string, error) {
	return "", nil
}

func (t DummyPlayable) GetName() string {
	return t.Name
}

func (t DummyPlayable) GetFullTitle() (string, string) {
	return t.Title, t.Artist
}
