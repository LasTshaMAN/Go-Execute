package examples

import (
	"fmt"

	"github.com/LasTshaMAN/Go-Execute/executor"
)

func enqueueAndWait() {
	exec := executor.New(4)

	for _, jobID := range []int64{1, 2, 3} {
		exec.Enqueue(func() {
			fmt.Printf("job: %d", jobID)
		})
	}
	exec.Wait()

	fmt.Println("All the jobs are done")
}
