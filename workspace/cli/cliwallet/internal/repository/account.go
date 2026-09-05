package repository

import (
	"cliwallet/internal/model"
	"cliwallet/internal/record"
	"context"

	"gorm.io/gorm"
)

type Account interface {
	Create(ctx context.Context, account model.Account) (model.Account, error)
}

type account struct {
	db *gorm.DB
}

func NewAccountRepo(db *gorm.DB) Account {
	return &account{
		db: db,
	}
}

func (r *account) Create(ctx context.Context, account model.Account) (model.Account, error) {
	record := record.NewAccount(account)

	if err := r.db.WithContext(ctx).Create(record).Error; err != nil {
		return model.Account{}, translateError(err)
	}

	return record.ToModel(), nil
}
