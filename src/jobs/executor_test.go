package jobs

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"time"
)

func TestLifecycle(t *testing.T) {
	t.Run("should be able to create Executor", func(t *testing.T) {
		executor := NewExecutor(4)

		assert.NotNil(t, executor)
	})

	t.Run("shouldn't be able to create Executor with 0 queue size", func(t *testing.T) {
		assert.Panics(t, func() {
			NewExecutor(0)
		})
	})

	t.Run("should be able to start Executor", func(t *testing.T) {
		executor := NewExecutor(4)

		assert.NotPanics(t, executor.Start)
	})

	t.Run("shouldn't be able to start already running Executor", func(t *testing.T) {
		executor := NewExecutor(4)
		executor.Start()

		assert.Panics(t, executor.Start)
	})

	t.Run("shouldn't be able to stop non-running Executor", func(t *testing.T) {
		executor := NewExecutor(4)

		assert.Panics(t, executor.Stop)
	})

	t.Run("should be able to stop running Executor", func(t *testing.T) {
		executor := NewExecutor(4)
		executor.Start()

		assert.NotPanics(t, executor.Stop)
	})

	t.Run("should be able to restart Executor after it was stopped", func(t *testing.T) {
		executor := NewExecutor(4)
		executor.Start()
		executor.Stop()

		assert.NotPanics(t, executor.Start)
	})
}

func TestJobExecution(t *testing.T) {
	t.Run("should be able to enqueue simple function for running", func(t *testing.T) {
		executor := NewExecutor(4)

		err := executor.Enqueue(func() {})

		assert.NoError(t, err)
	})

	t.Run("should be able to enqueue function with args for running", func(t *testing.T) {
		executor := NewExecutor(4)

		err := executor.Enqueue(func(string, int) {}, "one", 1)

		assert.NoError(t, err)
	})

	t.Run("should panic when supplied args do not match enqueued function signature", func(t *testing.T) {
		executor := NewExecutor(4)

		assert.Panics(t, func() {
			executor.Enqueue(func(int, string) {}, "one", 1)
		})
	})

	t.Run("should panic when too many args were supplied", func(t *testing.T) {
		executor := NewExecutor(4)

		assert.Panics(t, func() {
			executor.Enqueue(func(int, string) {}, 1, "one", "excessive")
		})
	})

	t.Run("should panic when too few args were supplied", func(t *testing.T) {
		executor := NewExecutor(4)

		assert.Panics(t, func() {
			executor.Enqueue(func(int, string) {}, 1)
		})
	})

	t.Run("shouldn't be able to enqueue function when Executor's queue is full", func(t *testing.T) {
		executor := NewExecutor(4)

		for i := 0; i < 4; i++ {
			err := executor.Enqueue(func() {})
			assert.NoError(t, err)
		}

		err := executor.Enqueue(func() {})
		assert.Error(t, err)
	})

	t.Run("shouldn't execute enqueued function while Executor is not running", func(t *testing.T) {
		executor := NewExecutor(4)

		out := make(chan bool)
		err := executor.Enqueue(func() {
			close(out)
		})

		assert.NoError(t, err)
		time.Sleep(time.Millisecond)
		select {
		case <-out:
			assert.FailNow(t, "enqueued function was executed")
		default:
		}
	})

	t.Run("should execute enqueued function when Executor is running", func(t *testing.T) {
		executor := NewExecutor(4)
		executor.Start()

		out := make(chan bool)
		err := executor.Enqueue(func() {
			out <- true
		})

		assert.NoError(t, err)
		ok := <-out
		assert.True(t, ok)
	})

	t.Run("should execute enqueued function when Executor is started after the function was enqueued", func(t *testing.T) {
		executor := NewExecutor(4)

		out := make(chan bool)
		err := executor.Enqueue(func() {
			out <- true
		})
		executor.Start()

		assert.NoError(t, err)
		ok := <-out
		assert.True(t, ok)
	})

	t.Run("should execute multiple enqueued functions when Executor is running", func(t *testing.T) {
		executor := NewExecutor(16)
		executor.Start()

		jobsAmount := 16
		out := make(chan bool)
		function := func() {
			out <- true
		}
		for i := 0; i < jobsAmount; i++ {
			err := executor.Enqueue(function)
			assert.NoError(t, err)
		}

		for i := 0; i < jobsAmount; i++ {
			ok := <-out
			assert.True(t, ok)
		}
	})

	t.Run("should execute multiple enqueued functions when Executor is started after the function was enqueued", func(t *testing.T) {
		executor := NewExecutor(16)

		jobsAmount := 16
		out := make(chan bool)
		function := func() {
			out <- true
		}
		for i := 0; i < jobsAmount; i++ {
			err := executor.Enqueue(function)
			assert.NoError(t, err)
		}
		executor.Start()

		for i := 0; i < jobsAmount; i++ {
			ok := <-out
			assert.True(t, ok)
		}
	})
}
