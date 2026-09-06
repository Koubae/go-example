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
	if limit <= 0 {
		return nil, ErrQueryInvalidLimit
	}
	if offset < 0 {
		return nil, ErrQueryInvalidOffSet
	}

	var records []record.Transaction
	err := r.db.WithContext(ctx).
		Where("wallet_id = ?", walletID).
		Limit(limit).
		Offset(offset).
		Find(&records).
		Error
	if err != nil {
		return nil, translateError(err)
	}

	models := make([]model.Transaction, 0, len(records))
	for _, rec := range records {
		models = append(models, rec.ToModel())
	}
	return models, nil
}
