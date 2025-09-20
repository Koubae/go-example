package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"strings"
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

	ReadJSONByToken()
	ReadJSONTokenByToken()

}

// @docs https://pkg.go.dev/encoding/json#NewDecoder
// @docs https://pkg.go.dev/encoding/json#Decoder.More
func ReadJSONByToken() {
	fmt.Println(
		`
==============================================
Read JSON by token	
==============================================`,
	)

	const jsonString = `
	[
		{"Name": "Ed", "Text": "Knock knock."},
		{"Name": "Sam", "Text": "Who's there?"},
		{"Name": "Ed", "Text": "Go fmt."},
		{"Name": "Sam", "Text": "Go fmt who?"},
		{"Name": "Ed", "Text": "Go fmt yourself!"}
	]
	`

	type Message struct {
		Name string `json:"Name"`
		Text string `json:"Text"`
	}

	decoder := json.NewDecoder(strings.NewReader(jsonString))

	// ---------------------------
	// Read root value/token of JSON it could be an object {} or and Array []
	// ---------------------------
	t, err := decoder.Token()
	if err != nil {
		log.Fatal(err)
	}
	delimiter, ok := t.(json.Delim)
	fmt.Printf("%T: %v (opening)\n", t, t)

	isObject := delimiter == '{'
	isArray := delimiter == '['
	fmt.Printf("ok=%v, Is Object: %v, Is Array: %v\n", ok, isObject, isArray)

	// Now we read each item inside the JSON Array
	message := make([]Message, 0)
	for decoder.More() {

		// t, err := decoder.Token()	// Doing this will throw an error since it will consume next token in JSON
		// fmt.Printf("Token: %v\n", t)

		var m Message
		err = decoder.Decode(&m) // Consume the entire Object till the next Object or Array
		if err != nil {
			log.Fatal(err)
		}

		fmt.Printf("%+v\n", m)
		message = append(message, m)
	}

	// read / consume closing token
	t, err = decoder.Token()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("%T: %v (closing)\n", t, t)

}

func ReadJSONTokenByToken() {
	fmt.Println(
		`
==============================================
Read Token by token
==============================================`,
	)

	const jsonString = `
{
    "manifest": {
        "name": "users",
        "created": "2025-09-20T11:48:25.777153552Z",
        "updated": "2025-09-20T11:48:25.777153552Z"
    },
    "records": {
      "key": "value",
      "key2": "value2",
      "key3": "value3"
    }
}
	`

	decoder := json.NewDecoder(strings.NewReader(jsonString))

	for i := 1; ; i++ {
		t, err := decoder.Token()
		if errors.Is(err, io.EOF) {
			fmt.Println("Reached EOF")
			break
		}
		if err != nil {
			log.Fatal(err)
		}
		hasMore := decoder.More()
		fmt.Printf("(iteration %d) -- %T: %v | more=%v\n", i, t, t, hasMore)
	}

	// for decoder.More() {
	// 	t, err := decoder.Token()
	// 	if err != nil {
	// 		log.Fatal(err)
	// 	}
	// 	fmt.Printf("(iteration %d) -- %T: %v\n", i, t, t)
	// 	i++
	// }

}
