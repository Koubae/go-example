package repository

import (
	"clean-net-http/internal/user/model"
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

var ErrNotFound = errors.New("resource not found")

type DBTX interface {
	Exec(ctx context.Context, sql string, args ...any) (pgconn.CommandTag, error)
	Query(ctx context.Context, sql string, args ...any) (pgx.Rows, error)
	QueryRow(ctx context.Context, sql string, args ...any) pgx.Row
}

type CreateUserParams struct {
	Email string
	Name  string
}

type UserRepository struct{}

func NewUserRepository() *UserRepository {
	return &UserRepository{}
}

func (r *UserRepository) Create(ctx context.Context, db DBTX, p CreateUserParams) (*model.User, error) {
	const q = `
		INSERT INTO users (email, name)
		VALUES ($1, $2)
		RETURNING id, email, name, created_at
	`

	var u model.User
	err := db.QueryRow(ctx, q, p.Email, p.Name).Scan(
		&u.ID,
		&u.Email,
		&u.Name,
		&u.CreatedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("create user: %w", err)
	}

	return &u, nil
}

func (r *UserRepository) GetByID(ctx context.Context, db DBTX, id int64) (*model.User, error) {
	const q = `
		SELECT id, email, name, created_at
		FROM users
		WHERE id = $1
	`

	var u model.User
	err := db.QueryRow(ctx, q, id).Scan(
		&u.ID,
		&u.Email,
		&u.Name,
		&u.CreatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrNotFound
		}
		return nil, fmt.Errorf("get user by id: %w", err)
	}

	return &u, nil
}

func (r *UserRepository) List(ctx context.Context, db DBTX, limit, offset int32) ([]model.User, error) {
	const q = `
		SELECT id, email, name, created_at
		FROM users
		ORDER BY id
		LIMIT $1 OFFSET $2
	`

	rows, err := db.Query(ctx, q, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("list users query: %w", err)
	}
	defer rows.Close()

	users := make([]model.User, 0, limit)

	for rows.Next() {
		var u model.User
		if err := rows.Scan(&u.ID, &u.Email, &u.Name, &u.CreatedAt); err != nil {
			return nil, fmt.Errorf("list users scan: %w", err)
		}
		users = append(users, u)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("list users rows err: %w", err)
	}

	return users, nil
}

func (r *UserRepository) InsertAudit(ctx context.Context, db DBTX, userID int64, action string) error {
	const q = `
		INSERT INTO audit_log (user_id, action)
		VALUES ($1, $2)
	`

	tag, err := db.Exec(ctx, q, userID, action)
	if err != nil {
		return fmt.Errorf("insert audit log: %w", err)
	}

	if tag.RowsAffected() != 1 {
		return fmt.Errorf("insert audit log: expected 1 row affected, got %d", tag.RowsAffected())
	}

	return nil
}
