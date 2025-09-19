package main

import (
	"context"
	"fmt"
	"time"
)

func main() {
	// exampleContextWithCancelAndIterateOverChannel()
	contextWithValue()
}

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
