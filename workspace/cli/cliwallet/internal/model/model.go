package model

import (
	"time"

	"github.com/google/uuid"
)

type TransactionType string

const (
	TransactionDeposit    TransactionType = "deposit"
	TransactionWithdrawal TransactionType = "withdrawal"
)

type Account struct {
	ID   uuid.UUID `json:"id"`
	Name string    `json:"name"`
}

type Wallet struct {
	ID        uuid.UUID `json:"id"`
	AccountID uuid.UUID `json:"account_id"`
	Name      string    `json:"name"`
	Currency  Currency  `json:"currency"`
	Balance   Money     `json:"balance"`
}

type Transaction struct {
	ID          uuid.UUID       `json:"id"`
	WalletID    uuid.UUID       `json:"wallet_id"`
	Type        TransactionType `json:"type"`
	Currency    Currency        `json:"currency"`
	Amount      Money           `json:"amount"`
	Description string          `json:"description"`
	CreatedAt   time.Time       `json:"created_at"`
}
