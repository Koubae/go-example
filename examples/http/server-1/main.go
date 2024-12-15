package main

import (
	"fmt"
	"net/http"
)

func main() {
	fmt.Println("Starting Server ...")

	// Register Routes
	http.HandleFunc("/", index)
	fmt.Println("Listening on port 8080 ...")
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		panic(fmt.Sprintf("Error starting server: %v", err))
	}
}

// ==================================
// Handlers
// ==================================
func index(w http.ResponseWriter, req *http.Request) {
	fmt.Println("[GET] / 200")
	_, err := fmt.Fprintf(w, "Welcome to Go-Example Simple Server 1\n")
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Printf("Something went wrong: %v\n", err)
		return
	}
}
