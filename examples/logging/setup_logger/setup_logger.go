package main

import (
	"fmt"
	"log"
	"os"
	"sync"
	"time"
)

func main() {
	// Main uses default logger
	log.SetFlags(0)
	log.SetPrefix("main: ")
	log.Println("Starting ...")

	var wg sync.WaitGroup

	wg.Add(3)

	for i := range 3 {
		go func() {
			defer wg.Done()

			name := fmt.Sprintf("queue_%d ", i)
			// create a logger for each worker with a custom prefix
			logger := log.New(os.Stderr, name, 0)
			logger.Printf("(worker %d) Starting ...\n", i)

			for range 3 {
				logger.Printf("(worker %d) Working ...\n", i)
				time.Sleep(500 * time.Millisecond)
			}

		}()
	}

	log.Println("Spawned 3 workers...")
	log.Println("Waiting for all workers to finish...")
	wg.Wait()

	log.Println("All workers finished!")

}
