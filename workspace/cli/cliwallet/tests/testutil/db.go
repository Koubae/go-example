package testutil

import (
	"testing"

	"github.com/stretchr/testify/require"
	"gorm.io/gorm"
)

type connector func() (*gorm.DB, error)
type migrator func(*gorm.DB) error

func SetupTestDB(t *testing.T, c connector, m migrator) *gorm.DB {
	t.Helper()

	db, err := c()
	require.NoError(t, err)

	require.NoError(t, m(db))

	sqlDB, err := db.DB()
	require.NoError(t, err)

	t.Cleanup(func() {
		_ = sqlDB.Close()
	})

	return db

}
