package jobs

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestJob(t *testing.T) {
	t.Run("should be able to create a Job", func(t *testing.T) {
		job := NewJob(func() {})

		assert.NotNil(t, job)
	})

	t.Run("should panic when args, supplied during Job creation, do not match function signature", func(t *testing.T) {
		assert.Panics(t, func() {
			NewJob(func(int, string) {}, "one", 1)
		})
	})

	t.Run("should panic when too many args were supplied", func(t *testing.T) {
		assert.Panics(t, func() {
			NewJob(func(int, string) {}, 1, "one", "excessive")
		})
	})

	t.Run("should panic when too few args were supplied", func(t *testing.T) {
		assert.Panics(t, func() {
			NewJob(func(int, string) {}, 1)
		})
	})

	t.Run("should panic when args, supplied during Job creation, do not match function signature", func(t *testing.T) {
		assert.Panics(t, func() {
			NewJob(func(int, string) {}, "one", 1)
		})
	})

	t.Run("should be able to execute a Job", func(t *testing.T) {
		var ok bool
		job := NewJob(func() {
			ok = true
		})

		job.Execute()

		assert.True(t, ok)
	})

	t.Run("args should be passed correctly during Job execution", func(t *testing.T) {
		var ok bool
		var integer int
		var str string
		job := NewJob(func(okArg bool, integerArg int, strArg string) {
			ok = okArg
			integer = integerArg
			str = strArg
		}, true, 123, "str")

		job.Execute()

		assert.True(t, ok)
		assert.Equal(t, integer, 123)
		assert.Equal(t, str, "str")
	})
}
