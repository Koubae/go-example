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

type UserService struct {
	pool *pgxpool.Pool
	repo *repository.UserRepository
}

func NewUserService(pool *pgxpool.Pool, repo *repository.UserRepository) *UserService {
	return &UserService{
		pool: pool,
		repo: repo,
	}
}

func (s *UserService) CreateUser(ctx context.Context, email, name string) (*model.User, error) {
	tx, err := s.pool.BeginTx(ctx, pgx.TxOptions{
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
	return s.repo.GetByID(ctx, s.pool, id)
}

func (s *UserService) ListUsers(ctx context.Context, limit, offset int32) ([]model.User, error) {
	return s.repo.List(ctx, s.pool, limit, offset)
}

func (s *UserService) IsNotFound(err error) bool {
	return errors.Is(err, repository.ErrNotFound)
}
