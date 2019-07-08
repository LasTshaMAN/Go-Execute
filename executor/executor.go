package executor

import (
	"context"
	"fmt"

	"golang.org/x/sync/semaphore"
)

type executor struct {
	workersCnt uint
	sem        *semaphore.Weighted
}

func newExecutor(workersCnt uint) *executor {
	return &executor{
		workersCnt: workersCnt,
		sem:        semaphore.NewWeighted(int64(workersCnt)),
	}
}

func (exec *executor) Enqueue(fn func()) {
	if fn == nil {
		return
	}

	_ = exec.sem.Acquire(context.Background(), 1)
	go func() {
		fn()
		exec.sem.Release(1)
	}()
}

func (exec *executor) TryEnqueue(fn func()) error {
	if fn == nil {
		return nil
	}

	success := exec.sem.TryAcquire(1)
	if !success {
		return fmt.Errorf("queue is full at the moment")
	}
	go func() {
		fn()
		exec.sem.Release(1)
	}()

	return nil
}

func (exec *executor) Wait() {
	_ = exec.sem.Acquire(context.Background(), int64(exec.workersCnt))
	exec.sem.Release(int64(exec.workersCnt))
}
