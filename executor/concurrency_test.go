package executor_test

import (
	"sync/atomic"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/LasTshaMAN/Go-Execute/executor"
)

func TestEnqueueConcurrent(t *testing.T) {
	t.Run("enqueued fn should wait for a free worker before execution", func(t *testing.T) {
		enqueuingDone := make(chan struct{})
		jobDone := make(chan struct{})

		exec := executor.New(1)

		go func() {
			exec.Enqueue(func() {
				close(jobDone)
			})
			enqueuingDone <- struct{}{}
		}()

		<-enqueuingDone

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
			go exec.Enqueue(func() {
				result <- struct{}{}
			})
		}
		for i := 0; i < jobsAmount; i++ {
			<-result
		}
	})
}

func TestTryEnqueueConcurrent(t *testing.T) {
	t.Run("enqueued fn should wait for a free worker before execution", func(t *testing.T) {
		enqueuingDone := make(chan struct{})
		jobDone := make(chan struct{})

		exec := executor.New(1)

		go func() {
			exec.Enqueue(func() {
				<-jobDone
			})
			enqueuingDone <- struct{}{}
		}()

		<-enqueuingDone

		err := exec.TryEnqueue(func() {
		})
		require.Error(t, err)
	})
	t.Run("should execute all the enqueued functions eventually", func(t *testing.T) {
		const jobsAmount = 16

		exec := executor.New(4)

		result := make(chan struct{}, 16)
		i := uint32(0)
		for atomic.LoadUint32(&i) < jobsAmount {
			go func() {
				err := exec.TryEnqueue(func() {
					result <- struct{}{}
				})
				if err != nil {
					return
				}
				atomic.AddUint32(&i, 1)
			}()
		}
		for j := 0; j < jobsAmount; j++ {
			<-result
		}
	})
}

func TestWaitConcurrent(t *testing.T) {
	t.Run("should execute all the enqueued functions eventually", func(t *testing.T) {
		exec := executor.New(4)

		result := make(chan struct{}, 2)
		scheduled := make(chan struct{}, 2)

		go func() {
			exec.Enqueue(func() {
				// Imitate some workload
				time.Sleep(time.Millisecond)

				result <- struct{}{}
			})

			scheduled <- struct{}{}
		}()

		go func() {
			exec.Enqueue(func() {
				// Imitate some workload
				time.Sleep(time.Millisecond)

				result <- struct{}{}
			})

			scheduled <- struct{}{}
		}()

		// Make sure some of the jobs were scheduled from another go-routines first.
		<-scheduled

		exec.Wait()

		require.Len(t, result, 2)
	})
}
