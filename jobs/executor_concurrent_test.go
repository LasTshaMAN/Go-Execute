package jobs_test

import (
	"github.com/LasTshaMAN/Go-Execute/jobs"
	"github.com/stretchr/testify/assert"
	"sync"
	"testing"
)

func TestJobEnqueueingConcurrent(t *testing.T) {
	t.Run("should be able to enqueue function for running in concurrent environment", func(t *testing.T) {
		executor := jobs.NewExecutor(4, 4)

		var wg sync.WaitGroup
		wg.Add(4)

		for i := 0; i < 4; i++ {
			go func() {
				executor.Enqueue(func() {})
				err := executor.TryToEnqueue(func() {})
				assert.NoError(t, err)
				wg.Done()
			}()
		}

		wg.Wait()
	})

	t.Run("should refuse an attempt to enqueue function when Executor's queue is full in concurrent environment", func(t *testing.T) {
		executor := jobs.NewExecutor(4, 4)

		var wg sync.WaitGroup
		wg.Add(8)

		for i := 0; i < 8; i++ {
			go func() {
				toughJob := func() {
					select {}
				}
				executor.Enqueue(toughJob)
				wg.Done()
			}()
		}

		wg.Wait()
		err := executor.TryToEnqueue(func() {})
		assert.Error(t, err)
	})
}

func TestJobExecutionConcurrent(t *testing.T) {
	t.Run("should eventually execute functions enqueued in concurrent environment", func(t *testing.T) {
		executor := jobs.NewExecutor(4, 4)

		out := make(chan bool)

		for i := 0; i < 8; i++ {
			go func() {
				executor.Enqueue(func() {
					out <- true
				})
			}()
		}

		for i := 0; i < 8; i++ {
			ok := <-out
			assert.True(t, ok)
		}
	})
}
