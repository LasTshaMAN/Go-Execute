package jobs

import "fmt"

// Executor is designed to free its User from mixing his business logic with necessary management of go-routines.
//
// Executor has a fixed amount of workers - go-routines that execute the actual work (you can specify their amount during Executor construction).
// Executor accepts jobs for execution in the form of functions(and function arguments, if they are required for function to work) that you can enqueue for execution.
// Enqueued function will eventually be executed.
// The order enqueued functions are run in is simply the order of executor.Enqueue(...) calls.
// Jobs can be executed in parallel - and they will, if you specify workersAmount to be > 1.
//
// Executor can be used in multi-threaded environment (you can call its methods from different go-routines in parallel and expect it to work correctly).
type Executor struct {
	jobQueue chan func()
}

// NewExecutor returns a new Executor object for you to run your jobs against.
// queueSize - specifies, how many jobs executor can guarantee to hold at any given time. You won't be able to enqueue new job for execution if job queue gets full. queueSize must be > 0.
// workersAmount - specifies, how many workers(go-routines) executor will use to handle jobs, sent for execution. Executor will run its workers in parallel. workersAmount must be > 0.
func NewExecutor(queueSize int, workersAmount int) *Executor {
	if queueSize < 1 {
		panic("queue size must be a positive number")
	}
	if workersAmount < 1 {
		panic("amount of workers must be a positive number")
	}

	workers := make([]*worker, workersAmount)
	for i := 0; i < workersAmount; i++ {
		workers = append(workers, NewWorker())
	}

	jobQueue := make(chan func(), queueSize)
	for _, worker := range workers {
		worker.Consume(jobQueue)
	}

	return &Executor{
		jobQueue: jobQueue,
	}
}

// Enqueue is a means to schedule a job for running through Executor.
// Enqueue returns an error if there already are too many jobs for Executor to handle at the moment. If Enqueue does return a error, you can try to enqueue (and succeed) your job at some time in the future.
// function - any standard Golang function - a unit of work that will be scheduled for execution. function cannot be 'nil'.
// args - arguments that function needs (could be none at all) in order to execute properly. args must match function signature.
func (executor *Executor) Enqueue(function func()) error {
	if function == nil {
		panic("cannot enqueue 'nil' function for execution")
	}
	if len(executor.jobQueue) == cap(executor.jobQueue) {
		return fmt.Errorf("executor queue is full at the moment")
	}

	executor.jobQueue <- function

	return nil
}
