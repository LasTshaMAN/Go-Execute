package jobs

import "time"

type worker struct {
}

func newWorker() *worker {
	return &worker{}
}

func (worker *worker) Consume(jobQueue chan func()) {
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
			default:
				// Saving CPU resources
				time.Sleep(time.Millisecond * 100)
			}
		}
	}()
}
