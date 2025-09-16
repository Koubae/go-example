package main

import (
	"encoding/json"
	"fmt"
	"os"
)

type Person struct {
	Name string `json:"name"`
	Age  int    `json:"age"`
}

func main() {
	p1 := Person{
		Name: "Federico",
		Age:  34,
	}

	jsonS, err := json.Marshal(p1)
	if err != nil {
		fmt.Printf("Error While marshalling json, error: %s\n", err.Error())
		return
	}

	fmt.Println(string(jsonS))
	fmt.Println(jsonS)

	_, err = os.Stdout.Write(jsonS)
	if err != nil {
		return
	}
	fmt.Println()

	var p2 Person
	err = json.Unmarshal(jsonS, &p2)
	if err != nil {
		fmt.Printf("Error While unmarshalling json, error: %s\n", err.Error())
		return
	}
	fmt.Println(p2)

}
