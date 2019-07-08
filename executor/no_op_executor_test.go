package executor_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/LasTshaMAN/Go-Execute/executor"
)

func TestExecutorNoOp(t *testing.T) {
	t.Run("creating new executor", func(t *testing.T) {
		exec := executor.New(0)

		require.NotNil(t, exec)
	})
}

func TestEnqueueNoOp(t *testing.T) {
	t.Run("should be able to enqueue simple function for running", func(t *testing.T) {
		exec := executor.New(0)

		exec.Enqueue(func() {})
	})
	t.Run("should be able to enqueue 'nil' function for running", func(t *testing.T) {
		exec := executor.New(0)

		exec.Enqueue(nil)
	})
	t.Run("enqueued fn should not be executed", func(t *testing.T) {
		exec := executor.New(0)

		exec.Enqueue(func() {
			t.Fatal("this fn shouldn't be executed")
		})
	})
	t.Run("enqueued fn should not wait for a free worker before execution", func(t *testing.T) {
		exec := executor.New(0)
		exec.Enqueue(func() {
		})
		exec.Enqueue(func() {
		})
	})
	t.Run("should not execute any of the enqueued functions", func(t *testing.T) {
		const jobsAmount = 16

		exec := executor.New(0)
		for i := 0; i < jobsAmount; i++ {
			exec.Enqueue(func() {
				t.Fatal("this fn shouldn't be executed")
			})
		}
	})
}

func TestTryEnqueueNoOp(t *testing.T) {
	t.Run("should be able to enqueue simple function for running", func(t *testing.T) {
		exec := executor.New(0)

		err := exec.TryEnqueue(func() {})
		require.NoError(t, err)
	})
	t.Run("should be able to enqueue 'nil' function for running", func(t *testing.T) {
		exec := executor.New(0)

		err := exec.TryEnqueue(nil)
		require.NoError(t, err)
	})
	t.Run("enqueued fn should not be executed", func(t *testing.T) {
		exec := executor.New(0)

		err := exec.TryEnqueue(func() {
			t.Fatal("this fn shouldn't be executed")
		})
		require.NoError(t, err)
	})
	t.Run("enqueued fn should not wait for a free worker before execution", func(t *testing.T) {
		exec := executor.New(0)

		exec.Enqueue(func() {
		})
		err := exec.TryEnqueue(func() {
		})
		require.NoError(t, err)
	})
	t.Run("should not execute any of the enqueued functions", func(t *testing.T) {
		const jobsAmount = 16

		exec := executor.New(0)

		for i := 0; i < jobsAmount; i++ {
			err := exec.TryEnqueue(func() {
				t.Fatal("this fn shouldn't be executed")
			})
			require.NoError(t, err)
		}
	})
}
