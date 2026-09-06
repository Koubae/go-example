package repository

import (
	"cliwallet/internal/model"
	"cliwallet/tests/testutil"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestWallet__Create(t *testing.T) {
	t.Parallel()

	db := testutil.SetupTestDB(t, DBSQLite, MigrateAll)
	accountRepo := NewAccountRepo(db)
	repository := NewWalletRepo(db)

	t.Run("on success", func(t *testing.T) {
		account, err := accountRepo.Create(t.Context(), model.Account{
			Name: "unit-test-account-" + uuid.NewString()[:8],
		})
		require.NoError(t, err)

		wallet := model.Wallet{
			AccountID: account.ID,
			Name:      "unit-test-wallet-" + uuid.NewString()[:8],
			Currency:  model.EUR,
			Balance:   model.NewMoney(1050, model.EUR),
		}

		created, err := repository.Create(t.Context(), wallet)

		require.NoError(t, err)
		require.NotEqual(t, uuid.Nil, created.ID)
		require.Equal(t, wallet.Name, created.Name)
		assert.Equal(t, int64(1050), created.Balance.Amount)
	})
	t.Run("account owns multiple wallets", func(t *testing.T) {
		account, err := accountRepo.Create(t.Context(), model.Account{
			Name: "unit-test-account-" + uuid.NewString()[:8],
		})
		require.NoError(t, err)

		_id := uuid.NewString()[:8]

		wallets := [5]model.Wallet{
			{Name: "multi-wallet-1-" + _id, Currency: model.EUR, Balance: model.NewMoney(1050, model.EUR)},
			{Name: "multi-wallet-2-" + _id, Currency: model.EUR, Balance: model.NewMoney(2500, model.EUR)},
			{Name: "multi-wallet-3-" + _id, Currency: model.EUR, Balance: model.NewMoney(600000, model.USD)},
			{Name: "multi-wallet-4-" + _id, Currency: model.EUR, Balance: model.NewMoney(300000, model.JPY)},
			{Name: "multi-wallet-5-" + _id, Currency: model.EUR, Balance: model.NewMoney(45000, model.EUR)},
		}
		for _, wallet := range wallets {
			wallet.AccountID = account.ID

			created, err := repository.Create(t.Context(), wallet)
			require.NoError(t, err)
			require.NotEqual(t, uuid.Nil, created.ID)
			require.Equal(t, wallet.Name, created.Name)
			assert.Equal(t, wallet.Balance.Amount, created.Balance.Amount)
			assert.Equal(t, wallet.Currency, created.Currency)
		}

	})
	t.Run("on duplicate id is ignored", func(t *testing.T) {
		account, err := accountRepo.Create(t.Context(), model.Account{
			Name: "unit-test-account-" + uuid.NewString()[:8],
		})
		require.NoError(t, err)

		for range 2 {
			wallet := model.Wallet{
				AccountID: account.ID,
				Name:      "unit-test-wallet-" + uuid.NewString()[:8],
				Currency:  model.EUR,
				Balance:   model.NewMoney(1050, model.EUR),
			}

			created, err := repository.Create(t.Context(), wallet)

			require.NoError(t, err)
			require.NotEqual(t, uuid.Nil, created.ID)
			require.Equal(t, wallet.Name, created.Name)
			assert.Equal(t, int64(1050), created.Balance.Amount)
		}

	})
	t.Run("on duplicate name", func(t *testing.T) {
		account, err := accountRepo.Create(t.Context(), model.Account{
			Name: "unit-test-account-" + uuid.NewString()[:8],
		})
		require.NoError(t, err)

		wallet := model.Wallet{
			AccountID: account.ID,
			Name:      "unit-test-wallet-" + uuid.NewString()[:8],
			Currency:  model.EUR,
			Balance:   model.NewMoney(1050, model.EUR),
		}

		created, err := repository.Create(t.Context(), wallet)

		require.NoError(t, err)
		require.NotEqual(t, uuid.Nil, created.ID)

		// Using same record!
		created2, err := repository.Create(t.Context(), wallet)

		require.Error(t, err)
		require.ErrorIs(t, err, ErrRecordDuplicate)
		require.Equal(t, uuid.Nil, created2.ID)

	})

}
func TestWallet__GetByID(t *testing.T) {
	t.Parallel()

	db := testutil.SetupTestDB(t, DBSQLite, MigrateAll)
	repository := NewWalletRepo(db)

	t.Run("on success", func(t *testing.T) {
		println(repository)
	})

}
func TestWallet__GetByAccountIDAndName(t *testing.T) {
	t.Parallel()

	db := testutil.SetupTestDB(t, DBSQLite, MigrateAll)
	repository := NewWalletRepo(db)

	t.Run("on success", func(t *testing.T) {
		println(repository)
	})

}
func TestWallet__ListByAccountID(t *testing.T) {
	t.Parallel()

	db := testutil.SetupTestDB(t, DBSQLite, MigrateAll)
	repository := NewWalletRepo(db)

	t.Run("on success", func(t *testing.T) {
		println(repository)
	})

}
