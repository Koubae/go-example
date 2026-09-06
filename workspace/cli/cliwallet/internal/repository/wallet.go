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
	ListByAccountID(ctx context.Context, accountID uuid.UUID, limit, offset int) ([]model.Wallet, error)

	UpdateBalance(ctx context.Context, tx *gorm.DB, id uuid.UUID, balance int64) error
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

func (r *wallet) ListByAccountID(ctx context.Context, accountID uuid.UUID, limit, offset int) ([]model.Wallet, error) {
	if limit <= 0 {
		return nil, ErrQueryInvalidLimit
	}
	if offset < 0 {
		return nil, ErrQueryInvalidOffSet
	}

	var records []record.Wallet
	err := r.db.WithContext(ctx).
		Where("account_id = ?", accountID).
		Limit(limit).
		Offset(offset).
		Find(&records).
		Error
	if err != nil {
		return nil, translateError(err)
	}

	models := make([]model.Wallet, 0, len(records))
	for _, rec := range records {
		models = append(models, rec.ToModel())
	}
	return models, nil

}

func (r *wallet) UpdateBalance(ctx context.Context, tx *gorm.DB, id uuid.UUID, balance int64) error {
	// TODO: Implement this
	return nil
}
