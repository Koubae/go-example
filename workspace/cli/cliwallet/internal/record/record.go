package record

import (
	"cliwallet/internal/model"
	"time"

	"github.com/google/uuid"
)

type Record interface {
	// ToModel -- for Reads/Load only
	ToModel() any
}

type Account struct {
	UUIDModel
	Name string `gorm:"not null;uniqueIndex:idx_accounts_name"`

	CreatedAt time.Time
	UpdatedAt time.Time

	Wallets []Wallet `gorm:"foreignKey:AccountID"`
}

// NewAccount -- for writes. Should not leak outside repository
func NewAccount(m model.Account) *Account {
	return &Account{
		Name: m.Name,
	}
}

func (r *Account) ToModel() model.Account {
	return model.Account{
		ID:   r.ID,
		Name: r.Name,
	}
}

type Wallet struct {
	UUIDModel

	AccountID uuid.UUID `gorm:"type:uuid;not null;index;uniqueIndex:idx_account_wallet_name"`
	Name      string    `gorm:"not null;uniqueIndex:idx_account_wallet_name"`

	Currency model.Currency `gorm:"not null"`
	Balance  int64          `gorm:"not null"`

	CreatedAt time.Time
	UpdatedAt time.Time

	Account      Account       `gorm:"constraint:OnUpdate:CASCADE,OnDelete:RESTRICT;"`
	Transactions []Transaction `gorm:"foreignKey:WalletID"`
}

// NewWallet -- for writes. Should not leak outside repository
func NewWallet(m model.Wallet) Wallet {
	return Wallet{
		AccountID: m.AccountID,
		Name:      m.Name,
		Currency:  m.Currency,
		Balance:   m.Balance.Amount,
	}
}

func (r Wallet) ToModel() model.Wallet {
	return model.Wallet{
		ID:        r.ID,
		AccountID: r.AccountID,
		Name:      r.Name,
		Currency:  r.Currency,
		Balance:   model.NewMoney(r.Balance, r.Currency),
	}
}

type Transaction struct {
	UUIDModel

	WalletID uuid.UUID `gorm:"type:uuid;not null;index"`

	Type     model.TransactionType `gorm:"not null"`
	Currency model.Currency        `gorm:"not null"`
	Amount   int64                 `gorm:"not null"`

	CreatedAt time.Time

	Wallet Wallet `gorm:"foreignKey:WalletID"`
}

// NewTransaction -- for writes. Should not leak outside repository
func NewTransaction(m model.Transaction) Transaction {
	return Transaction{
		WalletID: m.WalletID,
		Type:     m.Type,
		Currency: m.Currency,
		Amount:   m.Amount.Amount,
	}
}

func (r Transaction) ToModel() model.Transaction {
	return model.Transaction{
		ID:       r.ID,
		WalletID: r.WalletID,
		Type:     r.Type,
		Currency: r.Currency,
		Amount:   model.NewMoney(r.Amount, r.Currency),
	}
}
