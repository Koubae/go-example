/*
Source: How I write HTTP services in Go after 13 years
	https://grafana.com/blog/how-i-write-http-services-in-go-after-13-years/

*/

package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"sync/atomic"
	"syscall"
	"time"
)

const (
	port = ":8080"
)

func main() {
	log.SetFlags(log.Ldate | log.LUTC | log.Lmicroseconds | log.Lshortfile)
	logger := log.New(log.Writer(), "http: ", log.LstdFlags)

	ctx := context.Background()
	if err := run(ctx, logger); err != nil {
		log.Fatalf("error running server: %v\n", err)
		os.Exit(1)
	}

	log.Println("Server stopped")

}

func run(ctx context.Context, logger *log.Logger) error {
	ctx, cancel := signal.NotifyContext(ctx, os.Interrupt, syscall.SIGTERM)
	defer cancel()

	logger.Println("Server process up")

	mux := http.NewServeMux()

	mux.Handle("/", http.NotFoundHandler())
	mux.Handle("GET /users", handleGetUser(logger))
	mux.Handle("GET /users/{id}", handleGetUserByID(logger))
	mux.Handle("POST /account", isAdmin(handleCreateAccount(logger)))
	mux.Handle("GET /account", isAdmin(handleListAccount(logger)))

	srv := &http.Server{
		Addr:           port,
		Handler:        mux,
		ErrorLog:       logger,
		ReadTimeout:    5 * time.Minute, // 5 minutes
		WriteTimeout:   5 * time.Minute, // 5 minutes
		MaxHeaderBytes: 1 << 20,         // 1 MB
	}

	var runErr error

	go func() {
		log.Printf("Server starting listening on %s\n", port)
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("error listen and serve, error: %v\n", err)
			runErr = fmt.Errorf("error listen and serve, error: %w", err)
		}
	}()

	// simple wait for interrupt signal
	// quit := make(chan os.Signal, 1)
	// signal.Notify(quit, os.Interrupt)
	// <-quit
	// log.Println("Server shutting down...")

	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		<-ctx.Done()

		shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		if err := srv.Shutdown(shutdownCtx); err != nil {
			log.Fatalf("Server forced to shutdown: %v", err)
			runErr = fmt.Errorf("server forced to shutdown: %w", err)
		}

	}()
	wg.Wait()
	return runErr

}

// --------------------------------------------
//
//	Handlers Decorators
//
// --------------------------------------------
type Account struct {
	ID   string `json:"id"`
	Name string `json:"name"`
	Role string `json:"role"`
}

var accounts = map[int]Account{
	1:  {ID: "1", Name: "Alice", Role: "admin"},
	2:  {ID: "2", Name: "Bob", Role: "user"},
	3:  {ID: "3", Name: "Charlie", Role: "user"},
	10: {ID: "10", Name: "Default", Role: "user"},
}

const defaultAccountID = 10

func isAdmin(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		idRaw := r.Header.Get("X-Account-ID")
		id := defaultAccountID
		if idRaw != "" {
			_, _ = fmt.Sscanf(idRaw, "%d", &id)
		}

		account, ok := accounts[id]
		if !ok || account.Role != "admin" {
			_ = encode(w, r, http.StatusForbidden, map[string]string{"error": "Admin access required"})
			return
		}

		h.ServeHTTP(w, r)

	})
}

func handleGetUser(logger *log.Logger) http.Handler {
	type user struct {
		Name string `json:"name"`
		Age  int    `json:"age"`
	}

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		response := []user{
			{Name: "John Doe", Age: 45},
			{Name: "Alice", Age: 30},
			{Name: "Bob", Age: 25},
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		err := json.NewEncoder(w).Encode(response)
		if err != nil {
			logger.Println("Error encoding JSON:", err)
		}
	})
}

func handleGetUserByID(logger *log.Logger) http.Handler {
	type user struct {
		Name string `json:"name"`
		Age  int    `json:"age"`
	}

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		id := r.PathValue("id")

		logger.Printf("Fetching user with ID: %s\n", id)

		response := user{Name: "John Doe", Age: 45}
		if err := encode(w, r, http.StatusOK, response); err != nil {
			logger.Printf("Error encoding JSON response: %v\n", err)
		}
	})
}

// handleGetUser
// curl -X POST -H "X-Account-ID: 1" http://localhost:8080/account -d '{"name": "JohnWick", "role":"admin"}'
func handleCreateAccount(logger *log.Logger) http.Handler {
	type Request struct {
		Name string `json:"name"`
		Role string `json:"role"`
	}
	type Response struct {
		ID string `json:"id"`
	}

	var counter atomic.Int64
	counter.Store(int64(defaultAccountID))
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		request, err := decode[Request](r)
		if err != nil {
			logger.Printf("Error decoding request body: %v\n", err)
			http.Error(w, "Bad Request", http.StatusBadRequest)
			return
		}

		counter.Add(1)
		newID := int(counter.Load())

		account := Account{
			ID:   fmt.Sprintf("%d", newID),
			Name: request.Name,
			Role: request.Role,
		}
		accounts[newID] = account

		response := Response{ID: account.ID}
		if err := encode(w, r, http.StatusCreated, response); err != nil {
			logger.Printf("Error encoding JSON response: %v\n", err)
		}

	})
}

// handleListAccount
// curl -H "X-Account-ID: 1" http://localhost:8080/account
func handleListAccount(logger *log.Logger) http.Handler {

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if err := encode(w, r, http.StatusOK, accounts); err != nil {
			logger.Printf("Error encoding JSON response: %v\n", err)
		}
	})
}

func encode[T any](w http.ResponseWriter, _ *http.Request, status int, v T) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(v); err != nil {
		return fmt.Errorf("error encoding response: %w", err)
	}
	return nil
}

func decode[T any](r *http.Request) (T, error) {
	var v T
	if err := json.NewDecoder(r.Body).Decode(&v); err != nil {
		return v, fmt.Errorf("error decoding request body: %w", err)
	}
	return v, nil
}
