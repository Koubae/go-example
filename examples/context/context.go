package main

import (
	"context"
	"fmt"
	"time"
)

func main() {
	// exampleContextWithCancelAndIterateOverChannel()
	// contextWithValue()
	// manualCancel()
	timeout()
}

func dummyWorker1(ctx context.Context) {
	name := "dummyWorker1"
	fmt.Printf("%s started\n", name)
	for {
		select {
		case <-ctx.Done():
			fmt.Printf("%s topped, ctx error: %v\n", name, ctx.Err())
			return
		default:
			fmt.Printf("%s working...\n", name)
			time.Sleep(500 * time.Millisecond)
		}
	}
}

func dummyWorker2WithTimeout(ctx context.Context, timeout time.Duration) {
	name := "dummyWorker2WithTimeout"
	fmt.Printf("%s started\n", name)

	select {
	case <-time.After(timeout):
		fmt.Printf("%s Task has finished its work\n", name)
	case <-ctx.Done():
		fmt.Printf("%s Conext Timeouted out error: %v\n", name, ctx.Err())
	}
}

/*
	================================================================================
			Simple Examples
	================================================================================

*/

// -------------------------------------------------
// Cancel Manually

func manualCancel() {
	ctx, cancel := context.WithCancel(context.Background())

	go dummyWorker1(ctx)

	time.Sleep(2 * time.Second)
	cancel() // stop worker
	time.Sleep(1 * time.Second)
	fmt.Println("Done")
}

// -------------------------------------------------
// Timeout
func timeout() {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)

	// INCREASE timeoutCaller and DECREASE timeoutWorker: worker will complete its work before timeout
	// DECREASE timeoutCaller and INCREASE timeoutWorker: worker will time out
	timeoutCaller := 1 * time.Second
	timeoutWorker := 3 * time.Second

	go func() {
		defer cancel()
		go dummyWorker2WithTimeout(ctx, timeoutWorker)
		time.Sleep(timeoutCaller)
	}()

	time.Sleep(3 * time.Second)
	fmt.Println("Done")

}

// -------------------------------------------------
// WithValue
func contextWithValue() {
	type customKeyT string

	getValueFromContext := func(ctx context.Context, k customKeyT) {
		if v := ctx.Value(k); v != nil {
			fmt.Printf("Value found with key '%s' value=%v, type='%T'\n", k, v, v)
		} else {
			fmt.Printf("Value NOT found with key '%s'\n", k)
		}
	}

	k := customKeyT("special-secret-key")
	ctx := context.WithValue(context.Background(), k, "secret-value")

	getValueFromContext(ctx, k)
	getValueFromContext(ctx, customKeyT("another-key"))

	obj := struct {
		name  string
		value int
	}{
		name:  "ObjectName",
		value: 99999,
	}

	k2 := customKeyT("special-secret-key2")
	ctx2 := context.WithValue(context.TODO(), k2, obj)

	getValueFromContext(ctx2, k2)
	getValueFromContext(ctx2, customKeyT("another-key2"))

}

/*
https://pkg.go.dev/context#pkg-variables
*/
func exampleContextWithCancelAndIterateOverChannel() {
	geneFunc := func(ctx context.Context) <-chan int {
		destination := make(chan int)
		n := 1

		go func() {
			for {
				select {
				case <-ctx.Done():
					fmt.Println("geneFunc Done")
					return
				case destination <- n:
					n++
				}
			}
		}()

		return destination
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	for n := range geneFunc(ctx) {
		fmt.Println(n)
		time.Sleep(500 * time.Millisecond)
		if n == 3 {
			break
		}
	}
	// time.Sleep(500 * time.Millisecond)
	fmt.Println("Done")
}
