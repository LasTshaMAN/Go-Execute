package jobs

import "fmt"

type Executor struct {
	started  bool
	jobQueue chan *Job
}

func NewExecutor(queueSize int) *Executor {
	if queueSize < 1 {
		panic("queue size must be a positive number")
	}
	return &Executor{
		jobQueue: make(chan *Job, queueSize),
	}
}

func (executor *Executor) Start() {
	if executor.started {
		panic("couldn't start already running Executor!")
	}
	executor.started = true

	executor.consumeJobQueue()
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

func (executor *Executor) consumeJobQueue() {
	go func() {
		for {
			select {
			case job := <-executor.jobQueue:
				go job.Execute()
			}
		}
	}()
}
