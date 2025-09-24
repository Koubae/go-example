package main

import (
	"fmt"
	"runtime"
	"sync"
	"time"
)

func main() {
	// simple()
	second()
}

func simple() {
	var wg sync.WaitGroup
	c := make(chan int)

	wg.Add(3)
	go product(&wg, c, 2, 5)
	go product(&wg, c, 3, 7)
	go product(&wg, c, 4, 2)

	a := <-c
	b := <-c
	d := <-c

	fmt.Println(a)
	fmt.Println(b)
	fmt.Println(d)
	wg.Wait()
}

func product(wg *sync.WaitGroup, c chan<- int, v1 int, v2 int) {
	defer wg.Done()
	product := v1 * v2
	c <- product
}

func second() {
	cpu := runtime.NumCPU()
	before := runtime.GOMAXPROCS(cpu)
	fmt.Printf("Setting GOMAXPROCS to current CPU count (%d), previous value: %d\n", cpu, before)

	var wg sync.WaitGroup

	size := 3
	cb := make(chan struct{}, size*100)
	quit := make(chan struct{}, size)

	// Add exactly 3 BEFORE we start the workers. if the buffer was 3 AND we had 1 more, we would deadlock here!
	cb <- struct{}{}
	cb <- struct{}{}
	cb <- struct{}{}

	for i := range size {
		wg.Add(1)
		go func(wg *sync.WaitGroup, i int, cb <-chan struct{}, quit <-chan struct{}) {
			defer wg.Done()
			defer fmt.Printf("Worker #%d terminated\n", i)
			fmt.Printf("Worker #%d started\n", i)

			work := func() {
				fmt.Printf("Worker #%d working...\n", i)
				time.Sleep(250 * time.Millisecond)
			}

		workLoop:
			for {
				select {
				case <-quit:
					fmt.Printf("Worker #%d go quit signal, checking if channel has still values\n", i)

					for {
						select {
						case _, ok := <-cb:
							if !ok {
								fmt.Printf("Worker #%d channel has no more work to do, stopping....\n", i)
								break workLoop
							}
							fmt.Printf("Worker #%d channel has still values, working...\n", i)
							work()
						default:
							fmt.Printf("Worker #%d channel has no item in queue, stopping....\n", i)
							break workLoop
						}
					}

				case <-cb:
					work()
				}
			}

		}(&wg, i, cb, quit)
	}

	// Add exactly 3
	cb <- struct{}{}
	cb <- struct{}{}
	cb <- struct{}{}

	// Add few more...
	cb <- struct{}{}
	cb <- struct{}{}
	cb <- struct{}{}
	cb <- struct{}{}

	// let's fill up the channel
	for range 100 {
		cb <- struct{}{}
	}
	close(cb) // close the channel so that we know no more work will be added

	// time.Sleep(time.Second)
	for range size {
		quit <- struct{}{}
	}

	fmt.Println("Closing request channel...")
	time.Sleep(250 * time.Millisecond)
	fmt.Println("Closed request channel!")

	wg.Wait()
	fmt.Println("-- Second Terminated")

}
