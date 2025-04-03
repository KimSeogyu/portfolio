package comments

import (
	"context"
	"fmt"
	"math/rand/v2"
	"testing"
	"time"

	"github.com/kimseogyu/portfolio/backend/cmd/board/internal/postings"
	"github.com/kimseogyu/portfolio/backend/cmd/board/internal/testutils"
	boardServer "github.com/kimseogyu/portfolio/backend/internal/proto/board/v1"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/gorm"
)

// 테스트용 Helper 함수들
func createTestPosting(t *testing.T, db *gorm.DB) int64 {
	posting := &postings.Posting{
		Title:      "Test Posting for Comments",
		Content:    "This is a test posting content",
		AuthorID:   "test-author",
		AuthorName: "Test Author",
		Status:     boardServer.PostingStatus_POSTING_STATUS_PUBLISHED,
		CreatedAt:  time.Now(),
		UpdatedAt:  time.Now(),
	}

	err := db.Create(posting).Error
	require.NoError(t, err)

	return int64(posting.ID)
}

// DB에서 직접 댓글 수 확인
func getCommentCountFromDB(t *testing.T, db *gorm.DB, postID int64) int {
	var posting postings.Posting
	err := db.Where("id = ?", postID).First(&posting).Error
	require.NoError(t, err)
	return posting.CommentCount
}

// 디버깅을 위한 댓글 상태 출력
func printCommentsState(t *testing.T, db *gorm.DB, postID int64) {
	var comments []Comment
	db.Where("post_id = ?", postID).Find(&comments)

	t.Logf("댓글 상태 (게시물 ID: %d)", postID)
	for _, c := range comments {
		parentInfo := "최상위 댓글"
		if c.ParentID != nil {
			parentInfo = fmt.Sprintf("대댓글 (부모: %d)", *c.ParentID)
		}

		deleteInfo := "활성"
		if c.DeletedAt != nil {
			deleteInfo = "삭제됨"
		}

		t.Logf("ID: %d, %s, %s, 내용: %s", c.ID, parentInfo, deleteInfo, c.Content)
	}

	var count int
	db.Model(&postings.Posting{}).Where("id = ?", postID).Select("comment_count").Scan(&count)
	t.Logf("현재 DB의 comment_count: %d", count)
}

func TestRepository(t *testing.T) {
	fixture, err := testutils.SetupFixture(t)
	require.NoError(t, err)
	t.Cleanup(func() {
		require.NoError(t, fixture.DbContainer.Close())
		require.NoError(t, fixture.RedisContainer.Terminate(context.Background()))
	})

	repo := NewRepository(fixture.Conn, fixture.CacheStore)
	ctx := context.Background()

	t.Run("Save and GetByID", func(t *testing.T) {
		testutils.TruncateTables(t, fixture.Conn)
		comment := &Comment{
			PostID:    1,
			Content:   "Test comment",
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}

		err := repo.Save(ctx, comment)
		assert.NoError(t, err)
		assert.NotZero(t, comment.ID)

		retrieved, err := repo.GetByPostID(ctx, comment.PostID)
		assert.NoError(t, err)

		assert.Equal(t, comment.Content, retrieved[0].Content)
	})

	t.Run("Save with ParentID", func(t *testing.T) {
		testutils.TruncateTables(t, fixture.Conn)
		parent := &Comment{
			PostID:    1,
			Content:   "Parent comment",
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}

		err := repo.Save(ctx, parent)
		assert.NoError(t, err)
		assert.NotZero(t, parent.ID)

		child := &Comment{
			PostID:    1,
			Content:   "Child comment",
			ParentID:  &parent.ID,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}

		err = repo.Save(ctx, child)
		assert.NoError(t, err)
		assert.NotZero(t, child.ID)
		assert.Equal(t, parent.ID, *child.ParentID)
	})

	t.Run("GetByPostID", func(t *testing.T) {
		testutils.TruncateTables(t, fixture.Conn)
		// Create multiple comments for the same post
		post2Comments := []*Comment{
			{PostID: 2, Content: "First comment"},
			{PostID: 2, Content: "Second comment"},
		}

		for _, c := range post2Comments {
			err := repo.Save(ctx, c)
			require.NoError(t, err)
		}

		comments, err := repo.GetByPostID(ctx, 2)
		assert.NoError(t, err)
		assert.Len(t, comments, 2)
	})

	t.Run("Delete", func(t *testing.T) {
		comment := &Comment{
			PostID:  4,
			Content: "To be deleted",
		}

		err := repo.Save(ctx, comment)
		require.NoError(t, err)

		err = repo.Delete(ctx, comment.ID)
		assert.NoError(t, err)

		comments, err := repo.GetByPostID(ctx, comment.PostID)
		assert.NoError(t, err)
		assert.Len(t, comments, 0)
	})

	t.Run("GetThread", func(t *testing.T) {

	})

	t.Run("Delete when parent has children", func(t *testing.T) {
		// Create a parent comment
		parent := &Comment{
			PostID:  6,
			Content: "Parent to delete",
		}
		err := repo.Save(ctx, parent)
		require.NoError(t, err)

		// Create child comment
		child := &Comment{
			PostID:   6,
			ParentID: &parent.ID,
			Content:  "Child to delete",
		}
		err = repo.Save(ctx, child)
		require.NoError(t, err)

		// Delete parent
		err = repo.Delete(ctx, parent.ID)
		assert.NoError(t, err)

		// Verify parent is deleted and child is not deleted
		comments, err := repo.GetByPostID(ctx, parent.PostID)
		assert.NoError(t, err)
		assert.Len(t, comments, 1)
		assert.Equal(t, child.ID, comments[0].ID)
	})

	t.Run("GetThread", func(t *testing.T) {
		// 테스트 데이터 생성: 최상위 댓글 3개
		postID := rand.Int64() / 3
		rootComments := []*Comment{
			{PostID: postID, Content: "Root 1"},
			{PostID: postID, Content: "Root 2"},
			{PostID: postID, Content: "Root 3"},
		}

		for _, c := range rootComments {
			err := repo.Save(ctx, c)
			require.NoError(t, err)
		}

		// 각 최상위 댓글에 여러 개의 자식 댓글 추가
		for _, root := range rootComments {
			children := []*Comment{
				{PostID: postID, ParentID: &root.ID, Content: "Child 1 of " + root.Content},
				{PostID: postID, ParentID: &root.ID, Content: "Child 2 of " + root.Content},
				{PostID: postID, ParentID: &root.ID, Content: "Child 3 of " + root.Content},
			}
			for _, child := range children {
				err := repo.Save(ctx, child)
				require.NoError(t, err)
			}
		}

		comments, err := repo.GetThread(ctx, postID)
		require.NoError(t, err)
		assert.Len(t, comments, 3, "Should return 3 root comments")

		for _, comment := range comments {
			assert.Equal(t, 3, comment.ChildrenCount, fmt.Sprintf("Each root should have exactly 3 children: %d", comment.ChildrenCount))
		}
	})

	t.Run("GetChildComments", func(t *testing.T) {
		// 부모 댓글 생성
		postID := int64(11)
		parent := &Comment{
			PostID:  postID,
			Content: "Parent comment",
		}
		err := repo.Save(ctx, parent)
		require.NoError(t, err)

		// 자식 댓글 5개 생성
		for i := 1; i <= 5; i++ {
			child := &Comment{
				PostID:   postID,
				ParentID: &parent.ID,
				Content:  fmt.Sprintf("Child comment %d", i),
			}
			err := repo.Save(ctx, child)
			require.NoError(t, err)

			// 각 자식 댓글에 손자 댓글 2개씩 추가
			for j := 1; j <= 2; j++ {
				grandchild := &Comment{
					PostID:   postID,
					ParentID: &child.ID,
					Content:  fmt.Sprintf("Grandchild %d of child %d", j, i),
				}
				err := repo.Save(ctx, grandchild)
				require.NoError(t, err)
			}
		}

		t.Run("child comments", func(t *testing.T) {
			children, err := repo.GetChildComments(ctx, postID, parent.ID)
			require.NoError(t, err)
			assert.Len(t, children, 5, "Should return 5 child comments")
		})
	})

	t.Run("GetByID", func(t *testing.T) {
		testutils.TruncateTables(t, fixture.Conn)
		comment := &Comment{
			PostID:    1,
			Content:   "Test comment",
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}

		err := repo.Save(ctx, comment)
		require.NoError(t, err)

		retrieved, err := repo.GetByID(ctx, comment.ID)
		require.NoError(t, err)
		assert.Equal(t, comment.Content, retrieved.Content)
	})

}
