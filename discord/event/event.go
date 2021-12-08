package event

import "errors"

type EventType int

const (
	Ready = EventType(iota)
	MessageCreated
)

var (
	ErrEventNotSupported = errors.New("event not supported")
)
