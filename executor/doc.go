// Package executor implements "Thread-Pool" design pattern - https://en.wikipedia.org/wiki/Thread_pool.
//
// Its main purpose is to decouple business logic from the logic necessary for go-routines management.
//
// This package is designed to be go-routine-safe for usage in multi-go-routine environment.
// You can enqueue functions for execution from different go-routines and expect Executor to work correctly.
package executor
