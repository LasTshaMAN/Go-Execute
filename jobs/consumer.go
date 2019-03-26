package jobs

func consume(jobQueue <-chan func()) {
	if jobQueue == nil {
		return
	}

	go func() {
		for function := range jobQueue {
			function()
		}
	}()
}
