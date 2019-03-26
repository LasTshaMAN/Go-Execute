package jobs_test

import (
	"testing"

	"github.com/LasTshaMAN/Go-Execute/jobs"
	"github.com/stretchr/testify/require"
)

func TestBasicLifecycle(t *testing.T) {
	t.Run("should be able to create Executor", func(t *testing.T) {
		executor := jobs.NewExecutor(4, 4)

		require.NotNil(t, executor)
	})
	t.Run("shouldn't be able to create Executor with 0 amount of workers and 0 queue size", func(t *testing.T) {
		executor := jobs.NewExecutor(0, 0)

		require.NotNil(t, executor)
	})
}

func TestJobEnqueueing(t *testing.T) {
	t.Run("should be able to enqueue simple function for running", func(t *testing.T) {
		executor := jobs.NewExecutor(4, 4)

		executor.Enqueue(func() {})
		err := executor.TryToEnqueue(func() {})

		require.NoError(t, err)
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

		require.NoError(t, err)
	})

	t.Run("shouldn't be able to enqueue 'nil' function for running", func(t *testing.T) {
		executor := jobs.NewExecutor(4, 4)

		require.NotPanics(t, func() {
			executor.Enqueue(nil)
		})
		require.NotPanics(t, func() {
			executor.TryToEnqueue(nil)
		})
	})

	t.Run("should refuse an attempt to enqueue function when Executor's queue is full", func(t *testing.T) {
		t.Run("zero workers, zero queue size case", func(t *testing.T) {
			executor := jobs.NewExecutor(0, 0)

			err := executor.TryToEnqueue(func() {})
			require.Error(t, err)
		})
		t.Run("basic case", func(t *testing.T) {
			executor := jobs.NewExecutor(4, 4)

			toughJob := func() {
				select {}
			}
			for i := 0; i < 8; i++ {
				executor.Enqueue(toughJob)
			}

			err := executor.TryToEnqueue(func() {})
			require.Error(t, err)
		})
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
		require.True(t, ok)
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
			require.True(t, ok)
		}
	})
}
