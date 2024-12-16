// @credit: https://gobyexample.com/range-over-iterators
package main

import (
	"fmt"
	"iter"
)

func fibonacciIterator() iter.Seq[int] {
	return func(yield func(int) bool) {
		a, b := 1, 1

		for {
			if !yield(a) {
				return
			}
			a, b = b, a+b
		}
	}
}

func main() {
	for n := range fibonacciIterator() {
		if n >= 10 {
			break
		}
		fmt.Println(n)
	}
}
