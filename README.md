[![Build Status](https://travis-ci.org/LasTshaMAN/Go-Execute.svg?branch=master)](https://travis-ci.org/LasTshaMAN/Go-Execute)
[![Go Report Card](https://goreportcard.com/badge/github.com/LasTshaMAN/Go-Execute)](https://goreportcard.com/report/github.com/LasTshaMAN/Go-Execute)
[![codecov](https://codecov.io/gh/LasTshaMAN/Go-Execute/branch/master/graph/badge.svg)](https://codecov.io/gh/LasTshaMAN/Go-Execute)
[![GoDoc](https://godoc.org/github.com/LasTshaMAN/Go-Execute/jobs?status.svg)](https://godoc.org/github.com/LasTshaMAN/Go-Execute/jobs)

# Simple Executor for Golang to run your Jobs

This package implements ["Thread-Pool"](https://en.wikipedia.org/wiki/Thread_pool) design pattern for Golang.

Its main purpose is to decouple business logic from the logic necessary for go-routines management.

## Example (basic)

```Go
package main

import (
	"github.com/LasTshaMAN/Go-Execute/jobs"
	"fmt"
	"time"
)

func main()  {
	executor := jobs.NewExecutor(4, 4)

	executor.Enqueue(func() {
		fmt.Println("Will be executed at some point in the future asynchronously")
	})

	fmt.Println("This statement will probably be reached first")
	time.Sleep(1 * time.Second)
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

	for {
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

## Docs
https://godoc.org/github.com/LasTshaMAN/Go-Execute/jobs

## Installation
```
go get github.com/LasTshaMAN/Go-Execute/jobs
```