package main

import (
	"fmt"
	"time"
)

func sum(s []int, channel chan int) {
	sum := 0
	for _, v := range s {
		sum += v
	}
	channel <- sum
}

func main() {
	integers := []int{7, 2, 8, -9, 4, 0}

	channel := make(chan int)

	go sum(integers[:len(integers)/2], channel)
	go sum(integers[len(integers)/2:], channel)

	var x int
	var y int

	x = <-channel
	y = <-channel

	fmt.Printf("X=%d\nY=%d\nX+Y=%d\n", x, y, x+y)

	// Anonymous
	go func() {
		fmt.Println("Go routine 1")
		channel <- 1
	}()
	fmt.Printf("Go routine anonymous value %v", <-channel)

	// Buffered channel

	channel_buffer := make(chan int, 2)

	channel_buffer <- 1
	channel_buffer <- 2

	fmt.Println(<-channel_buffer)
	fmt.Println(<-channel_buffer)

	// Range & Close
	const FibonacciNumber = 10
	fibonacci_channel := make(chan int, FibonacciNumber)

	go fibonacci(FibonacciNumber, fibonacci_channel)
	for i := range fibonacci_channel {
		fmt.Printf("- fibonacci: %d\n", i)
	}

	// anonymous go routine
	fmt.Printf("\n\nanonymous go routine\n")
	fibonacci_channel_anonymous := make(chan int, FibonacciNumber)
	go func(n int, channel chan int) {
		var x = 0
		var y = 1

		for i := 0; i < n; i++ {
			channel <- x
			x, y = y, x+y
		}
		close(channel) // notify to receiver that channel is closed
	}(FibonacciNumber, fibonacci_channel_anonymous)

	for i := range fibonacci_channel_anonymous {
		fmt.Printf("- fibonacci: %d\n", i)
	}

	// Select
	fmt.Println("\nShowing select")
	channel1 := make(chan int)
	channelQuit := make(chan int)
	go func() {
		for i := 0; i < FibonacciNumber; i++ {
			fmt.Printf("NEXT_NUMBER: %d\n", <-channel1)
		}
		channelQuit <- 0
	}()

	func(channel chan int, quit chan int) {
		x, y := 0, 1
		var loop int = 0
		for {
			fmt.Printf("for select loop %d\n", loop)
			loop++
			select {
			case channel1 <- x:
				x, y = y, x+y
			case <-quit:
				fmt.Println("Channel is quitting")
				return
			}
		}
	}(channel1, channelQuit)

	// Default Select
	bomb()
}

func fibonacci(n int, channel chan int) {
	var x = 0
	var y = 1

	for i := 0; i < n; i++ {
		channel <- x
		x, y = y, x+y
	}
	close(channel) // notify to receiver that channel is closed
}

func bomb() {
	tick := time.Tick(100 * time.Millisecond)
	boom := time.After(500 * time.Millisecond)
	for {
		select {
		case <-tick:
			fmt.Println("tick.")
		case <-boom:
			fmt.Println("BOOM!")
			return
		default:
			fmt.Println("    .")
			time.Sleep(50 * time.Millisecond)
		}
	}
}
