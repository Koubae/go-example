// @doc https://pkg.go.dev/iter
// @doc https://bitfieldconsulting.com/posts/iterators
package main

import (
	"fmt"
	"iter"
)

func main() {
	// Simple
	for v := range iterItems {
		fmt.Println(v)
	}

	fmt.Println("---------------------------------------------------")

	for i, v := range GetItems() {
		fmt.Println(i, v)
	}

}

// Simple
func iterItems(yield func(int) bool) {
	items := [...]int{1, 2, 3, 4, 5}
	for _, v := range items {
		if !yield(v) {
			return
		}
	}

}

func GetItems() iter.Seq2[int, int] {
	return func(yield func(int, int) bool) {
		items := [...]int{1, 2, 3, 4, 5}
		for i, v := range items {
			if !yield(i, v) {
				return
			}
		}
	}
}
