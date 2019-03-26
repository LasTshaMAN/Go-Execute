package jobs

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestWorker(t *testing.T) {
	t.Run("should refuse to consume nil job queue", func(t *testing.T) {
		var jobQueue chan func()

		assert.Panics(t, func() {
			Consume(jobQueue)
		})
	})

	t.Run("should be able to consume empty job queue", func(t *testing.T) {
		jobQueue := make(chan func())

		assert.NotPanics(t, func() {
			Consume(jobQueue)
		})
	})

	t.Run("should take jobs from queue", func(t *testing.T) {
		jobQueue := make(chan func(), 1)
		jobQueue <- func() {}

		Consume(jobQueue)

		// Should not block
		jobQueue <- func() {}
	})

	t.Run("should terminate gracefully when job queue is closed", func(t *testing.T) {
		jobQueue := make(chan func(), 1)
		jobQueue <- func() {}

		Consume(jobQueue)

		close(jobQueue)
		time.Sleep(200 * time.Millisecond)
	})

	t.Run("should execute jobs taken from queue", func(t *testing.T) {
		jobQueue := make(chan func(), 1)
		out := make(chan bool)
		jobQueue <- func() {
			out <- true
		}

		Consume(jobQueue)

		ok := <-out
		assert.True(t, ok)
	})

	t.Run("should execute jobs taken from queue", func(t *testing.T) {
		jobQueue := make(chan func(), 1)
		out := make(chan bool)
		jobQueue <- func() {
			out <- true
		}

		Consume(jobQueue)

		ok := <-out
		assert.True(t, ok)
	})

	t.Run("should execute only one job at a time according to queue's order", func(t *testing.T) {
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

		Consume(jobQueue)

		ok := <-out
		assert.True(t, ok)
		assert.True(t, firstJobFinished)
	})
}
