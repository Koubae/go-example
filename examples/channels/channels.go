package main

import (
	"fmt"
)

func main() {
	iterateOverChannels()
}

/*
https://gobyexample.com/range-over-channels
*/
func iterateOverChannels() {
	queue := make(chan string, 4)

	queue <- "one"
	queue <- "two"
	queue <- "three"
	queue <- "four"
	close(queue) // To iterate, the channel should be closed or will throw a deadlock error.

	for item := range queue {
		fmt.Printf("received %s\n", item)
	}

}
