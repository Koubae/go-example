/*

How connection return works

With Pool.Exec, pgxpool acquires a connection for the call and returns it when Exec returns.
With Pool.Query, the connection is held by Rows, so you must defer rows.Close().
With Pool.QueryRow, the connection is returned when you call Scan().

// OK: connection returned automatically after Exec returns
_, err := pool.Exec(ctx, "...")

// OK: connection returned when rows.Close() happens
rows, err := pool.Query(ctx, "...")
if err != nil { ... }
defer rows.Close()

// OK: connection returned when Scan() happens
err := pool.QueryRow(ctx, "...").Scan(&x)


But this is bad: 🧨🧨🧨
rows, _ := pool.Query(ctx, "SELECT ...")
// forgot rows.Close() -> leaked checked-out connection until GC/finalization path

And this is also bad: 🧨🧨🧨

row := pool.QueryRow(ctx, "SELECT ...")
// forgot Scan() -> connection not returned yet
_ = row

====================
Transaction rules
====================
Begin / BeginTx acquire a pooled connection and start a transaction.
You must finish it with Commit() or Rollback().

Also, pgxpool docs explicitly note that context cancellation does not auto-rollback the transaction,
so the defer tx.Rollback(ctx) pattern is the right production pattern.
Rollback() is safe even if Commit() already succeeded.

==========================

In a real HTTP API, you usually do this once at startup:

func main() {
	ctx := context.Background()

	pool, err := NewPool(ctx, mustGetEnv("DATABASE_URL"))
	if err != nil {
		panic(err)
	}
	defer pool.Close()

	store := NewStore(pool)

	_ = store
	// pass store into services / handlers
}

// ------------ Then inside handlers, pass the request context down:
func (h *Handler) GetUser(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	user, err := h.store.GetUserByID(ctx, 123)
	if err != nil {
		// handle
		return
	}

	_ = user
}

*/

package db

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type User struct {
	ID        int64
	Email     string
	Name      string
	CreatedAt time.Time
}

type Store struct {
	Pool *pgxpool.Pool
}

func NewPool(ctx context.Context, databaseURL string) (*pgxpool.Pool, error) {
	cfg, err := pgxpool.ParseConfig(databaseURL)
	if err != nil {
		return nil, fmt.Errorf("parse pool config: %w", err)
	}

	// Pool tuning: example values, adjust to your DB/server capacity.
	cfg.MaxConns = 30
	cfg.MinIdleConns = 5
	cfg.MaxConnLifetime = 1 * time.Hour
	cfg.MaxConnLifetimeJitter = 5 * time.Minute
	cfg.MaxConnIdleTime = 30 * time.Minute
	cfg.HealthCheckPeriod = 1 * time.Minute

	// Optional: run per-connection setup when a new connection is created.
	cfg.AfterConnect = func(ctx context.Context, conn *pgx.Conn) error {
		_, err := conn.Exec(ctx, `
			SET application_name = 'game-hangar-api'
		`)
		return err
	}

	pool, err := pgxpool.NewWithConfig(ctx, cfg)
	if err != nil {
		return nil, fmt.Errorf("create pool: %w", err)
	}

	// Important: pool creation itself does not prove the DB is reachable.
	pingCtx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	if err := pool.Ping(pingCtx); err != nil {
		pool.Close()
		return nil, fmt.Errorf("ping db: %w", err)
	}

	return pool, nil
}

func NewStore(pool *pgxpool.Pool) *Store {
	return &Store{Pool: pool}
}

func (s *Store) Close() {
	if s.Pool != nil {
		s.Pool.Close()
	}
}

// ------------------------------------------------------------
// 1. INSERT
// ------------------------------------------------------------

func (s *Store) InsertUser(ctx context.Context, email, name string) (int64, error) {
	const q = `
		INSERT INTO users (email, name)
		VALUES ($1, $2)
		RETURNING id
	`

	var id int64
	if err := s.Pool.QueryRow(ctx, q, email, name).Scan(&id); err != nil {
		return 0, fmt.Errorf("insert user: %w", err)
	}

	return id, nil
}

// Alternative INSERT when you only care about rows affected.
func (s *Store) InsertAuditLog(ctx context.Context, userID int64, action string) error {
	const q = `
		INSERT INTO audit_log (user_id, action)
		VALUES ($1, $2)
	`

	tag, err := s.Pool.Exec(ctx, q, userID, action)
	if err != nil {
		return fmt.Errorf("insert audit log: %w", err)
	}

	if tag.RowsAffected() != 1 {
		return fmt.Errorf("insert audit log: expected 1 row affected, got %d", tag.RowsAffected())
	}

	return nil
}

// ------------------------------------------------------------
// 2. QUERY (multiple rows)
// ------------------------------------------------------------

func (s *Store) ListUsers(ctx context.Context, limit int32) ([]User, error) {
	const q = `
		SELECT id, email, name, created_at
		FROM users
		ORDER BY id
		LIMIT $1
	`

	rows, err := s.Pool.Query(ctx, q, limit)
	if err != nil {
		return nil, fmt.Errorf("list users query: %w", err)
	}
	defer rows.Close() // VERY IMPORTANT: returns connection to pool

	users := make([]User, 0, limit)

	for rows.Next() {
		var u User
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

// ------------------------------------------------------------
// 3. QUERYROW (single row)
// ------------------------------------------------------------

func (s *Store) GetUserByID(ctx context.Context, id int64) (*User, error) {
	const q = `
		SELECT id, email, name, created_at
		FROM users
		WHERE id = $1
	`

	var u User
	err := s.Pool.QueryRow(ctx, q, id).Scan(
		&u.ID,
		&u.Email,
		&u.Name,
		&u.CreatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil // or return your domain NotFound error
		}
		return nil, fmt.Errorf("get user by id: %w", err)
	}

	return &u, nil
}

// ------------------------------------------------------------
// 4. TRANSACTION + ROLLBACK
// ------------------------------------------------------------

func (s *Store) CreateUserWithAudit(ctx context.Context, email, name string) (int64, error) {
	tx, err := s.Pool.BeginTx(ctx, pgx.TxOptions{
		IsoLevel: pgx.ReadCommitted,
	})
	if err != nil {
		return 0, fmt.Errorf("begin tx: %w", err)
	}

	// Safe pattern: always defer rollback.
	// If Commit succeeds first, this rollback becomes a harmless no-op.
	defer func() {
		_ = tx.Rollback(ctx)
	}()

	var userID int64
	err = tx.QueryRow(ctx, `
		INSERT INTO users (email, name)
		VALUES ($1, $2)
		RETURNING id
	`, email, name).Scan(&userID)
	if err != nil {
		return 0, fmt.Errorf("insert user in tx: %w", err)
	}

	tag, err := tx.Exec(ctx, `
		INSERT INTO audit_log (user_id, action)
		VALUES ($1, $2)
	`, userID, "user_created")
	if err != nil {
		return 0, fmt.Errorf("insert audit in tx: %w", err)
	}
	if tag.RowsAffected() != 1 {
		return 0, fmt.Errorf("insert audit in tx: expected 1 row affected, got %d", tag.RowsAffected())
	}

	if err := tx.Commit(ctx); err != nil {
		if errors.Is(err, pgx.ErrTxCommitRollback) {
			return 0, fmt.Errorf("commit resulted in rollback: %w", err)
		}
		return 0, fmt.Errorf("commit tx: %w", err)
	}

	return userID, nil
}

// ------------------------------------------------------------
// BONUS: Manual Acquire / Release
// Use this only when you need a pinned connection.
// ------------------------------------------------------------

func (s *Store) DoPinnedConnectionWork(ctx context.Context) error {
	conn, err := s.Pool.Acquire(ctx)
	if err != nil {
		return fmt.Errorf("acquire conn: %w", err)
	}
	defer conn.Release()

	// Anything here uses the same underlying connection.
	_, err = conn.Exec(ctx, `SET LOCAL statement_timeout = '5s'`)
	if err != nil {
		return fmt.Errorf("set local timeout: %w", err)
	}

	var now time.Time
	if err := conn.QueryRow(ctx, `SELECT NOW()`).Scan(&now); err != nil {
		return fmt.Errorf("query now: %w", err)
	}

	return nil
}
