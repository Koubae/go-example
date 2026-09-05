package record

import (
	"cliwallet/internal/model"
	"time"
)

type Account struct {
	ID   uint   `gorm:"primaryKey"`
	Name string `gorm:"not null"`

	Wallets []Wallet `gorm:"foreignKey:AccountID"`
}

type Wallet struct {
	ID uint `gorm:"primaryKey"`

	AccountID uint   `gorm:"not null;index;uniqueIndex:idx_account_wallet_name"`
	Name      string `gorm:"not null;uniqueIndex:idx_account_wallet_name"`

	Currency model.Currency `gorm:"not null"`
	Balance  int64          `gorm:"not null"`

	CreatedAt time.Time
	UpdatedAt time.Time

	Account      Account       `gorm:"constraint:OnUpdate:CASCADE,OnDelete:RESTRICT;"`
	Transactions []Transaction `gorm:"foreignKey:WalletID"`
}

type Transaction struct {
	ID uint `gorm:"primaryKey"`

	WalletID uint `gorm:"not null;index"`

	Type   model.TransactionType `gorm:"not null"`
	Amount int64                 `gorm:"not null"`

	CreatedAt time.Time

	Wallet Wallet `gorm:"foreignKey:WalletID"`
}
