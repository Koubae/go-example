package repository

import (
	"cliwallet/internal/record"
	"fmt"

	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func DBSQLite() (*gorm.DB, error) {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	sqlDB, err := db.DB()
	if err != nil {
		return nil, err
	}

	sqlDB.SetMaxOpenConns(1)

	return db, nil
}

type PostgresConfig struct {
	User     string
	Password string
	DB       string
	Host     string
	Port     string
	SSLMode  string
}

func DBPostgres(config PostgresConfig) (*gorm.DB, error) {
	dsn := fmt.Sprintf(
		"user=%s password=%s dbname=%s host=%s port=%s sslmode=%s connect_timeout=5",
		config.User,
		config.Password,
		config.DB,
		config.Host,
		config.Port,
		config.SSLMode,
	)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	sqlDB, err := db.DB()
	if err != nil {
		return nil, err
	}

	if err := sqlDB.Ping(); err != nil {
		return nil, fmt.Errorf("ping postgres: %w", err)
	}

	return db, nil
}

func MigrateAll(db *gorm.DB) error {
	tables := []any{
		&record.Account{},
		&record.Wallet{},
		&record.Transaction{},
	}
	for _, table := range tables {
		if err := db.AutoMigrate(table); err != nil {
			return err
		}
	}

	return nil
}
