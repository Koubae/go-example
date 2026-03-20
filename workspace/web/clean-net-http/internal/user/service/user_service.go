package service

import (
	"clean-net-http/internal/user/model"
	"clean-net-http/internal/user/repository"
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Tx interface {
	repository.DBTX
	Commit(ctx context.Context) error
	Rollback(ctx context.Context) error
}

type TxManager interface {
	BeginTx(ctx context.Context, txOptions pgx.TxOptions) (Tx, error)
}

type UserRepository interface {
	Create(ctx context.Context, db repository.DBTX, p repository.CreateUserParams) (*model.User, error)
	GetByID(ctx context.Context, db repository.DBTX, id int64) (*model.User, error)
	List(ctx context.Context, db repository.DBTX, limit, offset int32) ([]model.User, error)
	InsertAudit(ctx context.Context, db repository.DBTX, userID int64, action string) error
}

type pgxTxManager struct {
	pool *pgxpool.Pool
}

func newPgxTxManager(pool *pgxpool.Pool) TxManager {
	return &pgxTxManager{pool: pool}
}

type pgxTxAdapter struct {
	pgx.Tx
}

func (m *pgxTxManager) BeginTx(ctx context.Context, txOptions pgx.TxOptions) (Tx, error) {
	tx, err := m.pool.BeginTx(ctx, txOptions)
	if err != nil {
		return nil, err
	}

	return &pgxTxAdapter{Tx: tx}, nil
}

type UserService struct {
	db        repository.DBTX
	txManager TxManager
	repo      UserRepository
}

func NewUserService(pool *pgxpool.Pool, repo UserRepository) *UserService {
	return &UserService{
		db:        pool,
		txManager: newPgxTxManager(pool),
		repo:      repo,
	}
}

func (s *UserService) CreateUser(ctx context.Context, email, name string) (*model.User, error) {
	tx, err := s.txManager.BeginTx(ctx, pgx.TxOptions{
		IsoLevel: pgx.ReadCommitted,
	})
	if err != nil {
		return nil, fmt.Errorf("begin tx: %w", err)
	}
	defer func() {
		_ = tx.Rollback(ctx)
	}()

	user, err := s.repo.Create(ctx, tx, repository.CreateUserParams{
		Email: email,
		Name:  name,
	})
	if err != nil {
		return nil, err
	}

	if err := s.repo.InsertAudit(ctx, tx, user.ID, "user_created"); err != nil {
		return nil, err
	}

	if err := tx.Commit(ctx); err != nil {
		if errors.Is(err, pgx.ErrTxCommitRollback) {
			return nil, fmt.Errorf("commit failed and transaction was rolled back: %w", err)
		}
		return nil, fmt.Errorf("commit tx: %w", err)
	}

	return user, nil
}

func (s *UserService) GetUser(ctx context.Context, id int64) (*model.User, error) {
	return s.repo.GetByID(ctx, s.db, id)
}

func (s *UserService) ListUsers(ctx context.Context, limit, offset int32) ([]model.User, error) {
	return s.repo.List(ctx, s.db, limit, offset)
}

func (s *UserService) IsNotFound(err error) bool {
	return errors.Is(err, repository.ErrNotFound)
}
