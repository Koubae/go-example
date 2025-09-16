package main

import (
	"errors"
	"fmt"
	"io"
	"iter"
	"os"
)

func main() {
	fileName := "./workspace/io/hello-world.txt"

	fileDescriptor, err := os.Open(fileName)
	if err != nil {
		panic(err)
	}
	defer fileDescriptor.Close()

	var totalData int
	var cycles int

	buffer := make([]byte, 10) // make very tiny buffer to show off the "idea" of reading small chunks of data
	for {
		content, err := fileDescriptor.Read(buffer)
		if errors.Is(err, io.EOF) {
			break
		} else if err != nil {
			panic(err.Error())
		}

		totalData += content
		cycles++

		data := string(buffer[:content])
		fmt.Print(data)
	}

	fmt.Printf("\n\n============================================================================================\n")
	fmt.Printf("- Total data read: %d bytes in %d cycles\n", totalData, cycles)

	fmt.Println("============================================================================================\n")
	reader, err := ReadFileWithBufferIterator(fileName, 10)
	if err != nil {
		fmt.Printf("Error initializing  ReadFileWithBufferIterator, error: %v\n", err)
		return
	}

	for data, err := range reader {
		if err != nil {
			panic(err)
		}
		fmt.Print(data)
	}

}

func ReadFileWithBufferIterator(fileName string, bufferSize int) (iter.Seq2[string, error], error) {
	fileDescriptor, err := os.Open(fileName)
	if err != nil {
		return nil, err
	}

	buffer := make([]byte, bufferSize)
	return func(yield func(string, error) bool) {
		defer fileDescriptor.Close() // TODO: we need to double check if this is ok to do. I mean would be a problem to defer the closing of the file herE? any resource leaking?

		for {
			content, err := fileDescriptor.Read(buffer)
			if errors.Is(err, io.EOF) {
				break
			} else if err != nil {
				yield("", err)
				return
			}

			data := string(buffer[:content])
			yield(data, nil)
		}

	}, nil

}
