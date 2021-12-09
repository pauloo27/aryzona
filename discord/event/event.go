package event

import "errors"

type EventType int

const (
	Ready = EventType(iota)
	MessageCreated
	VoiceStateUpdated
)

var (
	ErrEventNotSupported = errors.New("event not supported")
)
