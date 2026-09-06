package main

import (
	"cliwallet/internal/server"
	"context"
	"log/slog"
	"os"
)

func main() {
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))
	slog.SetDefault(logger)
	logger.Info("cliwallet server starting")

	ctx := context.Background()
	container, err := server.NewProductionContainer(ctx, logger)
	if err != nil {
		logger.Error("failed to create container", "error", err)
		os.Exit(1)
	}
	if err := server.Run(ctx, logger, container); err != nil {
		logger.Error("server failed", "error", err)
		os.Exit(1)
	}

}
