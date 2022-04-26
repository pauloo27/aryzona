package model

type PresenceType int

const (
	PresencePlaying = PresenceType(iota)
	PresenceListening
	PresenceStreaming
)

type Presence struct {
	Title, Extra string
	Type         PresenceType
}
