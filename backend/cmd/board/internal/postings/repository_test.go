package postings

import (
	"context"
	"testing"

	"github.com/kimseogyu/portfolio/backend/cmd/board/internal/testutils"
	boardServer "github.com/kimseogyu/portfolio/backend/internal/proto/board/v1"
	"github.com/stretchr/testify/require"
)

func TestRepository(t *testing.T) {
	fixture, err := testutils.SetupFixture(t)
	require.NoError(t, err)
	t.Cleanup(func() {
		require.NoError(t, fixture.DbContainer.Close())
		require.NoError(t, fixture.RedisContainer.Terminate(context.Background()))
	})

	repo := NewRepository(fixture.Conn, fixture.CacheStore)

	t.Run("CreatePosting", func(t *testing.T) {
		posting := &Posting{
			Title:  "Test Posting",
			Status: boardServer.PostingStatus_POSTING_STATUS_PUBLISHED,
		}

		err := repo.Save(context.Background(), posting)
		require.NoError(t, err)

		require.NotNil(t, posting.ID)
		require.Equal(t, "Test Posting", posting.Title)

		// check ID is not 0
		require.NotEqual(t, 0, posting.ID)

		// check ViewCount is 0
		require.Equal(t, 0, posting.ViewCount)

		// check CommentCount is 0
		require.Equal(t, 0, posting.CommentCount)

		// check Status is 1
		require.Equal(t, boardServer.PostingStatus_POSTING_STATUS_PUBLISHED, posting.Status)

		// check CreatedAt is not 0
		require.NotEqual(t, 0, posting.CreatedAt)

		// check UpdatedAt is not 0
		require.NotEqual(t, 0, posting.UpdatedAt)
	})
}
