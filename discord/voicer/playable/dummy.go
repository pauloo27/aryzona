package playable

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

func (DummyPlayable) Pause() error {
	return nil
}

func (DummyPlayable) Unpause() error {
	return nil
}

func (DummyPlayable) TogglePause() error {
	return nil
}

func (DummyPlayable) GetDirectURL() (string, error) {
	return "", nil
}

func (t DummyPlayable) GetName() string {
	return t.Name
}

func (t DummyPlayable) GetFullTitle() (string, string) {
	return t.Title, t.Artist
}
