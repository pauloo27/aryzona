package scheduler

import (
	"time"

	"github.com/Pauloo27/aryzona/utils"
	"github.com/Pauloo27/logger"
)

type Task struct {
	Time     time.Time
	Callback func(params ...interface{})
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
			if task.Time.After(time.Now()) {
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
			Unschedule(key)
		}
		time.Sleep(delay)
	}
}

func init() {
	utils.Go(func() { scheduleLoop(1 * time.Second) })
}
