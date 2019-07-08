package executor

type noOpExecutor struct {
}

func newNoOpExecutor() *noOpExecutor {
	return &noOpExecutor{}
}

func (exec *noOpExecutor) Enqueue(fn func()) {
}

func (exec *noOpExecutor) TryEnqueue(fn func()) error {
	return nil
}

func (exec *noOpExecutor) Wait() {
}
