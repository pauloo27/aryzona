package scheduler

import (
	"log/slog"
	"time"

	"github.com/pauloo27/aryzona/internal/core/routine"
)

type TaskCallback func(params ...any)

type Task struct {
	LastRunAt             time.Time
	RunAt                 time.Time
	Delay                 time.Duration
	Repeat, RepeatCounter int
	Callback              TaskCallback
}

var tasks = make(map[string]*Task)

func Schedule(key string, task *Task) {
	tasks[key] = task
}

func IsScheduled(key string) bool {
	_, ok := tasks[key]
	return ok
}

func Unschedule(key string) {
	delete(tasks, key)
}

func scheduleLoop(delay time.Duration) {
	for {
		for key, task := range tasks {
			now := time.Now()
			if task.RunAt.After(now) {
				continue
			}

			go func() {
				defer func() {
					if err := recover(); err != nil {
						slog.Error("scheduler caught error", "err", err)
					}
				}()
				task.Callback()
			}()

			task.RepeatCounter++
			if task.Repeat > 0 && task.Repeat-task.RepeatCounter == 0 {
				Unschedule(key)
			}
			task.RunAt = now.Add(task.Delay)
		}
		time.Sleep(delay)
	}
}

func init() {
	routine.GoAndRecover(func() { scheduleLoop(1 * time.Second) })
}
