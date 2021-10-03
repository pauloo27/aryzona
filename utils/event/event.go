package event

type EventType string
type EventListener func(params ...interface{})

type EventEmitter struct {
	listeners map[EventType][]EventListener
}

func NewEventEmitter() *EventEmitter {
	emitter := EventEmitter{
		listeners: make(map[EventType][]EventListener),
	}
	return &emitter
}

func (q *EventEmitter) On(event EventType, listener EventListener) {
	q.listeners[event] = append(q.listeners[event], listener)
}

func (q *EventEmitter) Emit(event EventType, params ...interface{}) {
	listeners, found := q.listeners[event]
	if !found {
		return
	}
	for _, listener := range listeners {
		listener(params...)
	}
}
