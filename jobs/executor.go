package jobs

import "fmt"

// Executor is an implementation of the "Thread-Pool" design pattern. Its main purpose is to decouple business logic from the logic necessary for Go-routines management.
//
// Executor has a fixed amount of workers - Go-routines that execute the actual work (you can specify their amount during Executor construction).
// Executor accepts simple, yet flexible function (func() {}) that you can enqueue for execution.
// Enqueued function will eventually be executed.
// The order enqueued functions are run in is simply the order of executor.EnqueueAsync(func() {}) calls.
// Functions can be executed in parallel - and they will, if you specify workersAmount to be > 1.
//
// Make sure that the function you enqueue for execution won't block forever (e.g. writing in channel that won't ever be read from).
// It will cause the corresponding worker to hang forever - thus leaking resources.
//
// Executor can be used in multi-threaded environment (you can enqueue functions from different Go-routines and expect Executor to work correctly).
type Executor struct {
	jobQueue chan func()
}

// NewExecutor returns a new Executor object - a means for you to enqueue your functions.
//
// workersAmount - specifies, how many workers(go-routines) Executor will use to handle functions, sent for execution. Executor will run its workers in parallel. workersAmount must be > 0.
// queueSize - specifies, how many functions executor can easily hold (either in queue or executing them) at any given time. queueSize must be > 0.
func NewExecutor(workersAmount int, queueSize int) *Executor {
	if queueSize < 1 {
		panic("queue size must be a positive number")
	}
	if workersAmount < 1 {
		panic("amount of workers must be a positive number")
	}

	workers := make([]*worker, 0, workersAmount)
	for i := 0; i < workersAmount; i++ {
		workers = append(workers, newWorker())
	}

	jobQueue := make(chan func(), queueSize)
	for _, worker := range workers {
		worker.Consume(jobQueue)
	}

	return &Executor{
		jobQueue: jobQueue,
	}
}

// EnqueueAsync is a way for you to schedule a function for execution.
// Enqueued function will eventually be executed at some point in the future.
//
// EnqueueAsync call doesn't block.
// EnqueueAsync returns an error if there already are too many functions for Executor to handle at the moment.
// If EnqueueAsync does return an error, you can try to enqueue your function (and succeed) at some time in the future.
//
// function - Golang function of the form "func() {}" - a unit of work that will be scheduled for execution as soon as there is a free worker to tackle it.
// function cannot be 'nil'.
func (executor *Executor) EnqueueAsync(function func()) error {
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
// function - Golang function of the form "func() {}" - a unit of work that will be scheduled for execution as soon as there is a free worker to tackle it.
// function cannot be 'nil'.
func (executor *Executor) Enqueue(function func()) {
	if function == nil {
		panic("cannot enqueue 'nil' function for execution")
	}
	executor.jobQueue <- function
}
