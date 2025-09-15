/*
Write a program which prompts the user to enter integers and stores the integers in a sorted slice.

The program should be written as a loop.
Before entering the loop, the program should create an empty integer slice of size (length) 3.

During each pass through the loop, the program prompts the user to enter an integer to be added to the slice.
The program adds the integer to the slice, sorts the slice, and prints the contents of the slice in sorted order.

The slice must grow in size to accommodate any number of integers which the user decides to enter.
The program should only quit (exiting the loop) when the user enters the character ‘X’ instead of an integer.
*/
package main

import (
	"fmt"
	"sort"
	"strconv"
)

const ExitInput = 'X'

func main() {
	var input string

	intergers := make([]int, 0)
	for {
		fmt.Scan(&input)

		if input == string(ExitInput) {
			break
		}

		n, err := strconv.Atoi(input)
		if err != nil {
			panic("Invalid input, error: " + err.Error())
		}

		intergers = append(intergers, n)
		sort.Ints(intergers)
		fmt.Println(intergers)

	}

}
