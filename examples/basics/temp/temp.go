package main

import (
	"fmt"
	"sync"
	"time"
)

func main() {
	var wg sync.WaitGroup

	wg.Add(1)

	go func() {
		defer wg.Done()
		fmt.Println("Go Routine A started -- (parent <main>)")

		go func() {
			fmt.Println("Go Routine A.B started -- (parent A)")

			go func() {
				fmt.Println("Go Routine A.B.C started -- (parent B)")

				go func() {
					fmt.Println("Go Routine A.B.C.D started -- (parent C)")

					time.Sleep(500 * time.Millisecond)
					fmt.Println("Go Routine A.B.C.D DONE -- (parent C)")
				}()

				time.Sleep(500 * time.Millisecond)
				fmt.Println("Go Routine A.B.C DONE -- (parent B)")
			}()

			time.Sleep(500 * time.Millisecond)
			fmt.Println("Go Routine A.B DONE -- (parent A)")
		}()

		// time.Sleep(3 * time.Second)
		time.Sleep(10 * time.Millisecond)
		fmt.Println("Go Routine A DONE -- (parent <main>)")

	}()

	wg.Wait()
	time.Sleep(3 * time.Second)
	fmt.Println("Main Thread Ended")
}
