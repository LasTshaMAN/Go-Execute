// Package jobs implements "Thread-Pool" design pattern - https://en.wikipedia.org/wiki/Thread_pool.
//
// Its main purpose is to decouple business logic from the logic necessary for go-routines management.
package jobs

import "fmt"

// Executor has a fixed amount of workers(go-routines) that execute the actual work.
//
// Executor accepts simple, yet flexible function (func() {}) that you can enqueue for execution.
//
// Enqueued function will eventually be executed.
// Functions are run in the order they were enqueued.
// Functions will be executed in parallel if you specify workersAmount to be > 1.
//
// Make sure that the function you enqueue for execution won't block forever (e.g. writing in channel that won't ever be read from).
// It will cause the corresponding worker to hang forever - thus leaking resources.
//
// Executor is safe to use in multi-threaded environment (you can enqueue functions from different go-routines and expect Executor to work correctly).
type Executor struct {
	jobQueue chan func()
}

// NewExecutor returns a new Executor object - a means to enqueue your functions for execution.
//
// workersAmount - specifies, how many workers(go-routines) Executor will use to handle functions, sent for execution. Executor will run its workers in parallel.
//
// queueSize - specifies, how many functions executor can easily hold to (either keeping them in queue or executing them) at any given time.
func NewExecutor(workersAmount uint, queueSize uint) *Executor {
	jobQueue := make(chan func(), queueSize)
	for i := uint(0); i < workersAmount; i++ {
		consume(jobQueue)
	}
	return &Executor{
		jobQueue: jobQueue,
	}
}

// TryToEnqueue is a way for you to schedule a function for execution.
// Enqueued function will eventually be executed at some point in the future.
//
// TryToEnqueue call doesn't block.
// TryToEnqueue returns an error if there already are too many functions for Executor to handle at the moment.
// If TryToEnqueue does return an error, you can try to enqueue your function (and succeed) at some point in the future.
//
// function - simple Golang function - a unit of work that will be scheduled for execution as soon as there is a free worker to tackle it.
// function cannot be 'nil'.
//
// TryToEnqueue is safe to use in multi-threaded environment
func (executor *Executor) TryToEnqueue(function func()) error {
	if len(executor.jobQueue) == cap(executor.jobQueue) {
		return fmt.Errorf("executor queue is full at the moment")
	}

	executor.Enqueue(function)

	return nil
}

// Enqueue is a way for you to schedule a function for execution.
// Enqueued function will eventually be executed at some point in the future.
//
// Enqueue call blocks until Executor is ready to accept the function you are trying to enqueue.
//
// function - simple Golang function - a unit of work that will be scheduled for execution as soon as there is a free worker to tackle it.
// function cannot be 'nil'.
//
// Enqueue is safe to use in multi-threaded environment
func (executor *Executor) Enqueue(function func()) {
	if function == nil {
		return
	}
	executor.jobQueue <- function
}
