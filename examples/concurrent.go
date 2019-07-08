package examples

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/LasTshaMAN/Go-Execute/executor"
)

func enqueueingFromMultipleGoRoutines() {
	rand.Seed(time.Now().UTC().UnixNano())
	exec := executor.New(4)

	out := make(chan int)
	for i := 0; i < 16; i++ {
		// Different go-routines use Executor to enqueue Jobs
		go func() {
			exec.Enqueue(func() {
				workAmount := rand.Intn(10)
				fmt.Printf("Finished processing %d\n", workAmount)
				out <- workAmount
			})
		}()
	}

	// Getting the results of executed Jobs
	for i := 0; i < 16; i++ {
		result := <-out
		fmt.Printf("result = %d\n", result)
	}
}
