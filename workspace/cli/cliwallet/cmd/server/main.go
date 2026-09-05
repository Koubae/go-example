package main

import (
	"cliwallet/internal/repository"
	"fmt"
	"os"
)

func main() {
	fmt.Print("Hello from the server!")

	db, err := repository.DBPostgres(repository.PostgresConfig{
		User:     "admin",
		Password: "admin",
		DB:       "cli_wallet",
		Host:     "127.0.0.1",
		Port:     "5432",
		SSLMode:  "disable",
	})
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	fmt.Println("Database connection Established", db)

	if err := repository.MigrateAll(db); err != nil {
		fmt.Printf("Error while migrate db, error: %s\n", err.Error())
		os.Exit(1)
	}
	fmt.Println("Database migration completed sucessfully")
}
