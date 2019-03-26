package jobs

func Consume(jobQueue <-chan func()) {
	if jobQueue == nil {
		panic("can't consume nil job queue")
	}
	go func() {
		for function := range jobQueue {
			function()
		}
	}()
}
