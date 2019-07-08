package examples

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/LasTshaMAN/Go-Execute/executor"
)

func gettingTheResultBack() {
	rand.Seed(time.Now().UTC().UnixNano())
	exec := executor.New(4)

	out := make(chan int)
	exec.Enqueue(func() {
		fmt.Println("Some work is done here ...")
		out <- rand.Intn(10)
	})

	result := <-out
	fmt.Printf("result = %d", result)
}
