package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
)

func main() {
	const inputFileName = "input.txt"
	const outputFileName = "output.txt"

	input, err := os.Open(inputFileName)
	if err != nil {
		panic(err)
	}
	output, err := os.Create(outputFileName)
	if err != nil {
		panic(err)
	}

	defer input.Close()
	defer output.Close()

	// Use a buffered writer to copy content
	reader := bufio.NewReader(input)
	writer := bufio.NewWriter(output)

	buffer := make([]byte, 1024)
	for {
		chunk, err := reader.Read(buffer)
		if err != nil && err != io.EOF {
			panic(err)
		}
		if chunk == 0 {
			break
		}

		// write
		_, err = writer.Write(buffer[:chunk])
		if err != nil {
			panic(err)
		}

	}

	if err := writer.Flush(); err != nil {
		panic(err)
	}

	fmt.Printf("Copy %s into %s\n", inputFileName, outputFileName)

}
