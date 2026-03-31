package main

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
)

const PasswordLength = 32

func main() {
	password := make([]byte, PasswordLength)
	if _, err := rand.Read(password); err != nil {
		panic(fmt.Sprintf("err while gen pass, err: %v", err))
	}

	encodedPassword := base64.StdEncoding.EncodeToString(password)

	fmt.Println(encodedPassword)
}
