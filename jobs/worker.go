package jobs

import "time"

type Worker struct {
}

func NewWorker() *Worker {
	return &Worker{}
}

func (worker *Worker) Consume(jobQueue chan *Job) {
	if jobQueue == nil {
		panic("can't consume nil job queue")
	}
	go func() {
		for {
			select {
			case job, ok := <-jobQueue:
				if !ok {
					return
				}
				job.Execute()
			default:
				// Saving CPU resources
				time.Sleep(time.Millisecond * 100)
			}
		}
	}()
}
