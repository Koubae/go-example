package server

import (
	"context"
	"log/slog"
	"net"
	"time"
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
		ServerHost:            "0.0.0.0",
		ServerPort:            "8080",
		ServerAddress:         net.JoinHostPort("0.0.0.0", "8080"),
		ServerReadTimeout:     1 * time.Minute,
		ServerWriteTimeout:    1 * time.Minute,
		ServerIdleTimeout:     1 * time.Minute,
		ServerMaxHeaderBytes:  1 << 20,
		ServerShutdownTimeout: 10 * time.Second,
	}

	app := NewApplication(container, logger, config)
	return app.Run(ctx)
}
