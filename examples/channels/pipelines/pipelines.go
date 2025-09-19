/*
https://go.dev/blog/pipelines
*/

package pipelines

import (
	"fmt"
	"sync"
)

func PipelineOneSimple() {
	fmt.Println("================================")
	fmt.Println("     PipelineOneSimple    ")
	fmt.Println("================================")
	numChannel := gen(2, 3, 5, 7, 11, 13)
	pipelineSquare := squareN(numChannel)

	for n := range pipelineSquare {
		fmt.Println(n)
	}

	for n := range squareN(squareN(squareN(gen(1, 2, 3)))) {
		fmt.Printf("Triple squared: %d\n", n)
	}

}

func PipelineTwoFanOutFanIn() {
	fmt.Println("================================")
	fmt.Println("     PipelineTwoFanOutFanIn    ")
	fmt.Println("================================")
	in := gen(1, 2, 3, 4, 5, 6, 7, 8, 9, 10)

	c1 := squareN(in)
	c2 := squareN(in)
	c3 := squareN(in)
	c4 := squareN(in)

	for n := range merge(c1, c2, c3, c4) {
		fmt.Println(n)
	}
}

func PipelineThreeExplicitCancel() {
	fmt.Println("================================")
	fmt.Println("     PipelineThreeExplicitCancel    ")
	fmt.Println("================================")

	done := make(chan struct{})
	defer close(done)

	// in := genBufferedWithCancel(done, 2, 3, 4, 5, 6, 7, 8, 9, 10)
	in := genWithCancel(done, 2, 3, 4, 5, 6, 7, 8, 9, 10)

	c1 := squareNWithCancel(done, in)
	c2 := squareNWithCancel(done, in)
	c3 := squareNWithCancel(done, in)

	out := mergeWithCancel(done, c1, c2, c3)
	fmt.Println(<-out) // consume 1st Value
	fmt.Println(<-out) // consume 2nd Value
	fmt.Println(<-out) // consume 3rd Value

	// done will be closed by the deferred call.

}

func gen(nums ...int) <-chan int {
	out := make(chan int)

	go func() {
		for _, n := range nums {
			out <- n
		}
		close(out)
	}()

	return out
}

func genWithCancel(done <-chan struct{}, nums ...int) <-chan int {
	out := make(chan int)

	go func() {
		defer close(out)
		for _, n := range nums {
			select {
			case out <- n:
			case <-done:
				return
			}
			out <- n
		}
	}()

	return out
}

func genBuffered(nums ...int) <-chan int {
	out := make(chan int, len(nums))

	for _, n := range nums {
		out <- n
	}
	close(out)
	return out
}

func genBufferedWithCancel(done <-chan struct{}, nums ...int) <-chan int {
	out := make(chan int, len(nums))
	defer close(out)

	for _, n := range nums {
		select {
		case out <- n:
		case <-done:
			return out
		}
	}
	return out
}

func squareN(in <-chan int) <-chan int {
	out := make(chan int)

	go func() {
		for n := range in {
			out <- n * n
		}
		close(out)
	}()

	return out
}

func squareNWithCancel(done <-chan struct{}, in <-chan int) <-chan int {
	out := make(chan int)

	go func() {
		defer close(out)
		for n := range in {
			select {
			case out <- n * n:
			case <-done:
				return
			}
		}
	}()

	return out
}

// Fan-Out | Fan-In => https://go.dev/blog/pipelines#fan-out-fan-in
func merge(outboundChannels ...<-chan int) <-chan int {
	var wg sync.WaitGroup
	out := make(chan int)

	pipeline := func(outbound <-chan int) {
		for n := range outbound {
			out <- n
		}
		wg.Done()
	}
	wg.Add(len(outboundChannels))
	for _, outbound := range outboundChannels {
		go pipeline(outbound)
	}

	go func() {
		wg.Wait()
		close(out)
	}()

	return out

}

// Fan-Out | Fan-In => https://go.dev/blog/pipelines#fan-out-fan-in
func mergeWithCancel(done <-chan struct{}, outboundChannels ...<-chan int) <-chan int {
	var wg sync.WaitGroup
	out := make(chan int)

	pipeline := func(outbound <-chan int) {
		defer wg.Done()
		for n := range outbound {
			select {
			case out <- n:
			case <-done:
				return
			}
		}
	}
	wg.Add(len(outboundChannels))
	for _, outbound := range outboundChannels {
		go pipeline(outbound)
	}

	go func() {
		wg.Wait()
		close(out)
	}()

	return out

}
