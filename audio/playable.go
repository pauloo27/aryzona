package audio

type Playable interface {
	CanPause() bool
	Pause() error
	Unpause() error
	TogglePause() error
	GetDirectURL() (string, error)
	GetName() string
	IsOppus() bool
	IsLocal() bool
	GetFullTitle() (title, artist string)
}
