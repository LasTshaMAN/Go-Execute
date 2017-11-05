package jobs

import "fmt"

type Executor struct {
	started  bool
	jobQueue chan *Job
	workers  []*Worker
}

func NewExecutor(queueSize int, workersAmount int) *Executor {
	if queueSize < 1 {
		panic("queue size must be a positive number")
	}
	if workersAmount < 1 {
		panic("amount of workers must be a positive number")
	}

	workers := make([]*Worker, workersAmount)
	for i := 0; i < workersAmount; i++ {
		workers = append(workers, NewWorker())
	}

	return &Executor{
		jobQueue: make(chan *Job, queueSize),
		workers:  workers,
	}
}

func (executor *Executor) Start() {
	if executor.started {
		panic("couldn't start already running Executor!")
	}
	executor.started = true

	for _, worker := range executor.workers {
		worker.Consume(executor.jobQueue)
	}
}

func (executor *Executor) Stop() {
	if !executor.started {
		panic("couldn't stop non-running Executor!")
	}
	executor.started = false
}

func (executor *Executor) Enqueue(function interface{}, args ...interface{}) error {
	if len(executor.jobQueue) == cap(executor.jobQueue) {
		return fmt.Errorf("executor queue is full at the moment")
	}

	job := NewJob(function, args...)
	executor.jobQueue <- job

	return nil
}
