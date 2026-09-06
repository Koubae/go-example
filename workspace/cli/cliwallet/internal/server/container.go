package server

import (
	"cliwallet/internal/repository"
	"context"
	"fmt"
	"log/slog"

	"gorm.io/gorm"
)

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
