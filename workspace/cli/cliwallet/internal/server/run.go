package server

import (
	"cliwallet/internal/repository"
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"os/signal"
	"syscall"
	"time"

	"gorm.io/gorm"
)

func Run(ctx context.Context, logger *slog.Logger, container *Container) error {
	defer func() {
		// Shutdown Dependencies
		if err := container.Close(); err != nil {
			logger.Error("error while closing container", "error", err)
		}
		logger.Info("container closed")
		logger.Info("server stopped")
	}()

	config := ApplicationConfig{ // TODO: load from config file
		ServerAddres:          ":8080",
		ServerReadTimeout:     1 * time.Minute,
		ServerWriteTimeout:    1 * time.Minute,
		ServerIdleTimeout:     1 * time.Minute,
		ServerMaxHeaderBytes:  1 << 20,
		ServerShutdownTimeout: 10 * time.Second,
	}

	app := NewApplication(container, logger, config)
	return app.Run(ctx)
}

type Container struct {
	db     *gorm.DB
	Logger *slog.Logger
}

func NewProductionContainer(ctx context.Context, logger *slog.Logger) (*Container, error) {
	db, err := repository.DBPostgres(repository.PostgresConfig{
		User:     "admin",
		Password: "admin",
		DB:       "cli_wallet",
		Host:     "127.0.0.1",
		Port:     "5432",
		SSLMode:  "disable",
	})
	if err != nil {
		return nil, fmt.Errorf("error while connecting to database: %w", err)
	}
	logger.Info("database connection established")

	return &Container{
		db:     db,
		Logger: logger,
	}, nil
}

func (c *Container) Close() error {
	db, err := c.db.DB()
	if err != nil {
		return fmt.Errorf("error while getting database connection: %w", err)
	}
	if err := db.Close(); err != nil {
		return fmt.Errorf("error while closing database connection: %w", err)
	}
	c.Logger.Info("database connection closed")
	return nil
}

type ApplicationConfig struct {
	ServerAddres          string
	ServerReadTimeout     time.Duration
	ServerWriteTimeout    time.Duration
	ServerIdleTimeout     time.Duration
	ServerMaxHeaderBytes  int
	ServerShutdownTimeout time.Duration
}

type Application struct {
	logger    *slog.Logger
	container *Container
	config    ApplicationConfig

	accountRepo     repository.Account
	walletRepo      repository.Wallet
	transactionRepo repository.Transaction
}

func NewApplication(container *Container, logger *slog.Logger, config ApplicationConfig) *Application {
	app := &Application{
		container:       container,
		logger:          logger,
		config:          config,
		accountRepo:     repository.NewAccountRepo(container.db),
		walletRepo:      repository.NewWalletRepo(container.db),
		transactionRepo: repository.NewTransactionRepo(container.db),
	}

	return app
}

func (a *Application) Run(ctx context.Context) error {
	ctx, cancel := signal.NotifyContext(ctx, syscall.SIGINT, syscall.SIGTERM)
	defer cancel()

	mux := http.NewServeMux()
	// TODO: setup handlers
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Hello, World!"))
	})

	server := &http.Server{
		Addr:           a.config.ServerAddres,
		ReadTimeout:    a.config.ServerReadTimeout,
		WriteTimeout:   a.config.ServerWriteTimeout,
		IdleTimeout:    a.config.ServerIdleTimeout,
		MaxHeaderBytes: a.config.ServerMaxHeaderBytes,
		Handler:        mux,
	}

	errCh := make(chan error, 1)
	go func() {
		a.logger.Info("starting server", "address", a.config.ServerAddres)

		errCh <- server.ListenAndServe()

	}()

	select {
	case err := <-errCh:
		if errors.Is(err, http.ErrServerClosed) {
			return nil
		}
		return fmt.Errorf("error listen and serve, error: %w", err)
	case <-ctx.Done():

		shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), a.config.ServerShutdownTimeout)
		defer shutdownCancel()

		if err := server.Shutdown(shutdownCtx); err != nil {
			return fmt.Errorf("shutdown server: %w", err)
		}
		a.logger.Info("server shutdown completed")

		err := <-errCh // ListenAndServe has returned
		if err != nil && !errors.Is(err, http.ErrServerClosed) {
			return fmt.Errorf("listen and serve: %w", err)
		}
		return nil
	}

}
