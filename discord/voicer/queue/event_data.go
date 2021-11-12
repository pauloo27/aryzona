package queue

import "github.com/Pauloo27/aryzona/discord/voicer/playable"

type EventAppendData struct {
	Queue  *Queue
	Index  int
	IsMany bool
	Items  []playable.Playable
}

type EventRemoveData struct {
	Queue *Queue
	Index int
}
