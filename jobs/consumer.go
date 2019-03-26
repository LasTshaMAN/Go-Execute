package jobs

func Consume(jobQueue <-chan func()) {
	if jobQueue == nil {
		panic("can't consume nil job queue")
	}
	go func() {
		for {
			select {
			case function, ok := <-jobQueue:
				if !ok {
					return
				}
				function()
			}
		}
	}()
}
