package discord

type PresenceType int

const (
	PresencePlaying = PresenceType(iota)
	PresenceListening
	PresenceStreaming
)

type Presence struct {
	Type         PresenceType
	Title, Extra string
}
