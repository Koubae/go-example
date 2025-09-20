package serialization

import (
	"errors"
	"fmt"
	"os"
)

func CreateFileWithSafeCloseDeferrable(path string) (*os.File, func(), error) {
	file, err := os.Create(path)
	if err != nil {
		return nil, nil, errors.New(fmt.Sprintf("Error creating file %s, error %v", path, err))
	}

	closer := func() {
		err := file.Close()
		if err != nil {
			fmt.Printf("Error closing file %s, error %v\n", path, err)
		}
	}
	return file, closer, nil
}

func OpenFileWithSafeCloseDeferrable(path string) (*os.File, func(), error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, nil, errors.New(fmt.Sprintf("Error opening file %s, error %v", path, err))
	}

	closer := func() {
		err := file.Close()
		if err != nil {
			fmt.Printf("Error closing file %s, error %v\n", path, err)
		}
	}
	return file, closer, nil

}
