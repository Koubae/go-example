package repository

import (
	"cliwallet/internal/model"
	"cliwallet/internal/record"
	"context"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Wallet interface {
	Create(ctx context.Context, wallet model.Wallet) (model.Wallet, error)
	GetByID(ctx context.Context, id uuid.UUID) (model.Wallet, error)

	GetByAccountIDAndName(ctx context.Context, accountID uuid.UUID, name string) (model.Wallet, error)
	ListByAccountID(ctx context.Context, accountID uuid.UUID) ([]model.Wallet, error)
}

type wallet struct {
	db *gorm.DB
}

func NewWalletRepo(db *gorm.DB) Wallet {
	return &wallet{
		db: db,
	}
}

func (r *wallet) Create(ctx context.Context, wallet model.Wallet) (model.Wallet, error) {
	record := record.NewWallet(wallet)

	if err := r.db.WithContext(ctx).Create(record).Error; err != nil {
		return model.Wallet{}, translateError(err)
	}

	return record.ToModel(), nil
}

func (r *wallet) GetByID(ctx context.Context, id uuid.UUID) (model.Wallet, error) {
	var record record.Wallet
	if err := r.db.WithContext(ctx).Where("id = ?", id).First(&record).Error; err != nil {
		return model.Wallet{}, translateError(err)
	}

	return record.ToModel(), nil
}

func (r *wallet) GetByAccountIDAndName(ctx context.Context, accountID uuid.UUID, name string) (model.Wallet, error) {
	var record record.Wallet

	err := r.db.WithContext(ctx).
		Where("account_id = ? AND name = ?", accountID, name).
		First(&record).
		Error
	if err != nil {
		return model.Wallet{}, translateError(err)
	}

	return record.ToModel(), nil
}

func (r *wallet) ListByAccountID(ctx context.Context, accountID uuid.UUID) ([]model.Wallet, error) {
	return []model.Wallet{}, nil
}
