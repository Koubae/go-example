/*
Problem Description:
Implement a Worker Queue with Synchronous Execution and Order Guarantee

You are tasked with implementing a worker queue that processes tasks using a pool of worker goroutines.
The queue should allow synchronous task execution while ensuring the following requirements:
1. Synchronous Execution: The caller of the queue must wait for the result of their task before proceeding. This means tasks are executed in a blocking manner for the caller.
2. Order Guarantee: Tasks submitted to the queue must be processed in the order they are received, even if workers complete their processing out of order.
3. Concurrent Processing: The queue should use a fixed number of workers to process tasks concurrently, but the caller should still receive results in the correct order.

The queue should support the following operations:
* Synchronous Task Submission (SyncExec):
- Accepts a task (request) and returns the result (response) after the task is processed.
- Blocks the caller until the task is processed and the result is available.
- If the queue has been stopped, calling this method should result in a panic.

* Queue Shutdown (Finish):
- Gracefully stops all workers and marks the queue as finished.
- After the queue is finished, any further calls to SyncExec should panic.
* Queue Initialization (New):
- Creates a new queue with a fixed number of workers and a transformation function that defines how tasks (requests) are processed to generate results (responses).

Example Scenario:
A developer is building a logging system that processes log messages in parallel but ensures they are written to a database in the order they are received.
Each log message is transformed into a formatted string before storage.
Using your implementation, the developer can submit log messages to the worker queue and retrieve the processed result in the correct order.

Constraints:
  - Tasks should be represented as a generic Request type, and results as a generic Response type
    (In the first approach, we remove this requirement to make it simpler and focus only on concurrency).
  - The number of workers should be configurable during initialization.
  - The solution should avoid starting unnecessary goroutines during synchronous execution.
*/
package main

import (
	"fmt"
	"math/rand"
	"sync"
	"time"
)

type workRequest struct {
	request  int
	response chan string
}

// Queue represents an object that can execute synchronously (but orderly) a given transformation
// function (function that gets a Request and returns a response). Note that on this version the
// actual transformation is defined when creating the object.
type Queue struct {
	workQueue chan workRequest
	// once      sync.Once // We could use this to ensure Finish only executes once
}

// SyncExec will execute a request synchronously. It will block the calling routine until the
// request is processed
func (w *Queue) SyncExec(req int) string {
	workRequest := workRequest{
		request:  req,
		response: make(chan string),
	}
	w.workQueue <- workRequest
	return <-workRequest.response
}

// Finish kills all the workers. After calling Finish, calling SyncExec will panic
func (w *Queue) Finish() {
	close(w.workQueue)
	/*
	   // We could have used w.once.Do to ensure it only executes once, but it’s not a requirement
	   w.once.Do(func() {
	      close(w.workChannel)
	   })
	*/
}

// New creates a new Queue object with a given transformation function.
func New(size int, transformation func(req int) string) *Queue {
	rv := &Queue{
		workQueue: make(chan workRequest),
	}

	for i := 0; i < size; i++ {
		go func() {
			for req := range rv.workQueue {
				// res := transformation(req.request)
				// fmt.Printf("Worker %d Processing request: %v Result: %s\n", i, req.request, res)
				req.response <- transformation(req.request)
			}
		}()
	}
	return rv
}

func main() {
	queue := New(
		3, func(req int) string {
			time.Sleep(time.Duration(rand.Intn(100)) * time.Millisecond) // We introduce a random delay so it’s clear that each consumer receives its own response
			return fmt.Sprintf("%d", req*2)
		},
	)

	wg := &sync.WaitGroup{}
	numRequests := 10
	wg.Add(numRequests)
	for i := 0; i < numRequests; i++ {
		go func(i int) {
			result := queue.SyncExec(i)
			fmt.Printf("Request: %d, Result: %s\n", i, result)
			wg.Done()
		}(i)
	}
	wg.Wait()
	queue.Finish()
}
