package playable_test

import "github.com/Pauloo27/aryzona/audio"

type TestPlayable struct {
	Name, Artist, Title string
}

var _ audio.Playable = TestPlayable{}

func (TestPlayable) CanPause() bool {
	return false
}

func (TestPlayable) IsOppus() bool {
	return true
}

func (TestPlayable) IsLocal() bool {
	return true
}

func (TestPlayable) Pause() error {
	return nil
}

func (TestPlayable) Unpause() error {
	return nil
}

func (TestPlayable) TogglePause() error {
	return nil
}

func (TestPlayable) GetDirectURL() (string, error) {
	return "", nil
}

func (t TestPlayable) GetName() string {
	return t.Name
}

func (t TestPlayable) GetFullTitle() (string, string) {
	return t.Title, t.Artist
}
