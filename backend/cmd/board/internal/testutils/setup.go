package testutils

import (
	"os"
	"testing"

	"github.com/kimseogyu/portfolio/backend/internal/cstore"
	"github.com/kimseogyu/portfolio/backend/internal/db"
	"github.com/kimseogyu/portfolio/backend/internal/redisutils"
	"github.com/kimseogyu/portfolio/backend/internal/testutils"
	"github.com/stretchr/testify/require"
	"github.com/testcontainers/testcontainers-go"
	"gorm.io/gorm"
)

type Fixture struct {
	Conn           *gorm.DB
	DbContainer    *testutils.PostgresTestContainer
	CacheStore     cstore.CacheStore
	RedisContainer testcontainers.Container
}

func SetupFixture(t *testing.T) (*Fixture, error) {
	pg, err := testutils.NewPostgresTestContainer(t.Context())
	require.NoError(t, err)

	endpoint, err := pg.Endpoint()
	require.NoError(t, err)

	db, err := db.NewPostgresDB(db.WithDSN(endpoint))
	require.NoError(t, err)

	sqlPath := "../../../../scripts/init.sql"
	content, err := os.ReadFile(sqlPath)
	require.NoError(t, err)

	db.Exec(string(content))

	// 테이블 truncate로 깨끗한 상태 유지
	TruncateTables(t, db)

	redisContainer, err := testutils.NewRedisTestContainer(t.Context())
	require.NoError(t, err)

	redisEndpoint, err := redisContainer.Endpoint(t.Context(), "")
	require.NoError(t, err)

	redisClient, err := redisutils.NewRedisClient(t.Context(), redisEndpoint)
	require.NoError(t, err)

	cacheStore := cstore.NewCacheStore(redisClient)

	return &Fixture{
		Conn:           db,
		DbContainer:    pg,
		CacheStore:     cacheStore,
		RedisContainer: redisContainer,
	}, nil
}

// 테이블 데이터 정리 함수
func TruncateTables(t *testing.T, db *gorm.DB) {
	// 외래 키 제약 조건 일시적으로 비활성화
	err := db.Exec("SET CONSTRAINTS ALL DEFERRED").Error
	require.NoError(t, err)

	// 테이블 truncate - 댓글 먼저 삭제 (외래 키 의존성 때문)
	err = db.Exec("TRUNCATE TABLE comments RESTART IDENTITY CASCADE").Error
	require.NoError(t, err)

	err = db.Exec("TRUNCATE TABLE postings RESTART IDENTITY CASCADE").Error
	require.NoError(t, err)

	// 외래 키 제약 조건 다시 활성화
	err = db.Exec("SET CONSTRAINTS ALL IMMEDIATE").Error
	require.NoError(t, err)
}
