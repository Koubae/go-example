package main

import (
	"errors"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"runtime"

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

func loadEnvFile(envFileName string) error {
	root, err := findProjectRoot()
	if err != nil {
		return errors.New("failed to find project root")
	}
	envPath := filepath.Join(root, envFileName)
	_ = godotenv.Load(envPath)

	return nil
}

// TODO: This func was AI created (cursor - Composer 1.5) in a rush
// TODO: Need to check whether there is a better way to handle this.. not sure
// I really don't like this approach.
func findProjectRoot() (string, error) {
	// 1. Explicit env var (production) – most reliable
	if root := os.Getenv("APP_ROOT"); root != "" {
		if abs, err := filepath.Abs(root); err == nil {
			return abs, nil
		}
	}
	// 2. Development: find go.mod from caller's locatio
	if _, file, _, ok := runtime.Caller(0); ok {
		dir := filepath.Dir(file)
		log.Printf("finding project root from caller's location %v\n", dir)
		for {
			if _, err := os.Stat(filepath.Join(dir, "go.mod")); err == nil {
				return dir, nil
			}
			parent := filepath.Dir(dir)

			log.Printf("finding project root from caller's location %v \n", parent)
			if parent == dir {
				break
			}
			dir = parent
		}
	}
	// 3. Fallback: directory of the executable (binary + .env in same dir)
	if execPath, err := os.Executable(); err == nil {
		return filepath.Dir(execPath), nil
	}
	// 4. Last resort: current working directory
	return os.Getwd()
}
