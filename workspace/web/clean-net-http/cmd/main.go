package main

import (
	"clean-net-http/internal/platform/postgres"
	"clean-net-http/internal/user/repository"
	"clean-net-http/internal/user/service"
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	httpuser "clean-net-http/internal/user/http"
)

func main() {
	ctx := context.Background()

	databaseURL := os.Getenv("DATABASE_URL")
	if databaseURL == "" {
		log.Fatal("DATABASE_URL is required")
	}

	pool, err := postgres.NewPool(ctx, databaseURL)
	if err != nil {
		log.Fatalf("db init failed: %v", err)
	}
	defer pool.Close()

	userRepo := repository.NewUserRepository()
	userService := service.NewUserService(pool, userRepo)
	userHandler := httpuser.NewHandler(userService)

	mux := http.NewServeMux()
	userHandler.RegisterRoutes(mux)

	server := &http.Server{
		Addr:              ":8080",
		Handler:           mux,
		ReadHeaderTimeout: 5 * time.Second,
	}

	go func() {
		log.Printf("listening on %s", server.Addr)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("http server failed: %v", err)
		}
	}()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)
	<-stop

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := server.Shutdown(shutdownCtx); err != nil {
		log.Printf("graceful shutdown failed: %v", err)
	}
}

