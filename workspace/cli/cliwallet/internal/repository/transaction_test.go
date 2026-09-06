package repository

import (
	"cliwallet/internal/model"
	"cliwallet/tests/testutil"
	"errors"
	"testing"

	"github.com/google/uuid"
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
