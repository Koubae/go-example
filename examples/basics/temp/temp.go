package main

import (
	"fmt"
	"math/rand"
	"sync"
	"time"
)

type chopStick struct {
	sync.Mutex
}

type philo struct {
	leftCS, rightCS *chopStick
	id              int
	eats            int
	finish          chan int
	start           chan int
	requestEat      chan int
}

func (p philo) eat() {
	for {
		// cant eat more then 3 times
		if p.eats == 3 {
			p.finish <- p.id
			gw.Done()
			return
		}
		// send start eating request to host
		p.requestEat <- p.id
		// block on startEat chanel if a value is (1) start eating else wait in a line.
		if eat := <-p.start; eat == 1 {
			p.leftCS.Lock()
			p.rightCS.Lock()
			fmt.Printf("[#] starting to eat %d\n", p.id)
			// random time to simulate eating time.
			time.Sleep(time.Millisecond * time.Duration(rand.Intn(p.id*1000)))
			fmt.Printf("[#] finishing eating %d\n", p.id)
			p.eats++
			p.finish <- p.id
			p.leftCS.Unlock()
			p.rightCS.Unlock()
		}
	}
}

var gw sync.WaitGroup

func main() {
	// init section
	gw.Add(5)
	chopSticks := make([]*chopStick, 5)
	philos := make([]*philo, 5)
	finishEat := make(chan int)
	startEat := make(chan int)
	requestEat := make(chan int)
	quitConf := make(chan int)
	quit := make(chan int)
	for i := 0; i < 5; i++ {
		chopSticks[i] = new(chopStick)
	}
	for i := 0; i < 5; i++ {
		philos[i] = &philo{
			leftCS:     chopSticks[i],
			rightCS:    chopSticks[(i+1)%5],
			id:         i + 1,
			eats:       0,
			finish:     finishEat,
			start:      startEat,
			requestEat: requestEat,
		}
	}

	// generate goroutines
	go host(requestEat, startEat, finishEat, quit, quitConf)
	for _, v := range philos {
		go v.eat()
	}
	// wait till philos to finish eating.
	gw.Wait()
	// send quit sign to host
	quit <- 1
	<-quitConf
}

func host(reqEat, startEat, finishEat, quit, quitConf chan int) {
	clientsID := make(map[int]string)
	defer func() {
		quitConf <- 1
	}()
	for {
		select {
		case id := <-reqEat:
			if len(clientsID) < 2 {
				clientsID[id] = "Eating"
				startEat <- 1
			} else {
				startEat <- 0
			}
		case id := <-finishEat:
			delete(clientsID, id)
		case <-quit:
			return
		}
	}

}
