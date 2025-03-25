package comments

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/kimseogyu/portfolio/backend/cmd/board/internal/postings"
	"github.com/kimseogyu/portfolio/backend/internal/db"
	boardServer "github.com/kimseogyu/portfolio/backend/internal/proto/board/v1"
	"github.com/kimseogyu/portfolio/backend/internal/testutils"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/gorm"
)

func setupTestDB(t *testing.T) *gorm.DB {
	pg, err := testutils.NewPostgresTestContainer(context.Background())
	require.NoError(t, err)

	t.Cleanup(func() {
		err = pg.Close()
		require.NoError(t, err)
	})

	endpoint, err := pg.Endpoint()
	require.NoError(t, err)

	db, err := db.NewPostgresDB(db.WithDSN(endpoint))
	require.NoError(t, err)

	// Auto-migrate the Comment 및 Posting 모델
	err = db.AutoMigrate(&Comment{}, &postings.Posting{})
	require.NoError(t, err)

	// 테이블 truncate로 깨끗한 상태 유지
	truncateTables(t, db)

	return db
}

// 테이블 데이터 정리 함수
func truncateTables(t *testing.T, db *gorm.DB) {
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
	db := setupTestDB(t)
	repo := NewRepository(db)
	ctx := context.Background()

	t.Run("Save and GetByID", func(t *testing.T) {
		truncateTables(t, db)
		comment := &Comment{
			PostID:    1,
			Content:   "Test comment",
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}

		err := repo.Save(ctx, comment)
		assert.NoError(t, err)
		assert.NotZero(t, comment.ID)

		retrieved, err := repo.GetByID(ctx, comment.ID)
		assert.NoError(t, err)
		assert.Equal(t, comment.Content, retrieved.Content)
	})

	t.Run("Save with ParentID", func(t *testing.T) {
		truncateTables(t, db)
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
		truncateTables(t, db)
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

		_, err = repo.GetByID(ctx, comment.ID)
		assert.Error(t, err) // Should return error as comment is deleted
	})

	t.Run("GetThread", func(t *testing.T) {

	})

	t.Run("Delete with Children", func(t *testing.T) {
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

		// Verify both parent and child are deleted
		_, err = repo.GetByID(ctx, parent.ID)
		assert.Error(t, err)
		_, err = repo.GetByID(ctx, child.ID)
		assert.Error(t, err)
	})

	t.Run("GetThread with Pagination", func(t *testing.T) {
		// 테스트 데이터 생성: 최상위 댓글 3개
		postID := int64(10)
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

		// 페이지네이션 테스트
		t.Run("First page", func(t *testing.T) {
			comments, err := repo.GetThread(ctx, postID, 1, 2)
			require.NoError(t, err)
			assert.Len(t, comments, 2, "Should return 2 root comments")

			// 각 루트 댓글은 첫 번째 자식 댓글만 가져와야 함
			for _, comment := range comments {
				assert.Equal(t, 3, comment.ChildrenCount, "Each root should have exactly 3 children")
			}
		})

		t.Run("Second page", func(t *testing.T) {
			comments, err := repo.GetThread(ctx, postID, 2, 2)
			require.NoError(t, err)
			assert.Len(t, comments, 1, "Should return 1 root comment")
		})
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

		t.Run("Paginated child comments", func(t *testing.T) {
			// 첫 페이지 테스트 (2개)
			children, err := repo.GetChildComments(ctx, postID, parent.ID, 1, 2)
			require.NoError(t, err)
			assert.Len(t, children, 2, "Should return 2 child comments")

			// 두 번째 페이지 테스트 (2개)
			children, err = repo.GetChildComments(ctx, postID, parent.ID, 2, 2)
			require.NoError(t, err)
			assert.Len(t, children, 2, "Should return 2 child comments")
		})
	})

	// 새로운 테스트 케이스 추가
	t.Run("CommentCount Updates", func(t *testing.T) {
		truncateTables(t, db)
		// 새 게시물 생성
		postID := createTestPosting(t, db)

		// 초기 댓글 수 확인 (DB 직접 조회)
		count := getCommentCountFromDB(t, db, postID)
		assert.Equal(t, 0, count, "Initial comment count should be 0")

		// 1. 최상위 댓글 생성 시 댓글 수 증가 확인
		rootComment := &Comment{
			PostID:     postID,
			Content:    "Root comment for counting test",
			AuthorID:   "test-user",
			AuthorName: "Test User",
			CreatedAt:  time.Now(),
			UpdatedAt:  time.Now(),
			Status:     boardServer.CommentStatus_COMMENT_STATUS_PUBLISHED,
		}

		err := repo.Save(ctx, rootComment)
		require.NoError(t, err)

		// 댓글 수가 1로 증가했는지 확인
		count = getCommentCountFromDB(t, db, postID)
		assert.Equal(t, 1, count, "Comment count should increase to 1 after adding root comment")

		// 2. 대댓글 생성 시 댓글 수 변화 없는지 확인
		childComment := &Comment{
			PostID:     postID,
			ParentID:   &rootComment.ID,
			Content:    "Child comment for counting test",
			AuthorID:   "test-user",
			AuthorName: "Test User",
			CreatedAt:  time.Now(),
			UpdatedAt:  time.Now(),
			Status:     boardServer.CommentStatus_COMMENT_STATUS_PUBLISHED,
		}

		err = repo.Save(ctx, childComment)
		require.NoError(t, err)

		// 댓글 수가 여전히 1인지 확인 (대댓글은 카운트에 영향 없음)
		count = getCommentCountFromDB(t, db, postID)
		assert.Equal(t, 1, count, "Comment count should remain 1 after adding child comment")

		// 3. 두 번째 최상위 댓글 생성
		secondRootComment := &Comment{
			PostID:     postID,
			Content:    "Second root comment",
			AuthorID:   "test-user",
			AuthorName: "Test User",
			CreatedAt:  time.Now(),
			UpdatedAt:  time.Now(),
			Status:     boardServer.CommentStatus_COMMENT_STATUS_PUBLISHED,
		}

		err = repo.Save(ctx, secondRootComment)
		require.NoError(t, err)

		// 댓글 수가 2로 증가했는지 확인
		count = getCommentCountFromDB(t, db, postID)
		assert.Equal(t, 2, count, "Comment count should increase to 2 after adding second root comment")

		// 4. 최상위 댓글 삭제 시 댓글 수 감소 확인
		err = repo.Delete(ctx, rootComment.ID)
		require.NoError(t, err)

		// 댓글 수가 1로 감소했는지 확인
		count = getCommentCountFromDB(t, db, postID)
		assert.Equal(t, 1, count, "Comment count should decrease to 1 after deleting root comment")

		// 5. 강제 업데이트 함수 테스트
		// 먼저 DB에 직접 잘못된 값 설정
		err = db.Exec("UPDATE postings SET comment_count = 10 WHERE id = ?", postID).Error
		require.NoError(t, err)

		// 값이 변경되었는지 확인
		count = getCommentCountFromDB(t, db, postID)
		assert.Equal(t, 10, count, "Comment count should be manually set to 10")

		// UpdatePostingCommentCount 함수 호출
		err = repo.UpdatePostingCommentCount(ctx, postID)
		require.NoError(t, err)

		// 값이 올바르게 수정되었는지 확인 (실제로 남아있는 최상위 댓글은 1개)
		count = getCommentCountFromDB(t, db, postID)
		assert.Equal(t, 1, count, "Comment count should be corrected to 1 after forced update")
	})
}
