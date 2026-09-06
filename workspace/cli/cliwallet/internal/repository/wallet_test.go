package repository

import (
	"cliwallet/internal/model"
	"cliwallet/tests/testutil"
	"fmt"
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
	t.Run("on injecting duplicate id new one is replaces", func(t *testing.T) {
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
	accountRepo := NewAccountRepo(db)
	repository := NewWalletRepo(db)

	account, err := accountRepo.Create(t.Context(), model.Account{
		Name: "unit-test-account-" + uuid.NewString()[:8],
	})
	require.NoError(t, err)

	t.Run("on success", func(t *testing.T) {
		created, err := repository.Create(t.Context(), model.Wallet{
			AccountID: account.ID,
			Name:      "unit-test-wallet-" + uuid.NewString()[:8],
			Currency:  model.EUR,
			Balance:   model.NewMoney(1050, model.EUR),
		})

		require.NoError(t, err)
		require.NotEqual(t, uuid.Nil, created.ID)

		model, err := repository.GetByID(t.Context(), created.ID)

		require.NoError(t, err)
		require.Equal(t, created, model)
	})
	t.Run("on not found", func(t *testing.T) {
		record, err := repository.GetByID(t.Context(), uuid.New())

		require.ErrorIs(t, err, ErrRecordNotFound)
		assert.Equal(t, uuid.Nil, record.ID)
		require.Equal(t, model.Wallet{}, record)
	})
	t.Run("get matching account", func(t *testing.T) {
		a, err := accountRepo.Create(t.Context(), model.Account{Name: "unit-test-" + uuid.NewString()[:8]})
		require.NoError(t, err)
		require.NotEqual(t, uuid.Nil, a.ID)
		b, err := accountRepo.Create(t.Context(), model.Account{Name: "unit-test-" + uuid.NewString()[:8]})
		require.NoError(t, err)
		require.NotEqual(t, uuid.Nil, b.ID)

		walletA, err := repository.Create(t.Context(), model.Wallet{
			AccountID: a.ID,
			Name:      "unit-test-wallet-" + uuid.NewString()[:8],
			Currency:  model.EUR,
			Balance:   model.NewMoney(1050, model.EUR),
		})
		require.NoError(t, err)
		require.NotEqual(t, uuid.Nil, walletA.ID)

		walletB, err := repository.Create(t.Context(), model.Wallet{
			AccountID: b.ID,
			Name:      "unit-test-wallet-" + uuid.NewString()[:8],
			Currency:  model.USD,
			Balance:   model.NewMoney(55500, model.EUR),
		})
		require.NoError(t, err)
		require.NotEqual(t, uuid.Nil, walletB.ID)

		modelA, err := repository.GetByID(t.Context(), walletA.ID)
		require.NoError(t, err)
		modelB, err := repository.GetByID(t.Context(), walletB.ID)
		require.NoError(t, err)

		require.Equal(t, a.ID, modelA.AccountID)
		require.Equal(t, model.EUR, modelA.Currency)

		require.Equal(t, b.ID, modelB.AccountID)
		require.Equal(t, model.USD, modelB.Currency)

	})

	t.Run("get matcing wallet withing single account", func(t *testing.T) {
		_id := uuid.NewString()[:8]

		wallets := [5]model.Wallet{
			{Name: "multi-wallet-1-" + _id, Currency: model.EUR, Balance: model.NewMoney(1050, model.EUR)},
			{Name: "multi-wallet-2-" + _id, Currency: model.EUR, Balance: model.NewMoney(2500, model.EUR)},
			{Name: "multi-wallet-3-" + _id, Currency: model.EUR, Balance: model.NewMoney(600000, model.USD)},
			{Name: "multi-wallet-4-" + _id, Currency: model.EUR, Balance: model.NewMoney(300000, model.JPY)},
			{Name: "multi-wallet-5-" + _id, Currency: model.EUR, Balance: model.NewMoney(45000, model.EUR)},
		}
		walletIDs := make([]uuid.UUID, len(wallets))
		for i, wallet := range wallets {
			wallet.AccountID = account.ID

			created, err := repository.Create(t.Context(), wallet)
			require.NoError(t, err)
			require.NotEqual(t, uuid.Nil, created.ID)
			require.Equal(t, wallet.Name, created.Name)
			assert.Equal(t, wallet.Balance.Amount, created.Balance.Amount)
			assert.Equal(t, wallet.Currency, created.Currency)

			walletIDs[i] = created.ID
		}

		for _, id := range walletIDs {
			model, err := repository.GetByID(t.Context(), id)
			require.NoError(t, err)
			assert.Equal(t, account.ID, model.AccountID)
			assert.Equal(t, id, model.ID)
		}

	})

}
func TestWallet__GetByAccountIDAndName(t *testing.T) {
	t.Parallel()

	db := testutil.SetupTestDB(t, DBSQLite, MigrateAll)
	accountRepo := NewAccountRepo(db)
	repository := NewWalletRepo(db)

	account, err := accountRepo.Create(t.Context(), model.Account{
		Name: "unit-test-account-" + uuid.NewString()[:8],
	})
	require.NoError(t, err)

	t.Run("on success", func(t *testing.T) {
		created, err := repository.Create(t.Context(), model.Wallet{
			AccountID: account.ID,
			Name:      "unit-test-wallet-" + uuid.NewString()[:8],
			Currency:  model.EUR,
			Balance:   model.NewMoney(1050, model.EUR),
		})

		require.NoError(t, err)
		require.NotEqual(t, uuid.Nil, created.ID)

		model, err := repository.GetByAccountIDAndName(t.Context(), account.ID, created.Name)

		require.NoError(t, err)
		require.Equal(t, created, model)
	})
	t.Run("on not found", func(t *testing.T) {
		record, err := repository.GetByAccountIDAndName(t.Context(), account.ID, "wallet-does-not-exists")

		require.Error(t, err)
		require.ErrorIs(t, err, ErrRecordNotFound)
		assert.Equal(t, uuid.Nil, record.ID)
		require.Equal(t, model.Wallet{}, record)
	})
	t.Run("get matching account", func(t *testing.T) {
		a, err := accountRepo.Create(t.Context(), model.Account{Name: "unit-test-" + uuid.NewString()[:8]})
		require.NoError(t, err)
		require.NotEqual(t, uuid.Nil, a.ID)
		b, err := accountRepo.Create(t.Context(), model.Account{Name: "unit-test-" + uuid.NewString()[:8]})
		require.NoError(t, err)
		require.NotEqual(t, uuid.Nil, b.ID)

		walletA, err := repository.Create(t.Context(), model.Wallet{
			AccountID: a.ID,
			Name:      "unit-test-wallet-" + uuid.NewString()[:8],
			Currency:  model.EUR,
			Balance:   model.NewMoney(1050, model.EUR),
		})
		require.NoError(t, err)
		require.NotEqual(t, uuid.Nil, walletA.ID)

		walletB, err := repository.Create(t.Context(), model.Wallet{
			AccountID: b.ID,
			Name:      "unit-test-wallet-" + uuid.NewString()[:8],
			Currency:  model.USD,
			Balance:   model.NewMoney(55500, model.EUR),
		})
		require.NoError(t, err)
		require.NotEqual(t, uuid.Nil, walletB.ID)

		modelA, err := repository.GetByAccountIDAndName(t.Context(), a.ID, walletA.Name)
		require.NoError(t, err)
		modelB, err := repository.GetByAccountIDAndName(t.Context(), b.ID, walletB.Name)
		require.NoError(t, err)

		require.Equal(t, a.ID, modelA.AccountID)
		require.Equal(t, walletA.Name, modelA.Name)
		require.Equal(t, model.EUR, modelA.Currency)

		require.Equal(t, b.ID, modelB.AccountID)
		require.Equal(t, walletB.Name, modelB.Name)
		require.Equal(t, model.USD, modelB.Currency)

	})

	t.Run("get matcing wallet withing single account", func(t *testing.T) {
		_id := uuid.NewString()[:8]

		wallets := [5]model.Wallet{
			{Name: "multi-wallet-1-" + _id, Currency: model.EUR, Balance: model.NewMoney(1050, model.EUR)},
			{Name: "multi-wallet-2-" + _id, Currency: model.EUR, Balance: model.NewMoney(2500, model.EUR)},
			{Name: "multi-wallet-3-" + _id, Currency: model.EUR, Balance: model.NewMoney(600000, model.USD)},
			{Name: "multi-wallet-4-" + _id, Currency: model.EUR, Balance: model.NewMoney(300000, model.JPY)},
			{Name: "multi-wallet-5-" + _id, Currency: model.EUR, Balance: model.NewMoney(45000, model.EUR)},
		}
		walletNames := make([]string, len(wallets))
		for i, wallet := range wallets {
			wallet.AccountID = account.ID

			created, err := repository.Create(t.Context(), wallet)
			require.NoError(t, err)
			require.NotEqual(t, uuid.Nil, created.ID)
			require.Equal(t, wallet.Name, created.Name)
			assert.Equal(t, wallet.Balance.Amount, created.Balance.Amount)
			assert.Equal(t, wallet.Currency, created.Currency)

			walletNames[i] = created.Name
		}

		for _, name := range walletNames {
			model, err := repository.GetByAccountIDAndName(t.Context(), account.ID, name)
			require.NoError(t, err)
			assert.Equal(t, account.ID, model.AccountID)
			assert.Equal(t, name, model.Name)
		}

	})

}
func TestWallet__ListByAccountID(t *testing.T) {
	t.Parallel()

	db := testutil.SetupTestDB(t, DBSQLite, MigrateAll)
	accountRepo := NewAccountRepo(db)
	repository := NewWalletRepo(db)

	_id := uuid.NewString()[:8]
	account, err := accountRepo.Create(t.Context(), model.Account{
		Name: "unit-test-account-" + _id,
	})
	require.NoError(t, err)

	wallets := [5]model.Wallet{
		{Name: "multi-wallet-1-" + _id, Currency: model.EUR, Balance: model.NewMoney(1050, model.EUR)},
		{Name: "multi-wallet-2-" + _id, Currency: model.EUR, Balance: model.NewMoney(2500, model.EUR)},
		{Name: "multi-wallet-3-" + _id, Currency: model.EUR, Balance: model.NewMoney(600000, model.USD)},
		{Name: "multi-wallet-4-" + _id, Currency: model.EUR, Balance: model.NewMoney(300000, model.JPY)},
		{Name: "multi-wallet-5-" + _id, Currency: model.EUR, Balance: model.NewMoney(45000, model.EUR)},
	}
	idToName := make(map[uuid.UUID]string, len(wallets))
	for _, wallet := range wallets {
		wallet.AccountID = account.ID

		created, err := repository.Create(t.Context(), wallet)
		require.NoError(t, err)
		require.NotEqual(t, uuid.Nil, created.ID)

		idToName[created.ID] = wallet.Name
	}

	t.Run("on success", func(t *testing.T) {
		records, err := repository.ListByAccountID(t.Context(), account.ID, len(wallets), 0)

		require.NoError(t, err)
		assert.Equal(t, len(wallets), len(records))
		for _, record := range records {
			expected, ok := idToName[record.ID]
			require.True(t, ok)

			assert.Equal(t, expected, record.Name)
		}
	})
	t.Run("on not found", func(t *testing.T) {
		records, err := repository.ListByAccountID(t.Context(), uuid.New(), 10, 0)

		require.NoError(t, err)
		assert.Equal(t, []model.Wallet{}, records)
	})
	t.Run("limit & offset respected", func(t *testing.T) {

		seen := make(map[uuid.UUID]struct{}, len(wallets))
		limit := 1
		offset := 0
		for i := range len(wallets) {
			records, err := repository.ListByAccountID(t.Context(), account.ID, limit, offset)

			require.NoError(t, err)
			require.Len(t, records, 1, fmt.Sprintf("unexpected records found in loop %d", i+1))

			_, duplicate := seen[records[0].ID]
			require.False(t, duplicate)
			seen[records[0].ID] = struct{}{}

			offset++
		}
		assert.Len(t, seen, len(wallets))

		// offset now is len(wallets) + 1 so cursor must go beyond actual records
		records, err := repository.ListByAccountID(t.Context(), account.ID, limit, offset)

		require.NoError(t, err)
		assert.Equal(t, 0, len(records))
		assert.Equal(t, []model.Wallet{}, records)

	})

	t.Run("limit & offset validation", func(t *testing.T) {
		record, err := repository.ListByAccountID(t.Context(), account.ID, 0, 0)
		require.Nil(t, record)
		require.Error(t, err, ErrQueryInvalidLimit)

		record, err = repository.ListByAccountID(t.Context(), account.ID, -1, 0)
		require.Nil(t, record)
		require.ErrorIs(t, err, ErrQueryInvalidLimit)

		record, err = repository.ListByAccountID(t.Context(), account.ID, 1, -1)
		require.Nil(t, record)
		require.ErrorIs(t, err, ErrQueryInvalidOffSet)

	})

}
