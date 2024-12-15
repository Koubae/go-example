package main

import "fmt"

type Number interface {
	int64 | float64
}

func main() {
	fmt.Println("Hello World 2")

	// Initialize a map of integers
	ints := map[string]int64{
		"first":  34,
		"second": 12,
	}

	floats := map[string]float64{
		"first":  34.5,
		"second": 26.99,
	}

	fmt.Printf("Non-Generic Sums: %v and %v\n",
		SumIntegers(ints),
		SumFloats(floats))

	fmt.Printf("Generic Sums: %v and %v\n",
		SumNumbers[string, int64](ints),
		SumNumbers[string, float64](floats))

	fmt.Printf("Generic Sums, type parameters inferred: %v and %v\n",
		SumNumbers(ints),
		SumNumbers(floats))

	fmt.Printf("Generic Sums with Constraint: %v and %v\n",
		SumNumbers(ints),
		SumNumbers(floats))

}

func SumIntegers(m map[string]int64) int64 {
	var sum int64
	for _, v := range m {
		sum += v
	}
	return sum
}

func SumFloats(m map[string]float64) float64 {
	var sum float64
	for _, v := range m {
		sum += v
	}
	return sum
}

// Here a function that uses interface
func SumNumbers[K comparable, V Number](m map[K]V) V {
	var sum V
	for _, v := range m {
		sum += v
	}
	return sum
}
