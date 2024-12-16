package main

import (
	"fmt"
	"sync"
)

type Counter struct {
	mutex sync.Mutex
	value map[string]int
}

func (c *Counter) Inc(key string) {
	c.mutex.Lock()
	c.value[key]++
	c.mutex.Unlock()
}

func (c *Counter) Value(key string) int {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	return c.value[key]
}

// @doc: https://go.dev/tour/concurrency/9
func main() {
	counter := Counter{value: make(map[string]int)}

	const counterKey = "some-keys"
	for i := 0; i < 1000; i++ {
		go counter.Inc(counterKey)
	}

	// wait for go routines to terminate
	//time.Sleep(time.Second)
	// Just to show the progression, you would never ever do something like this unless you are on something heavy!
	for value := counter.Value(counterKey); value != 1000; value = counter.Value(counterKey) {
		fmt.Println(counter.Value(counterKey))
	}
	fmt.Println(counter.Value(counterKey))
}
