package executor

import (
	"fmt"
	"sync"
)

type executor struct {
	workers chan struct{}
	wg      *sync.WaitGroup
}

func newExecutor(workersCnt uint) *executor {
	return &executor{
		workers: make(chan struct{}, workersCnt),
		wg:      &sync.WaitGroup{},
	}
}

func (exec *executor) Enqueue(fn func()) {
	if fn == nil {
		return
	}

	exec.workers <- struct{}{}
	exec.wg.Add(1)

	go func() {
		fn()

		<-exec.workers
		exec.wg.Done()
	}()
}

func (exec *executor) TryEnqueue(fn func()) error {
	if fn == nil {
		return nil
	}

	select {
	case exec.workers <- struct{}{}:
		exec.wg.Add(1)
	default:
		return fmt.Errorf("queue is full at the moment")
	}

	go func() {
		fn()

		<-exec.workers
		exec.wg.Done()
	}()

	return nil
}

func (exec *executor) Wait() {
	exec.wg.Wait()
}
