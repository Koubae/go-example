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
)

type ApplicationConfig struct {
	ServerHost            string
	ServerPort            string
	ServerAddress         string
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

	server := newServer(a)

	errCh := make(chan error, 1)
	go func() {
		a.logger.Info("starting server", "address", a.config.ServerAddress)

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
