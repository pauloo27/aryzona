package queue

import "github.com/Pauloo27/aryzona/internal/discord/voicer/playable"

type EventAppendData struct {
	Items  []playable.Playable
	Queue  *Queue
	Index  int
	IsMany bool
}

type EventRemoveData struct {
	Queue *Queue
	Index int
}
