package scheduler_test

import (
	"testing"
	"time"

	"github.com/Pauloo27/aryzona/utils/scheduler"
	"github.com/stretchr/testify/assert"
)

func TestRunLater(t *testing.T) {
	delay := time.Second

	now := time.Now()
	c := make(chan bool, 1)
	task := scheduler.NewRunLaterTask(delay, func(params ...interface{}) {
		c <- time.Since(now).Truncate(time.Second) == time.Second
	})
	scheduler.Schedule("test run later", task)
	assert.True(t, <-c)
}

func TestRepeatingRun(t *testing.T) {
	delay := time.Second
	repeat := 3

	now := time.Now()
	c := make(chan bool, 1)
	counter := 0
	task := scheduler.NewRepeatingTask(delay, repeat, func(params ...interface{}) {
		counter++
		c <- time.Since(now).Truncate(time.Second) == time.Duration(int(time.Second))
		now = time.Now()
	})
	scheduler.Schedule("test repeating", task)
	for {
		assert.True(t, <-c)
		if counter == repeat {
			break
		}
	}
}
