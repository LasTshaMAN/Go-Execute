package executor_test

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/LasTshaMAN/Go-Execute/executor"
)

func TestExecutor(t *testing.T) {
	t.Run("creating new executor", func(t *testing.T) {
		exec := executor.New(1)

		require.NotNil(t, exec)
	})
}

func TestEnqueue(t *testing.T) {
	t.Run("should be able to enqueue simple function for running", func(t *testing.T) {
		exec := executor.New(1)

		exec.Enqueue(func() {})
	})
	t.Run("should be able to enqueue 'nil' function for running", func(t *testing.T) {
		exec := executor.New(1)

		exec.Enqueue(nil)
	})
	t.Run("enqueued fn should be executed", func(t *testing.T) {
		exec := executor.New(1)

		result := make(chan struct{})
		exec.Enqueue(func() {
			result <- struct{}{}
		})

		<-result
	})
	t.Run("enqueued fn should wait for a free worker before execution", func(t *testing.T) {
		jobDone := make(chan struct{})

		exec := executor.New(1)
		exec.Enqueue(func() {
			close(jobDone)
		})
		exec.Enqueue(func() {
		})

		if _, ok := <-jobDone; ok {
			require.Fail(t, "jobDone channel should've already been closed")
		}
	})
	t.Run("should execute all the enqueued functions eventually", func(t *testing.T) {
		const jobsAmount = 16

		exec := executor.New(4)

		result := make(chan struct{}, jobsAmount)
		for i := 0; i < jobsAmount; i++ {
			exec.Enqueue(func() {
				result <- struct{}{}
			})
		}
		for i := 0; i < jobsAmount; i++ {
			<-result
		}
	})
}

func TestTryEnqueue(t *testing.T) {
	t.Run("should be able to enqueue simple function for running", func(t *testing.T) {
		exec := executor.New(1)

		err := exec.TryEnqueue(func() {})
		require.NoError(t, err)
	})
	t.Run("should be able to enqueue 'nil' function for running", func(t *testing.T) {
		exec := executor.New(1)

		err := exec.TryEnqueue(nil)
		require.NoError(t, err)
	})
	t.Run("enqueued fn should be executed", func(t *testing.T) {
		exec := executor.New(1)

		result := make(chan struct{})
		err := exec.TryEnqueue(func() {
			result <- struct{}{}
		})
		require.NoError(t, err)

		<-result
	})
	t.Run("enqueued fn should wait for a free worker before execution", func(t *testing.T) {
		exec := executor.New(1)

		jobDone := make(chan struct{})
		exec.Enqueue(func() {
			<-jobDone
		})
		err := exec.TryEnqueue(func() {
		})
		require.Error(t, err)
	})
	t.Run("should execute all the enqueued functions eventually", func(t *testing.T) {
		const jobsAmount = 16

		exec := executor.New(4)

		result := make(chan struct{}, 16)
		i := 0
		for i < jobsAmount {
			err := exec.TryEnqueue(func() {
				result <- struct{}{}
			})
			if err != nil {
				continue
			}
			i++
		}
		for i := 0; i < jobsAmount; i++ {
			<-result
		}
	})
}

func TestWait(t *testing.T) {
	t.Run("should execute all the enqueued functions eventually", func(t *testing.T) {
		exec := executor.New(4)

		result := make(chan struct{}, 2)

		exec.Enqueue(func() {
			// Imitate some workload
			time.Sleep(time.Millisecond)

			result <- struct{}{}
		})
		exec.Enqueue(func() {
			// Imitate some workload
			time.Sleep(time.Millisecond)

			result <- struct{}{}
		})

		exec.Wait()

		require.Len(t, result, 2)
	})
}
