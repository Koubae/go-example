package service

import (
	"clean-net-http/internal/user/model"
	"clean-net-http/internal/user/repository"
	"context"
	"testing"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockDB struct {
	mock.Mock
}

func (m *MockDB) Exec(ctx context.Context, sql string, args ...any) (pgconn.CommandTag, error) {
	panic("unexpected call to MockDB.Exec")
}

func (m *MockDB) Query(ctx context.Context, sql string, args ...any) (pgx.Rows, error) {
	panic("unexpected call to MockDB.Query")
}

func (m *MockDB) QueryRow(ctx context.Context, sql string, args ...any) pgx.Row {
	panic("unexpected call to MockDB.QueryRow")
}

type MockTx struct {
	mock.Mock
}

func (m *MockTx) Commit(ctx context.Context) error {
	args := m.Called(ctx)
	return args.Error(0)
}

func (m *MockTx) Rollback(ctx context.Context) error {
	args := m.Called(ctx)
	return args.Error(0)
}

func (m *MockTx) Exec(ctx context.Context, sql string, args ...any) (pgconn.CommandTag, error) {
	panic("unexpected call to MockTx.Exec")
}

func (m *MockTx) Query(ctx context.Context, sql string, args ...any) (pgx.Rows, error) {
	panic("unexpected call to MockTx.Query")
}

func (m *MockTx) QueryRow(ctx context.Context, sql string, args ...any) pgx.Row {
	panic("unexpected call to MockTx.QueryRow")
}

type MockTxManager struct {
	mock.Mock
}

func (m *MockTxManager) BeginTx(ctx context.Context, txOptions pgx.TxOptions) (Tx, error) {
	args := m.Called(ctx, txOptions)
	tx, _ := args.Get(0).(Tx)
	return tx, args.Error(1)
}

type MockUserRepository struct {
	mock.Mock
}

func (m *MockUserRepository) Create(ctx context.Context, db repository.DBTX, p repository.CreateUserParams) (*model.User, error) {
	args := m.Called(ctx, db, p)
	user, _ := args.Get(0).(*model.User)
	return user, args.Error(1)
}

func (m *MockUserRepository) GetByID(ctx context.Context, db repository.DBTX, id int64) (*model.User, error) {
	args := m.Called(ctx, db, id)
	user, _ := args.Get(0).(*model.User)
	return user, args.Error(1)
}

func (m *MockUserRepository) List(ctx context.Context, db repository.DBTX, limit, offset int32) ([]model.User, error) {
	args := m.Called(ctx, db, limit, offset)
	users, _ := args.Get(0).([]model.User)
	return users, args.Error(1)
}

func (m *MockUserRepository) InsertAudit(ctx context.Context, db repository.DBTX, userID int64, action string) error {
	args := m.Called(ctx, db, userID, action)
	return args.Error(0)
}

func TestUserService_CreateUser_Success(t *testing.T) {
	ctx := context.Background()

	db := new(MockDB)
	tx := new(MockTx)
	txManager := new(MockTxManager)
	repo := new(MockUserRepository)

	svc := &UserService{
		db:        db,
		txManager: txManager,
		repo:      repo,
	}

	expectedUser := &model.User{
		ID:        1,
		Email:     "alice@example.com",
		Name:      "Alice",
		CreatedAt: time.Now(),
	}

	params := repository.CreateUserParams{
		Email: "alice@example.com",
		Name:  "Alice",
	}

	txManager.
		On("BeginTx", ctx, pgx.TxOptions{IsoLevel: pgx.ReadCommitted}).
		Return(tx, nil).
		Once()

	repo.
		On("Create", ctx, tx, params).
		Return(expectedUser, nil).
		Once()

	repo.
		On("InsertAudit", ctx, tx, int64(1), "user_created").
		Return(nil).
		Once()

	tx.
		On("Commit", ctx).
		Return(nil).
		Once()

	// Deferred rollback still runs after Commit; service ignores its error.
	tx.
		On("Rollback", ctx).
		Return(nil).
		Once()

	user, err := svc.CreateUser(ctx, "alice@example.com", "Alice")

	assert.NoError(t, err)
	assert.Equal(t, expectedUser, user)

	txManager.AssertExpectations(t)
	repo.AssertExpectations(t)
	tx.AssertExpectations(t)
}

func TestUserService_CreateUser_BeginTxError(t *testing.T) {
	ctx := context.Background()

	db := new(MockDB)
	txManager := new(MockTxManager)
	repo := new(MockUserRepository)

	svc := &UserService{
		db:        db,
		txManager: txManager,
		repo:      repo,
	}

	txManager.
		On("BeginTx", ctx, pgx.TxOptions{IsoLevel: pgx.ReadCommitted}).
		Return(nil, assert.AnError).
		Once()

	user, err := svc.CreateUser(ctx, "alice@example.com", "Alice")

	assert.Nil(t, user)
	assert.Error(t, err)
	assert.ErrorContains(t, err, "begin tx")
	assert.ErrorIs(t, err, assert.AnError)

	repo.AssertNotCalled(t, "Create", mock.Anything, mock.Anything, mock.Anything)
	repo.AssertNotCalled(t, "InsertAudit", mock.Anything, mock.Anything, mock.Anything, mock.Anything)
	txManager.AssertExpectations(t)
}

func TestUserService_CreateUser_CreateError(t *testing.T) {
	ctx := context.Background()

	db := new(MockDB)
	tx := new(MockTx)
	txManager := new(MockTxManager)
	repo := new(MockUserRepository)

	svc := &UserService{
		db:        db,
		txManager: txManager,
		repo:      repo,
	}

	params := repository.CreateUserParams{
		Email: "alice@example.com",
		Name:  "Alice",
	}

	txManager.
		On("BeginTx", ctx, pgx.TxOptions{IsoLevel: pgx.ReadCommitted}).
		Return(tx, nil).
		Once()

	repo.
		On("Create", ctx, tx, params).
		Return(nil, assert.AnError).
		Once()

	tx.
		On("Rollback", ctx).
		Return(nil).
		Once()

	user, err := svc.CreateUser(ctx, "alice@example.com", "Alice")

	assert.Nil(t, user)
	assert.Error(t, err)
	assert.ErrorIs(t, err, assert.AnError)

	repo.AssertNotCalled(t, "InsertAudit", mock.Anything, mock.Anything, mock.Anything, mock.Anything)
	tx.AssertNotCalled(t, "Commit", mock.Anything)

	txManager.AssertExpectations(t)
	repo.AssertExpectations(t)
	tx.AssertExpectations(t)
}

func TestUserService_CreateUser_InsertAuditError(t *testing.T) {
	ctx := context.Background()

	db := new(MockDB)
	tx := new(MockTx)
	txManager := new(MockTxManager)
	repo := new(MockUserRepository)

	svc := &UserService{
		db:        db,
		txManager: txManager,
		repo:      repo,
	}

	expectedUser := &model.User{
		ID:    1,
		Email: "alice@example.com",
		Name:  "Alice",
	}

	params := repository.CreateUserParams{
		Email: "alice@example.com",
		Name:  "Alice",
	}

	txManager.
		On("BeginTx", ctx, pgx.TxOptions{IsoLevel: pgx.ReadCommitted}).
		Return(tx, nil).
		Once()

	repo.
		On("Create", ctx, tx, params).
		Return(expectedUser, nil).
		Once()

	repo.
		On("InsertAudit", ctx, tx, int64(1), "user_created").
		Return(assert.AnError).
		Once()

	tx.
		On("Rollback", ctx).
		Return(nil).
		Once()

	user, err := svc.CreateUser(ctx, "alice@example.com", "Alice")

	assert.Nil(t, user)
	assert.Error(t, err)
	assert.ErrorIs(t, err, assert.AnError)

	tx.AssertNotCalled(t, "Commit", mock.Anything)

	txManager.AssertExpectations(t)
	repo.AssertExpectations(t)
	tx.AssertExpectations(t)
}

func TestUserService_CreateUser_CommitRollbackError(t *testing.T) {
	ctx := context.Background()

	db := new(MockDB)
	tx := new(MockTx)
	txManager := new(MockTxManager)
	repo := new(MockUserRepository)

	svc := &UserService{
		db:        db,
		txManager: txManager,
		repo:      repo,
	}

	expectedUser := &model.User{
		ID:    1,
		Email: "alice@example.com",
		Name:  "Alice",
	}

	params := repository.CreateUserParams{
		Email: "alice@example.com",
		Name:  "Alice",
	}

	txManager.
		On("BeginTx", ctx, pgx.TxOptions{IsoLevel: pgx.ReadCommitted}).
		Return(tx, nil).
		Once()

	repo.
		On("Create", ctx, tx, params).
		Return(expectedUser, nil).
		Once()

	repo.
		On("InsertAudit", ctx, tx, int64(1), "user_created").
		Return(nil).
		Once()

	tx.
		On("Commit", ctx).
		Return(pgx.ErrTxCommitRollback).
		Once()

	tx.
		On("Rollback", ctx).
		Return(nil).
		Once()

	user, err := svc.CreateUser(ctx, "alice@example.com", "Alice")

	assert.Nil(t, user)
	assert.Error(t, err)
	assert.ErrorIs(t, err, pgx.ErrTxCommitRollback)
	assert.ErrorContains(t, err, "commit failed and transaction was rolled back")

	txManager.AssertExpectations(t)
	repo.AssertExpectations(t)
	tx.AssertExpectations(t)
}

func TestUserService_GetUser_Success(t *testing.T) {
	ctx := context.Background()

	db := new(MockDB)
	txManager := new(MockTxManager)
	repo := new(MockUserRepository)

	svc := &UserService{
		db:        db,
		txManager: txManager,
		repo:      repo,
	}

	expectedUser := &model.User{
		ID:    7,
		Email: "bob@example.com",
		Name:  "Bob",
	}

	repo.
		On("GetByID", ctx, db, int64(7)).
		Return(expectedUser, nil).
		Once()

	user, err := svc.GetUser(ctx, 7)

	assert.NoError(t, err)
	assert.Equal(t, expectedUser, user)

	repo.AssertExpectations(t)
}

func TestUserService_ListUsers_Success(t *testing.T) {
	ctx := context.Background()

	db := new(MockDB)
	txManager := new(MockTxManager)
	repo := new(MockUserRepository)

	svc := &UserService{
		db:        db,
		txManager: txManager,
		repo:      repo,
	}

	expectedUsers := []model.User{
		{ID: 1, Email: "a@example.com", Name: "A"},
		{ID: 2, Email: "b@example.com", Name: "B"},
	}

	repo.
		On("List", ctx, db, int32(20), int32(0)).
		Return(expectedUsers, nil).
		Once()

	users, err := svc.ListUsers(ctx, 20, 0)

	assert.NoError(t, err)
	assert.Equal(t, expectedUsers, users)

	repo.AssertExpectations(t)
}

func TestUserService_IsNotFound(t *testing.T) {
	db := new(MockDB)
	txManager := new(MockTxManager)
	repo := new(MockUserRepository)

	svc := &UserService{
		db:        db,
		txManager: txManager,
		repo:      repo,
	}

	assert.True(t, svc.IsNotFound(repository.ErrNotFound))
	assert.False(t, svc.IsNotFound(assert.AnError))
}
