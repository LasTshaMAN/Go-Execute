package jobs

import (
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestWorker(t *testing.T) {
	t.Run("should be able to create a Worker", func(t *testing.T) {
		worker := NewWorker()

		assert.NotNil(t, worker)
	})

	t.Run("should refuse to consume nil job queue", func(t *testing.T) {
		worker := NewWorker()
		var jobQueue chan func()

		assert.Panics(t, func() {
			worker.Consume(jobQueue)
		})
	})

	t.Run("should be able to consume empty job queue", func(t *testing.T) {
		worker := NewWorker()
		jobQueue := make(chan func())

		assert.NotPanics(t, func() {
			worker.Consume(jobQueue)
		})
	})

	t.Run("should take jobs from queue", func(t *testing.T) {
		worker := NewWorker()
		jobQueue := make(chan func(), 1)
		jobQueue <- func() {}

		worker.Consume(jobQueue)

		// Should not block
		jobQueue <- func() {}
	})

	t.Run("should terminate gracefully when job queue is closed", func(t *testing.T) {
		worker := NewWorker()
		jobQueue := make(chan func(), 1)
		jobQueue <- func() {}

		worker.Consume(jobQueue)

		close(jobQueue)
		time.Sleep(200 * time.Millisecond)
	})

	t.Run("should execute jobs taken from queue", func(t *testing.T) {
		worker := NewWorker()
		jobQueue := make(chan func(), 1)
		out := make(chan bool)
		jobQueue <- func() {
			out <- true
		}

		worker.Consume(jobQueue)

		ok := <-out
		assert.True(t, ok)
	})

	t.Run("should execute jobs taken from queue", func(t *testing.T) {
		worker := NewWorker()
		jobQueue := make(chan func(), 1)
		out := make(chan bool)
		jobQueue <- func() {
			out <- true
		}

		worker.Consume(jobQueue)

		ok := <-out
		assert.True(t, ok)
	})

	t.Run("should execute only one job at a time according to queue's order", func(t *testing.T) {
		worker := NewWorker()
		jobQueue := make(chan func(), 4)
		firstJobFinished := false
		jobQueue <- func() {
			time.Sleep(200 * time.Millisecond)
			firstJobFinished = true
		}
		out := make(chan bool)
		jobQueue <- func() {
			out <- true
		}

		worker.Consume(jobQueue)

		ok := <-out
		assert.True(t, ok)
		assert.True(t, firstJobFinished)
	})
}
