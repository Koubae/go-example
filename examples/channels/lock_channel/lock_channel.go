package main

import (
	"fmt"
	"sync"
	"time"
)

type ChanLock chan struct{}

func NewChanLock() ChanLock {
	l := make(chan struct{}, 1)
	l <- struct{}{}
	return l
}

func (l ChanLock) Lock() {
	<-l
}

func (l ChanLock) Unlock() {
	l <- struct{}{}
}

func main() {
	lock := NewChanLock()

	wg := &sync.WaitGroup{}
	wg.Add(3)

	for i := range 3 {
		go func(id int) {
			defer wg.Done()
			fmt.Printf("Worker %d: Starting... Acquiring Lock\n", id+1)

			lock.Lock()
			defer lock.Unlock()
			time.Sleep(500 * time.Millisecond) // Lock a bit to show that other workers will block here

			fmt.Printf("Worker %d: Working....\n", id+1)
			time.Sleep(1 * time.Second)
			fmt.Printf("Worker %d: Done!\n", id+1)
		}(i)
	}

	wg.Wait()
	fmt.Println("All jobs are done")

}
