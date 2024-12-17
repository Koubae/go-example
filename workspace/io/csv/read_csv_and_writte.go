package main

import (
	"encoding/csv"
	"fmt"
	"os"
	"strings"
)

func main() {
	const csvFilePath = "username.csv"
	const csvOutFilePath = "username_out.csv"

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

	output, err := os.Create(csvOutFilePath)
	if err != nil {
		panic(err)
	}
	defer output.Close()

	writer := csv.NewWriter(output)
	writer.Comma = ','

	err = writer.WriteAll(rows)
	if err != nil {
		panic(err)
	}

	fmt.Printf("New csv file created: %s\n", csvOutFilePath)

}
