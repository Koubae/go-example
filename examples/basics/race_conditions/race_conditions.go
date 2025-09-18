/*
Why does race condition occur?

Go routine 1 will increase the value of x before Go Routine 2 can read it.
I intentionally make go routine 1 sleep for 500ms to make it more obvious.

This is what happens in more details:

Go Routine 1:

1. Go Routine 1 starts
2. Starts a loop in range 10 times
3. Each loop increases x by 1
4. Sleeps for 500ms

Go Routine 2:

1. Go Routine 2 starts
2. Starts a loop in range 10 times
3. Each loop prints the value of x
4. Sleeps for 1 second

The sleep difference is to make the race condition more obvious but since there is no lock when increasing X
from 1 and no waiting in go routine 2, go routine 2 could skip printing some increased values

*/

package main

import (
	"fmt"
	"sync"
	"time"
)

func main() {
	fmt.Println("Main Thread Started")

	x := 1

	var wg sync.WaitGroup

	wg.Add(2)

	// Start 1st Go Routine
	go func() {
		defer wg.Done()

		fmt.Println("Go Routine 1 Started")
		for _ = range 10 {
			x++
			time.Sleep(500 * time.Millisecond)
		}
		fmt.Println("Go Routine 1 Ended")
	}()

	// Start the 2nd Go Routine
	go func() {
		defer wg.Done()

		fmt.Println("Go Routine 2 Started")

		for _ = range 10 {
			fmt.Printf("X=%d\n", x)
			time.Sleep(1 * time.Second)
		}

		fmt.Println("Go Routine 2 Ended")

	}()

	wg.Wait()
	fmt.Println("Main Thread Ended")

}
