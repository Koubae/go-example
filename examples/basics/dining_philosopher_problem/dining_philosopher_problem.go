package main

import (
	"fmt"
	"sync"
)

const TotalPhilosophers = 5
const EatSameTimeMax = 2
const TotalConsumptionPerPhilosopher = 3

type Chopstick struct {
	sync.Mutex
}

type Philosopher struct {
	Name           string
	ID             int
	LeftChopstick  *Chopstick
	RightChopstick *Chopstick
	EatRequestAck  chan struct{}
	EatTurnEnd     chan struct{}
	EatRequest     chan int
}

func (p *Philosopher) Eat(wg *sync.WaitGroup) {
	defer wg.Done()

	for _ = range TotalConsumptionPerPhilosopher {
		p.EatRequest <- p.ID // send eat request to host
		<-p.EatRequestAck    // wait host confirmation

		p.LeftChopstick.Lock()
		p.RightChopstick.Lock()

		fmt.Printf("starting to eat %s\n", p.Name)
		fmt.Printf("finishing eating %s\n", p.Name)

		p.RightChopstick.Unlock()
		p.LeftChopstick.Unlock()

		p.EatTurnEnd <- struct{}{} // notify host done eating

	}

}

type Host struct {
	Chopsticks   []*Chopstick
	Philosophers []*Philosopher
}

func (h *Host) Serve(done <-chan struct{}, eatTurnEnd chan struct{}, eatRequests chan int) {
	currentEating := 0
	queue := make([]int, 0)
	grant := func() {
		for currentEating < EatSameTimeMax && len(queue) > 0 {
			next := queue[0]
			queue = queue[1:]
			h.Philosophers[next].EatRequestAck <- struct{}{}
			currentEating++
		}
	}

	for {
		select {
		case <-done:
			return
		case id := <-eatRequests:
			queue = append(queue, id)
			grant()
		case <-eatTurnEnd:
			if currentEating > 0 {
				currentEating--
			}
			grant()
		}
	}

}

func main() {
	var wg sync.WaitGroup
	eatTurnEnd := make(chan struct{}, EatSameTimeMax)
	eatRequests := make(chan int, EatSameTimeMax)

	chopsticks := make([]*Chopstick, TotalPhilosophers)
	for i := 0; i < TotalPhilosophers; i++ {
		chopsticks[i] = new(Chopstick)
	}

	philosophers := make([]*Philosopher, TotalPhilosophers)
	for i := 0; i < TotalPhilosophers; i++ {
		name := fmt.Sprintf("%d", i+1)
		philosophers[i] = &Philosopher{
			Name:           name,
			ID:             i,
			EatRequestAck:  make(chan struct{}),
			EatTurnEnd:     eatTurnEnd,
			EatRequest:     eatRequests,
			LeftChopstick:  chopsticks[i],
			RightChopstick: chopsticks[(i+1)%TotalPhilosophers],
		}

	}

	for _, p := range philosophers {
		wg.Add(1)
		go p.Eat(&wg)
	}

	done := make(chan struct{})

	host := &Host{
		Chopsticks:   chopsticks,
		Philosophers: philosophers,
	}
	go host.Serve(done, eatTurnEnd, eatRequests)

	wg.Wait()
	close(done)

}
