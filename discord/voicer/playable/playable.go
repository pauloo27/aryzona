package playable

type Playable interface {
	CanPause() bool
	GetDirectURL() (string, error)
	GetName() string
	IsOppus() bool
	IsLocal() bool
	GetFullTitle() (title, artist string)
}
