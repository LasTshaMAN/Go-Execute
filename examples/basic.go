package examples

import (
	"fmt"
	"time"

	"github.com/LasTshaMAN/Go-Execute/executor"
)

func basicBlockingEnqueueing() {
	exec := executor.New(4)

	// Will block current go-routine if Executor is busy
	exec.Enqueue(func() {
		fmt.Println("World")
	})

	fmt.Println("Hello")
	time.Sleep(time.Millisecond)
}

func basicNonBlockingEnqueueing() {
	exec := executor.New(4)

	// Will block current go-routine if Executor is busy
	err := exec.TryEnqueue(func() {
		fmt.Println("World")
	})
	if err != nil {
		fmt.Println("Executor is full, can't enqueue more jobs at the moment ...")
		time.Sleep(1 * time.Millisecond)
	}

	fmt.Println("Hello")
	time.Sleep(time.Millisecond)
}
