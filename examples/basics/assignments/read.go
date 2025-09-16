package assignments

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

type Person struct {
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
}

func Read() {
	fmt.Printf(">>> File name: ")
	var fileName string // examples/basics/names.txt
	fmt.Scan(&fileName)

	fileDescriptor, err := os.Open(fileName)
	if err != nil {
		panic(err)
	}
	defer fileDescriptor.Close()

	fmt.Println()

	names := make([]Person, 0)
	scanner := bufio.NewScanner(fileDescriptor)
	for scanner.Scan() {

		line := scanner.Text()

		fullName := strings.Fields(line)

		firstName := fullName[0]
		lastName := fullName[1]

		p := Person{
			FirstName: firstName,
			LastName:  lastName,
		}
		names = append(names, p)
	}
	if err := scanner.Err(); err != nil {
		panic(err)
	}

	for _, p := range names {
		fmt.Printf("%s %s\n", p.FirstName, p.LastName)
	}
}
