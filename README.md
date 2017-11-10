[![Build Status](https://travis-ci.org/LasTshaMAN/Go-Execute.svg?branch=master)](https://travis-ci.org/LasTshaMAN/Go-Execute)
[![Go Report Card](https://goreportcard.com/badge/github.com/LasTshaMAN/Go-Execute)](https://goreportcard.com/report/github.com/LasTshaMAN/Go-Execute)
[![codecov](https://codecov.io/gh/LasTshaMAN/Go-Execute/branch/master/graph/badge.svg)](https://codecov.io/gh/LasTshaMAN/Go-Execute)
[![GoDoc](https://godoc.org/github.com/LasTshaMAN/Go-Execute/jobs?status.svg)](https://godoc.org/github.com/LasTshaMAN/Go-Execute/jobs)

# Simple Executor for Golang to run your Jobs

This package implements ["Thread-Pool"](https://en.wikipedia.org/wiki/Thread_pool) design pattern for Golang.

Its main purpose is to decouple business logic from the logic necessary for go-routines management.

Go-Execute is designed to be simple and lightweight yet flexible enough to suit your needs.

## Example (basic enqueueing)

```Go
package main

import (
	"github.com/LasTshaMAN/Go-Execute/jobs"
	"fmt"
	"time"
)

func main()  {
	executor := jobs.NewExecutor(4, 4)

	// Will block current go-routine if Executor is busy
	executor.Enqueue(func() {
		fmt.Println("World")
	})

	fmt.Println("Hello")
	time.Sleep(time.Second)
}
```

## Example (non-blocking enqueueing)

```Go
package main

import (
	"github.com/LasTshaMAN/Go-Execute/jobs"
	"fmt"
	"time"
)

func main()  {
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
```

## Example (returning value from enqueued function)

```Go
package main

import (
	"github.com/LasTshaMAN/Go-Execute/jobs"
	"fmt"
	"time"
)

func main()  {
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
```

## Example (more elaborate ones)

For more real-world examples check out [examples](https://godoc.org/github.com/LasTshaMAN/Go-Execute/examples) directory.

## Docs

https://godoc.org/github.com/LasTshaMAN/Go-Execute/jobs

## Installation

```
go get github.com/LasTshaMAN/Go-Execute/jobs
```

## Contributing

Feel free to submit issues, fork the repository and send pull requests!

Your suggestions on how to extend the functionality of Go-Execute to cover possible use-cases are also welcome!