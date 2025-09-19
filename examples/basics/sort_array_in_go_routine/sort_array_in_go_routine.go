package main

import (
	"fmt"
	"sync"

	// "sync"
	"slices"
)

func main() {
	arrayInput := []int{2, 5, 23, 2, 6, 7, 8, 10, 11, 12}
	arrayL := len(arrayInput)

	partitionL := arrayL / 4
	reminder := arrayL % 4
	partitions := make([][]int, 0)
	index := 0
	for i := 0; i < 4; i++ {
		currentPartition := make([]int, 0)
		for j := 0; j < partitionL; j++ {
			value := arrayInput[index]
			currentPartition = append(currentPartition, value)
			index++
		}
		partitions = append(partitions, currentPartition)

	}
	for j := 0; j < reminder; j++ {
		value := arrayInput[index]
		partitions[j] = append(partitions[j], value)
		index++
	}

	var wg sync.WaitGroup
	c := make(chan []int, 4)

	wg.Add(4)

	for _, partition := range partitions {
		go sortPartitionArray(&wg, c, partition)
	}

	wg.Wait()

	results := make([]int, 0)
	for i := 0; i < 4; i++ {
		result := <-c
		results = append(results, result...)
	}
	slices.Sort(results)

	// fmt.Printf("Array Lenght %d, partitions size %d with reminder %d \n", arrayL, partitionL, reminder)
	// fmt.Println(partitions)
	fmt.Println(results)

}

func sortPartitionArray(wg *sync.WaitGroup, c chan<- []int, array []int) {
	defer wg.Done()

	slices.Sort(array)
	fmt.Println(array)

	c <- array

}
