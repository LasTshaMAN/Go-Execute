package jobs

func consume(jobQueue <-chan func()) {
	if jobQueue == nil {
		panic("can't consume nil job queue")
	}
	go func() {
		for function := range jobQueue {
			function()
		}
	}()
}
