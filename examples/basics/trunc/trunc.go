package main

import (
	"fmt"
	"strconv"
	"strings"
)

func main() {
	var input float64
	var inputString string
	var inputInt int

	fmt.Scanln(&input)
	inputString = fmt.Sprintf("%v", input)
	inputString = strings.Split(inputString, ".")[0]

	inputInt, _ = strconv.Atoi(inputString)
	fmt.Println(inputInt)
}
