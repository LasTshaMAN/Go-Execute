package executor

// Executor provides a simple interface for job-execution by a fixed-size pool of go-routines.
type Executor interface {

	// Enqueue provides a way to schedule a function fn for execution.
	// Enqueued fn will eventually be executed at some point in the future.
	// Enqueue call blocks until Executor is ready to accept the function you are trying to enqueue.
	//
	// fn - is a Golang function - a unit of work that will be scheduled for execution as soon as there is a free worker to tackle it.
	// If 'nil' is passed as fn, Executor silently throws it away.
	//
	// Make sure that the function fn you enqueue for execution doesn't block forever.
	// It will cause the corresponding worker to hang forever.
	//
	// A call to Enqueue method "Happens Before" func fn executes. All other bets with respect to Golang memory model are off.
	// That means all variables func fn has access to shouldn't be modified from any place (other than from func fn)
	// after func fn has been enqueued. Otherwise such modifications will result in a "racy" behavior by func fn.
	Enqueue(fn func())

	// TryEnqueue provides a way to schedule a function fn for execution.
	// Enqueued fn will eventually be executed at some point in the future.
	//
	// TryEnqueue call doesn't block.
	// TryEnqueue returns an error if there already are too many functions for Executor to handle at the moment.
	// If TryEnqueue does return an error, you can try to enqueue your fn (and succeed) at some point in the future.
	//
	// fn - is a Golang function - a unit of work that will be scheduled for execution as soon as there is a free worker to tackle it.
	// If 'nil' is passed as fn, Executor silently throws it away.
	//
	// Make sure that the function fn you enqueue for execution won't block forever.
	// It will cause the corresponding worker to hang forever.
	//
	// A call to TryEnqueue method "Happens Before" func fn executes. All other bets with respect to Golang memory model are off.
	// That means all variables func fn has access to shouldn't be modified from any place (other than from func fn)
	// after func fn has been enqueued. Otherwise such modifications will result in a "racy" behavior by func fn.
	TryEnqueue(fn func()) error

	// Wait blocks until all the workers are finished with all ongoing jobs.
	Wait()
}

// New creates and returns a new Executor object.
//
// workersCnt - specifies, how many workers(go-routines) simultaneously Executor will use to handle functions, sent for execution.
func New(workersCnt uint) Executor {
	if workersCnt == 0 {
		return newNoOpExecutor()
	}
	return newExecutor(workersCnt)
}
