package repository

import (
	"cliwallet/internal/model"
	"cliwallet/tests/testutil"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
)

func TestAccount__Create(t *testing.T) {
	t.Parallel()

	db := testutil.SetupTestDB(t, DBSQLite, MigrateAll)
	repository := NewAccountRepo(db)

	t.Run("on success", func(t *testing.T) {
		account := model.Account{
			Name: "unit-test-" + uuid.NewString()[:8],
		}
		created, err := repository.Create(t.Context(), account)

		require.NoError(t, err)
		require.NotEqual(t, uuid.Nil, created.ID)
		require.Equal(t, account.Name, created.Name)
	})
	t.Run("on duplicate id is ignored", func(t *testing.T) {
		account := model.Account{
			Name: "unit-test-" + uuid.NewString()[:8],
		}
		created, err := repository.Create(t.Context(), account)

		require.NoError(t, err)
		require.NotEqual(t, uuid.Nil, created.ID)

		account2 := model.Account{
			// same id
			ID: created.ID,
			// New Name
			Name: "unit-test-" + uuid.NewString()[:8],
		}
		created2, err := repository.Create(t.Context(), account2)

		require.NoError(t, err)
		require.NotEqual(t, uuid.Nil, created2.ID)
		require.Equal(t, account2.Name, created2.Name)
	})
	t.Run("on duplicate name", func(t *testing.T) {
		account := model.Account{
			Name: "unit-test-" + uuid.NewString()[:8],
		}
		created, err := repository.Create(t.Context(), account)

		require.NoError(t, err)
		require.NotEqual(t, uuid.Nil, created.ID)
		require.Equal(t, account.Name, created.Name)

		// Using same account!
		created2, err := repository.Create(t.Context(), account)

		require.Error(t, err)
		require.ErrorIs(t, err, ErrRecordDuplicate)
		require.Equal(t, uuid.Nil, created2.ID)
	})

}
