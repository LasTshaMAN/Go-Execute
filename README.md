[![Build Status](https://travis-ci.org/LasTshaMAN/Go-Execute.svg?branch=master)](https://travis-ci.org/LasTshaMAN/Go-Execute)
[![Go Report Card](https://goreportcard.com/badge/github.com/LasTshaMAN/Go-Execute)](https://goreportcard.com/report/github.com/LasTshaMAN/Go-Execute)
[![codecov](https://codecov.io/gh/LasTshaMAN/Go-Execute/branch/master/graph/badge.svg)](https://codecov.io/gh/LasTshaMAN/Go-Execute)
[![GoDoc](https://godoc.org/github.com/LasTshaMAN/Go-Execute/executor?status.svg)](https://godoc.org/github.com/LasTshaMAN/Go-Execute/executor)

# Simple Executor for Golang to run your Jobs

This package implements ["Thread-Pool"](https://en.wikipedia.org/wiki/Thread_pool) design pattern for Golang.

Its main purpose is to decouple business logic from the logic necessary for go-routines management.

Go-Execute is designed to be simple and lightweight yet flexible enough to suit your needs. It is go-routine-safe (you can enqueue Jobs from different go-routines and expect Executor to work correctly).

## Example (basic enqueueing)

```Go
package main

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
```

## Example (non-blocking enqueueing)

```Go
package main

import (
	"fmt"
	"time"

	"github.com/LasTshaMAN/Go-Execute/executor"
)

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
```

## Example (returning value from enqueued function)

```Go
package main

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
```

## Example (waiting for all the jobs to finish)

```Go
package main

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
```

## Example (more elaborate ones)

For more real-world examples check out [examples](https://github.com/LasTshaMAN/Go-Execute/tree/master/examples) directory.

## Docs

https://godoc.org/github.com/LasTshaMAN/Go-Execute/executor

## Installation

```
go get github.com/LasTshaMAN/Go-Execute/executor
```

## Contributing

Feel free to submit issues, fork the repository and send pull requests!

Your suggestions on how to extend the functionality of Go-Execute to cover possible use-cases are also welcome!
