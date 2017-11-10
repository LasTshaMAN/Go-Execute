package examples

import (
	"fmt"
	"github.com/LasTshaMAN/Go-Execute/jobs"
	"math/rand"
	"time"
)

func BasicEnqueueing() {
	executor := jobs.NewExecutor(4, 4)

	// Will block current go-routine if Executor is busy
	executor.Enqueue(func() {
		fmt.Println("World")
	})

	fmt.Println("Hello")
	time.Sleep(time.Second)
}

func NonBlockingEnqueueing() {
	rand.Seed(time.Now().UTC().UnixNano())
	executor := jobs.NewExecutor(4, 4)

	// Tasks keep coming ...
	for {
		// Will not block current go-routine
		err := executor.TryToEnqueue(func() {
			time.Sleep(time.Duration(rand.Intn(3)) * time.Second)
			fmt.Println("Some task has finished")
		})
		if err != nil {
			fmt.Println("Executor is full, can't enqueue more jobs at the moment ...")
			time.Sleep(1 * time.Second)
		}
	}
}

func GettingResultBack() {
	rand.Seed(time.Now().UTC().UnixNano())
	executor := jobs.NewExecutor(4, 4)

	out := make(chan int)
	executor.Enqueue(func() {
		fmt.Println("Some work is done here ...")
		out <- rand.Intn(10)
	})

	result := <-out
	fmt.Printf("result = %d", result)
}
