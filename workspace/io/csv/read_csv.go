package main

import (
	"encoding/csv"
	"fmt"
	"os"
	"strings"
)

func main() {
	const csvFilePath = "username.csv"

	csvFile, err := os.Open(csvFilePath)
	if err != nil {
		panic(err)
	}
	defer csvFile.Close()

	reader := csv.NewReader(csvFile)

	rows, err := reader.ReadAll()
	if err != nil {
		panic(err)
	}

	headers := rows[0]
	for _, header := range headers {
		fmt.Printf("%s, ", header)
	}
	fmt.Println("\n===========================================")
	for i, row := range rows[1:] {
		fmt.Printf("%d) %s\n", i+1, strings.Join(row, ", "))
	}
}
