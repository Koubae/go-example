package main

import (
	"fmt"
	"sync"
)

func main() {
	simple()

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
