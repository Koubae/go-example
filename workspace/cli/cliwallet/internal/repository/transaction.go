package repository

import (
	"cliwallet/internal/model"
	"cliwallet/internal/record"
	"context"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Transaction interface {
	Create(ctx context.Context, tx *gorm.DB, transaction model.Transaction) (model.Transaction, error)
	GetByID(ctx context.Context, id uuid.UUID) (model.Transaction, error)
	ListByWalletID(ctx context.Context, walletID uuid.UUID, limit, offset int) ([]model.Transaction, error)
}

type transaction struct {
	db *gorm.DB
}

func NewTransactionRepo(db *gorm.DB) Transaction {
	return &transaction{
		db: db,
	}
}

func (r *transaction) Create(ctx context.Context, tx *gorm.DB, transaction model.Transaction) (model.Transaction, error) {
	record := record.NewTransaction(transaction)
	if err := tx.WithContext(ctx).Create(record).Error; err != nil {
		return model.Transaction{}, translateError(err)
	}

	return record.ToModel(), nil
}
func (r *transaction) GetByID(ctx context.Context, id uuid.UUID) (model.Transaction, error) {
	var record record.Transaction

	if err := r.db.WithContext(ctx).Where("id = ?", id).First(&record).Error; err != nil {
		return model.Transaction{}, translateError(err)
	}

	return record.ToModel(), nil
}
func (r *transaction) ListByWalletID(ctx context.Context, walletID uuid.UUID, limit, offset int) ([]model.Transaction, error) {
	return nil, nil
}
