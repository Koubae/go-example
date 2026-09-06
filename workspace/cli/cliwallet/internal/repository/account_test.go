package repository

import (
	"cliwallet/internal/model"
	"cliwallet/tests/testutil"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
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

		// Using same record!
		created2, err := repository.Create(t.Context(), account)

		require.Error(t, err)
		require.ErrorIs(t, err, ErrRecordDuplicate)
		require.Equal(t, uuid.Nil, created2.ID)
	})

}

func TestAccount__GetByID(t *testing.T) {
	t.Parallel()

	db := testutil.SetupTestDB(t, DBSQLite, MigrateAll)
	repository := NewAccountRepo(db)

	t.Run("on success", func(t *testing.T) {
		created, err := repository.Create(t.Context(), model.Account{
			Name: "unit-test-" + uuid.NewString()[:8],
		})
		require.NoError(t, err)
		require.NotEqual(t, uuid.Nil, created.ID)

		model, err := repository.GetByID(t.Context(), created.ID)

		require.NoError(t, err)
		require.Equal(t, created, model)

	})
	t.Run("on not found", func(t *testing.T) {
		record, err := repository.GetByID(t.Context(), uuid.New())

		require.Error(t, err)
		require.ErrorIs(t, err, ErrRecordNotFound)
		assert.Equal(t, uuid.Nil, record.ID)
		require.Equal(t, model.Account{}, record)
	})
	t.Run("returns the matching account", func(t *testing.T) {
		a, err := repository.Create(t.Context(), model.Account{Name: "unit-test-" + uuid.NewString()[:8]})
		require.NoError(t, err)
		b, err := repository.Create(t.Context(), model.Account{Name: "unit-test-" + uuid.NewString()[:8]})
		require.NoError(t, err)

		recordA, err := repository.GetByID(t.Context(), a.ID)

		require.NoError(t, err)
		require.Equal(t, a, recordA)
		require.NotEqual(t, b.ID, recordA.ID)

		recordB, err := repository.GetByID(t.Context(), b.ID)

		require.NoError(t, err)
		require.Equal(t, b, recordB)
		require.NotEqual(t, a.ID, recordB.ID)
	})
}
