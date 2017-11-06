package jobs

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"time"
)

func TestBasicLifecycle(t *testing.T) {
	t.Run("should be able to create Executor", func(t *testing.T) {
		executor := NewExecutor(4, 4)

		assert.NotNil(t, executor)
	})

	t.Run("shouldn't be able to create Executor with 0 queue size", func(t *testing.T) {
		assert.Panics(t, func() {
			NewExecutor(0, 4)
		})
	})

	t.Run("shouldn't be able to create Executor with 0 amount of workers", func(t *testing.T) {
		assert.Panics(t, func() {
			NewExecutor(4, 0)
		})
	})
}

func TestJobEnqueueing(t *testing.T) {
	t.Run("should be able to enqueue simple function for running", func(t *testing.T) {
		executor := NewExecutor(4, 4)

		err := executor.EnqueueAsync(func() {})

		assert.NoError(t, err)
	})

	t.Run("should be able to enqueue function with args for running", func(t *testing.T) {
		executor := NewExecutor(4, 4)

		someFunction := func(string, int) {}
		err := executor.EnqueueAsync(func() {
			someFunction("one", 1)
		})

		assert.NoError(t, err)
	})

	t.Run("shouldn't be able to enqueue 'nil' function for running", func(t *testing.T) {
		executor := NewExecutor(4, 4)

		assert.Panics(t, func() {
			executor.EnqueueAsync(nil)
		})
	})

	t.Run("shouldn't be able to enqueue function when Executor's queue is full", func(t *testing.T) {
		executor := NewExecutor(4, 4)

		for i := 0; i < 4; i++ {
			err := executor.EnqueueAsync(func() {})
			assert.NoError(t, err)
		}

		err := executor.EnqueueAsync(func() {})
		assert.Error(t, err)
	})
}

func TestJobExecution(t *testing.T) {
	t.Run("should execute enqueued function eventually", func(t *testing.T) {
		executor := NewExecutor(4, 4)

		out := make(chan bool)
		err := executor.EnqueueAsync(func() {
			out <- true
		})

		assert.NoError(t, err)
		ok := <-out
		assert.True(t, ok)
	})

	t.Run("should execute multiple enqueued functions eventually", func(t *testing.T) {
		executor := NewExecutor(16, 4)

		jobsAmount := 16
		out := make(chan bool)
		function := func() {
			out <- true
		}
		for i := 0; i < jobsAmount; i++ {
			err := executor.EnqueueAsync(function)
			assert.NoError(t, err)
		}

		for i := 0; i < jobsAmount; i++ {
			ok := <-out
			assert.True(t, ok)
		}
	})

	t.Run("shouldn't be able to enqueue function while workers are busy and queue is full", func(t *testing.T) {
		executor := NewExecutor(2, 3)

		toughJob := func() {
			select {}
		}
		for i := 0; i < 5; {
			err := executor.EnqueueAsync(toughJob)
			if err == nil {
				time.Sleep(time.Millisecond)
			}
			i++
		}

		err := executor.EnqueueAsync(toughJob)
		assert.Error(t, err)
	})
}
