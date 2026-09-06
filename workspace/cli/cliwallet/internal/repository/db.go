package repository

import (
	"cliwallet/internal/record"
	"errors"
	"fmt"

	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var (
	ErrRecordDuplicate = errors.New("record already exists")
	ErrRecordNotFound  = errors.New("record not found")
	ErrInvalidRef      = errors.New("invalid reference")

	ErrQueryInvalidLimit  = errors.New("invalid query limit")
	ErrQueryInvalidOffSet = errors.New("invalid query offset")
)

func DBSQLite() (*gorm.DB, error) {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{
		TranslateError: true,
	})
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

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		TranslateError: true,
	})
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

func translateError(err error) error {
	switch {
	case errors.Is(err, gorm.ErrDuplicatedKey):
		return fmt.Errorf("%w: %v", ErrRecordDuplicate, err)

	case errors.Is(err, gorm.ErrForeignKeyViolated):
		return fmt.Errorf("%w: %v", ErrInvalidRef, err)

	case errors.Is(err, gorm.ErrRecordNotFound):
		return ErrRecordNotFound

	default:
		return err
	}
}
