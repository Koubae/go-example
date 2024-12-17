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

	inputFile, err := os.Open(inputFileName)
	if err != nil {
		panic(err)
	}
	defer inputFile.Close()

	outputFile, err := os.Create(outputFileName)
	if err != nil {
		panic(err)
	}
	defer outputFile.Close()

	// Use a buffered writer to copy content
	reader := bufio.NewReader(inputFile)
	writer := bufio.NewWriter(outputFile)

	_, err = io.Copy(writer, reader)
	if err != nil {
		fmt.Printf("Error copying data: %v\n", err)
		return
	}

	// Ensure all data is flushed to the destination file
	err = writer.Flush()
	if err != nil {
		fmt.Printf("Error flushing data to destination file: %v\n", err)
		return
	}

	fmt.Println("File content copied successfully!")

}
