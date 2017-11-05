package jobs

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

// TODO - write tests for using Executor in multi-threaded environment

func TestLifecycleConcurrent(t *testing.T) {
	t.Run("shouldn't be able to start already running Executor in parallel", func(t *testing.T) {
		executor := NewExecutor(4, 4)
		firstStarted := startExecutorInBackground(executor)

		secondStarted := startExecutorInBackground(executor)

		successfulStarts := 0
		ok := <-firstStarted
		if ok {
			successfulStarts++
		}
		ok = <-secondStarted
		if ok {
			successfulStarts++
		}
		assert.Equal(t, successfulStarts, 1)
	})
}

func startExecutorInBackground(executor *Executor) chan bool {
	success := make(chan bool)

	go func() {
		defer func() {
			if r := recover(); r != nil {
				success <- false
			}
		}()

		executor.Start()

		success <- true
	}()

	return success
}
