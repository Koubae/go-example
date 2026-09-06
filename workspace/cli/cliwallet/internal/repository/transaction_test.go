package repository

import (
	"cliwallet/internal/model"
	"cliwallet/tests/testutil"
	"errors"
	"fmt"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func createAccountAndWallet(t *testing.T, accountRepo Account, walletRepo Wallet) (model.Account, model.Wallet) {
	t.Helper()

	_id := uuid.NewString()[:8]
	account, err := accountRepo.Create(t.Context(), model.Account{
		Name: "unit-test-account-" + _id,
	})
	require.NoError(t, err)
	require.NotEqual(t, uuid.Nil, account.ID)

	wallet, err := walletRepo.Create(t.Context(), model.Wallet{
		AccountID: account.ID,
		Name:      "unit-test-wallet-" + _id,
		Currency:  model.EUR,
		Balance:   model.NewMoney(0, model.EUR),
	})
	require.NoError(t, err)
	require.NotEqual(t, uuid.Nil, wallet.ID)

	return account, wallet
}

func TestTransaction__Create(t *testing.T) {
	t.Parallel()

	db := testutil.SetupTestDB(t, DBSQLite, MigrateAll)
	accountRepo := NewAccountRepo(db)
	walletRepo := NewWalletRepo(db)
	repository := NewTransactionRepo(db)
	tx := NewDBTx(db)

	_, wallet := createAccountAndWallet(t, accountRepo, walletRepo)

	t.Run("on success", func(t *testing.T) {
		err := tx.Do(t.Context(), func(tx Tx) error {
			transaction, err := repository.Create(t.Context(), tx.Session(), model.Transaction{
				WalletID:    wallet.ID,
				Type:        model.TransactionDeposit,
				Currency:    model.EUR,
				Amount:      model.NewMoney(100, model.EUR),
				Description: "unit-test-transaction",
			})
			require.NoError(t, err)
			require.NotEqual(t, uuid.Nil, transaction.ID)
			require.Equal(t, model.TransactionDeposit, transaction.Type)
			require.Equal(t, model.EUR, transaction.Currency)
			require.Equal(t, model.NewMoney(100, model.EUR), transaction.Amount)
			require.Equal(t, "unit-test-transaction", transaction.Description)
			require.NotZero(t, transaction.CreatedAt)
			return nil
		})
		require.NoError(t, err)
	})

	t.Run("on injecting duplicate id new one is replaces", func(t *testing.T) {
		var transactionID uuid.UUID

		err := tx.Do(t.Context(), func(tx Tx) error {
			transaction, err := repository.Create(t.Context(), tx.Session(), model.Transaction{
				WalletID:    wallet.ID,
				Type:        model.TransactionDeposit,
				Currency:    model.EUR,
				Amount:      model.NewMoney(3000, model.EUR),
				Description: "unit-test-transaction",
			})
			require.NoError(t, err)
			require.NotEqual(t, uuid.Nil, transaction.ID)
			require.Equal(t, model.TransactionDeposit, transaction.Type)
			require.Equal(t, model.EUR, transaction.Currency)
			require.Equal(t, model.NewMoney(3000, model.EUR), transaction.Amount)
			require.Equal(t, "unit-test-transaction", transaction.Description)
			require.NotZero(t, transaction.CreatedAt)

			transactionID = transaction.ID
			return nil
		})
		require.NoError(t, err)

		require.NotEqual(t, uuid.Nil, transactionID)

		// Inserting once more but injecting the same id
		err = tx.Do(t.Context(), func(tx Tx) error {
			transaction, err := repository.Create(t.Context(), tx.Session(), model.Transaction{
				ID:          transactionID,
				WalletID:    wallet.ID,
				Type:        model.TransactionDeposit,
				Currency:    model.EUR,
				Amount:      model.NewMoney(3000, model.EUR),
				Description: "unit-test-transaction",
			})
			require.NoError(t, err)
			require.NotEqual(t, uuid.Nil, transaction.ID)
			require.Equal(t, model.TransactionDeposit, transaction.Type)
			require.Equal(t, model.EUR, transaction.Currency)
			require.Equal(t, model.NewMoney(3000, model.EUR), transaction.Amount)
			require.Equal(t, "unit-test-transaction", transaction.Description)
			require.NotZero(t, transaction.CreatedAt)

			transactionID = transaction.ID
			return nil
		})
		require.NoError(t, err)
	})

	t.Run("rolled back when Do returns error", func(t *testing.T) {
		var id uuid.UUID
		err := tx.Do(t.Context(), func(tx Tx) error {
			created, err := repository.Create(t.Context(), tx.Session(), model.Transaction{
				WalletID:    wallet.ID,
				Type:        model.TransactionDeposit,
				Currency:    model.EUR,
				Amount:      model.NewMoney(100, model.EUR),
				Description: "unit-test-rollback",
			})
			require.NoError(t, err)
			id = created.ID
			return errors.New("force rollback")
		})
		require.Error(t, err)
		require.NotEqual(t, uuid.Nil, id)
		got, err := repository.GetByID(t.Context(), id)
		require.ErrorIs(t, err, ErrRecordNotFound)
		require.Equal(t, model.Transaction{}, got)
	})

}

func TestTransaction__GetByID(t *testing.T) {
	t.Parallel()

	db := testutil.SetupTestDB(t, DBSQLite, MigrateAll)
	accountRepo := NewAccountRepo(db)
	walletRepo := NewWalletRepo(db)
	repository := NewTransactionRepo(db)
	tx := NewDBTx(db)

	_, wallet := createAccountAndWallet(t, accountRepo, walletRepo)

	var transaction model.Transaction
	err := tx.Do(t.Context(), func(tx Tx) error {
		created, err := repository.Create(t.Context(), tx.Session(), model.Transaction{
			WalletID:    wallet.ID,
			Type:        model.TransactionDeposit,
			Currency:    model.EUR,
			Amount:      model.NewMoney(100, model.EUR),
			Description: "unit-test-transaction",
		})
		transaction = created

		require.NoError(t, err)
		require.NotEqual(t, uuid.Nil, transaction.ID)
		require.Equal(t, model.TransactionDeposit, transaction.Type)
		require.Equal(t, model.EUR, transaction.Currency)
		require.Equal(t, model.NewMoney(100, model.EUR), transaction.Amount)
		require.Equal(t, "unit-test-transaction", transaction.Description)
		require.NotZero(t, transaction.CreatedAt)
		return nil
	})
	require.NoError(t, err)
	require.NotEqual(t, uuid.Nil, transaction.ID)

	t.Run("on success", func(t *testing.T) {
		record, err := repository.GetByID(t.Context(), transaction.ID)
		require.NoError(t, err)
		assert.Equal(t, transaction, record)
	})

	t.Run("on not found", func(t *testing.T) {
		record, err := repository.GetByID(t.Context(), uuid.New())
		require.ErrorIs(t, err, ErrRecordNotFound)
		require.Equal(t, model.Transaction{}, record)
	})
}

func TestTransaction__ListByWalletID(t *testing.T) {
	t.Parallel()

	db := testutil.SetupTestDB(t, DBSQLite, MigrateAll)
	accountRepo := NewAccountRepo(db)
	walletRepo := NewWalletRepo(db)
	repository := NewTransactionRepo(db)
	tx := NewDBTx(db)

	_, wallet := createAccountAndWallet(t, accountRepo, walletRepo)

	deposits := [10]int64{100, 200, 300, 400, 500, 600, 700, 800, 900, 1000}
	transactionIDs := make(map[uuid.UUID]string, len(deposits))
	err := tx.Do(t.Context(), func(tx Tx) error {
		for i, amount := range deposits {
			created, err := repository.Create(t.Context(), tx.Session(), model.Transaction{
				WalletID:    wallet.ID,
				Type:        model.TransactionDeposit,
				Currency:    model.EUR,
				Amount:      model.NewMoney(amount, model.EUR),
				Description: fmt.Sprintf("unit-test-transaction-%d", i),
			})
			require.NoError(t, err)
			transactionIDs[created.ID] = created.Description
		}
		return nil

	})
	require.NoError(t, err)

	t.Run("on success", func(t *testing.T) {
		transactions, err := repository.ListByWalletID(t.Context(), wallet.ID, 10, 0)

		require.NoError(t, err)
		assert.Len(t, transactions, len(deposits))
		for _, transaction := range transactions {
			description, ok := transactionIDs[transaction.ID]

			assert.True(t, ok)
			assert.Equal(t, description, transaction.Description)
		}
	})

	t.Run("on not found", func(t *testing.T) {
		transactions, err := repository.ListByWalletID(t.Context(), uuid.New(), 10, 0)
		require.NoError(t, err)
		assert.Empty(t, transactions)
	})

	t.Run("limit & offset respected", func(t *testing.T) {

		seen := make(map[uuid.UUID]struct{}, len(deposits))
		limit := 1
		offset := 0
		for i := range len(deposits) {
			records, err := repository.ListByWalletID(t.Context(), wallet.ID, limit, offset)

			require.NoError(t, err)
			require.Len(t, records, 1, fmt.Sprintf("unexpected records found in loop %d", i+1))

			_, duplicate := seen[records[0].ID]
			require.False(t, duplicate)
			seen[records[0].ID] = struct{}{}

			offset++
		}
		assert.Len(t, seen, len(deposits))

		// offset now is len(wallets) + 1 so cursor must go beyond actual records
		records, err := repository.ListByWalletID(t.Context(), wallet.ID, limit, offset)

		require.NoError(t, err)
		assert.Equal(t, 0, len(records))
		assert.Equal(t, []model.Transaction{}, records)

	})

	t.Run("limit & offset validation", func(t *testing.T) {
		record, err := repository.ListByWalletID(t.Context(), wallet.ID, 0, 0)
		require.Nil(t, record)
		require.Error(t, err, ErrQueryInvalidLimit)

		record, err = repository.ListByWalletID(t.Context(), wallet.ID, -1, 0)
		require.Nil(t, record)
		require.ErrorIs(t, err, ErrQueryInvalidLimit)

		record, err = repository.ListByWalletID(t.Context(), wallet.ID, 1, -1)
		require.Nil(t, record)
		require.ErrorIs(t, err, ErrQueryInvalidOffSet)

	})

}
