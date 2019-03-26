package jobs

import "errors"

func consume(jobQueue <-chan func()) error {
	if jobQueue == nil {
		return errors.New("can't consume nil job queue")
	}

	go func() {
		for function := range jobQueue {
			function()
		}
	}()

	return nil
}
