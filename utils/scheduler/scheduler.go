package scheduler

import (
	"time"

	"github.com/Pauloo27/aryzona/utils"
	"github.com/Pauloo27/logger"
)

type TaskCallback func(params ...interface{})

type Task struct {
	Delay                 time.Duration
	Repeat, RepeatCounter int
	Callback              TaskCallback
	LastRunAt             time.Time
	RunAt                 time.Time
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
	utils.Go(func() { scheduleLoop(1 * time.Second) })
}
