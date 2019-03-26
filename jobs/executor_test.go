package jobs_test

import (
	"github.com/LasTshaMAN/Go-Execute/jobs"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBasicLifecycle(t *testing.T) {
	t.Run("should be able to create Executor", func(t *testing.T) {
		executor := jobs.NewExecutor(4, 4)

		assert.NotNil(t, executor)
	})

	t.Run("shouldn't be able to create Executor with 0 amount of workers", func(t *testing.T) {
		assert.Panics(t, func() {
			jobs.NewExecutor(0, 4)
		})
	})
	
	t.Run("shouldn't be able to create Executor with 0 queue size", func(t *testing.T) {
		assert.Panics(t, func() {
			jobs.NewExecutor(4, 0)
		})
	})
}

func TestJobEnqueueing(t *testing.T) {
	t.Run("should be able to enqueue simple function for running", func(t *testing.T) {
		executor := jobs.NewExecutor(4, 4)

		executor.Enqueue(func() {})
		err := executor.TryToEnqueue(func() {})

		assert.NoError(t, err)
	})

	t.Run("should be able to enqueue function with args for running", func(t *testing.T) {
		executor := jobs.NewExecutor(4, 4)

		someFunction := func(string, int) {}
		executor.Enqueue(func() {
			someFunction("one", 1)
		})
		err := executor.TryToEnqueue(func() {
			someFunction("one", 1)
		})

		assert.NoError(t, err)
	})

	t.Run("shouldn't be able to enqueue 'nil' function for running", func(t *testing.T) {
		executor := jobs.NewExecutor(4, 4)

		assert.Panics(t, func() {
			executor.Enqueue(nil)
		})
		assert.Panics(t, func() {
			executor.TryToEnqueue(nil)
		})
	})

	t.Run("should refuse an attempt to enqueue function when Executor's queue is full", func(t *testing.T) {
		executor := jobs.NewExecutor(4, 4)

		toughJob := func() {
			select {}
		}
		for i := 0; i < 8; i++ {
			executor.Enqueue(toughJob)
		}

		err := executor.TryToEnqueue(func() {})
		assert.Error(t, err)
	})
}

func TestJobExecution(t *testing.T) {
	t.Run("should execute enqueued function eventually", func(t *testing.T) {
		executor := jobs.NewExecutor(4, 4)

		out := make(chan bool)
		executor.Enqueue(func() {
			out <- true
		})

		ok := <-out
		assert.True(t, ok)
	})

	t.Run("should execute multiple enqueued functions eventually", func(t *testing.T) {
		executor := jobs.NewExecutor(4, 4)

		jobsAmount := 16
		out := make(chan bool, 8)
		function := func() {
			out <- true
		}
		for i := 0; i < jobsAmount; i++ {
			executor.Enqueue(function)
		}

		for i := 0; i < jobsAmount; i++ {
			ok := <-out
			assert.True(t, ok)
		}
	})
}
