package queue

type EventAppendData struct {
	Items  []*QueueEntry
	Queue  *Queue
	Index  int
	IsMany bool
}

type EventRemoveData struct {
	Queue *Queue
	Index int
	Item  *QueueEntry
}
