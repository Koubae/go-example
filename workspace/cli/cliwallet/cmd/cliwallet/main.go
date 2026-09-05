package main

import (
	"cliwallet/internal/cli"
	"fmt"
	"os"
)

func main() {
	if err := cli.Cli.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
