package scheduler

import "time"

func NewRunLaterTask(delay time.Duration, callback TaskCallback) *Task {
	return NewRepeatingTask(delay, 1, callback)
}

func NewRepeatingTask(delay time.Duration, repeat int, callback TaskCallback) *Task {
	return &Task{
		RunAt:    time.Now().Add(delay),
		Repeat:   repeat,
		Delay:    delay,
		Callback: callback,
	}
}
