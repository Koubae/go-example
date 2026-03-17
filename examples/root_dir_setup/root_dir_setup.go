package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/joho/godotenv"
	vars "github.com/koubae/go-example"
)

// go build -o app .
// go build -o app examples/root_dir_setup/root_dir_setup.go
func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	f, _ := os.Getwd()
	gopath := os.Getenv("GOPATH")
	log.Printf("GOPATH: %s\n", gopath)
	log.Printf("Base: %s\n", filepath.Base(f))
	log.Printf("Dir: %s\n", filepath.Dir(f))

	log.Printf("RootDir: %s\n", vars.RootDir)

	envPath := filepath.Join(vars.RootDir, ".env.example")
	log.Printf("EnvPath: %s\n", envPath)
	err := godotenv.Load(envPath)
	if err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}

	envVar := os.Getenv("ENV_VAR_EXAMPLE")
	fmt.Println("ENV_VAR_EXAMPLE:", envVar)

}
