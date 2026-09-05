package repository

import "context"

type Wallet interface {
	Create(ctx context.Context) error
	Get(ctx context.Context) error
	ListTransactionHistory(ctx context.Context) error
}
