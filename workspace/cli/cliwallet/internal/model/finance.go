package model

type TransactionType string

const (
	TransactionDeposit    TransactionType = "deposit"
	TransactionWithdrawal TransactionType = "withdrawal"
)
