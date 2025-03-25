package db

import (
	"context"
	"testing"

	"github.com/kimseogyu/portfolio/backend/internal/testutils"
	"github.com/stretchr/testify/assert"
)

func TestNewDB(t *testing.T) {
	pg, err := testutils.NewPostgresTestContainer(context.Background())
	assert.NoError(t, err)

	t.Cleanup(func() {
		err = pg.Close()
		assert.NoError(t, err)
	})

	endpoint, err := pg.Endpoint()
	assert.NoError(t, err)

	db, err := NewDB(WithDBType(DBTypePostgres), WithPostgresOptions(WithDSN(endpoint)))
	assert.NoError(t, err)

	t.Cleanup(func() {
		sqlDB, err := db.DB()
		assert.NoError(t, err)

		err = sqlDB.Close()
		assert.NoError(t, err)
	})
}
