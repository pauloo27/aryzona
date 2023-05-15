package scheduler

import (
	"time"

	"github.com/pauloo27/aryzona/internal/core/routine"
	"github.com/pauloo27/logger"
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
						logger.Errorf("scheduler caught %v", err)
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
	routine.Go(func() { scheduleLoop(1 * time.Second) })
}
