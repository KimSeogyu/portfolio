package comments

import (
	"context"
	"fmt"
	"time"

	"github.com/kimseogyu/portfolio/backend/cmd/board/internal/postings"
	"gorm.io/gorm"
)

type Repository interface {
	Save(ctx context.Context, comment *Comment) error
	GetByID(ctx context.Context, id int64) (*Comment, error)
	GetByPostID(ctx context.Context, postID int64) ([]Comment, error)
	Delete(ctx context.Context, id int64) error
	GetThread(ctx context.Context, postID int64, page, pageSize int) ([]Comment, error)
	GetChildComments(ctx context.Context, postID int64, parentID int64, page, pageSize int) ([]Comment, error)
	UpdatePostingCommentCount(ctx context.Context, postID int64) error
}

type repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) Repository {
	return &repository{db: db}
}

func (r *repository) Save(ctx context.Context, comment *Comment) error {
	isNew := comment.ID == 0

	if err := r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// 댓글 저장
		if err := tx.Save(comment).Error; err != nil {
			return fmt.Errorf("failed to save comment: %w", err)
		}

		// 새 댓글이고 최상위 댓글(parent_id가 nil)인 경우에만 게시물의 댓글 수 증가
		if isNew && comment.ParentID == nil {
			// GORM 메서드 체이닝으로 업데이트
			if err := tx.Model(&postings.Posting{}).
				Where("id = ?", comment.PostID).
				UpdateColumn("comment_count", gorm.Expr("comment_count + ?", 1)).Error; err != nil {
				return fmt.Errorf("failed to update posting comment count: %w", err)
			}
		}

		return nil
	}); err != nil {
		return err
	}

	return nil
}

func (r *repository) GetByID(ctx context.Context, id int64) (*Comment, error) {
	var comment Comment
	if err := r.db.WithContext(ctx).First(&comment, id).Error; err != nil {
		return nil, fmt.Errorf("failed to get comment by ID: %w", err)
	}

	if comment.DeletedAt != nil {
		return nil, fmt.Errorf("comment not found")
	}

	return &comment, nil
}

func (r *repository) GetByPostID(ctx context.Context, postID int64) ([]Comment, error) {
	var comments []Comment
	if err := r.db.WithContext(ctx).Where("post_id = ?", postID).Find(&comments).Error; err != nil {
		return nil, fmt.Errorf("failed to get comments by post ID: %w", err)
	}

	for i := range comments {
		if comments[i].DeletedAt != nil {
			comments = append(comments[:i], comments[i+1:]...)
		}
	}
	return comments, nil
}

func (r *repository) Delete(ctx context.Context, id int64) error {
	// 댓글과 그 하위 댓글들을 모두 소프트 삭제
	now := time.Now()
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// 먼저 삭제할 댓글 정보 가져오기
		comment := Comment{}
		if err := tx.First(&comment, id).Error; err != nil {
			return fmt.Errorf("failed to find comment: %w", err)
		}

		// 최상위 댓글인지 확인
		isRootComment := comment.ParentID == nil
		postID := comment.PostID

		// 하위 댓글들을 재귀적으로 소프트 삭제
		var childComments []Comment
		if err := tx.Where("parent_id = ? AND deleted_at IS NULL", id).Find(&childComments).Error; err != nil {
			return fmt.Errorf("failed to get child comments: %w", err)
		}

		for _, child := range childComments {
			if err := r.Delete(ctx, child.ID); err != nil {
				return fmt.Errorf("failed to delete child comment: %w", err)
			}
		}

		// 현재 댓글 소프트 삭제
		comment.DeletedAt = &now
		if err := tx.Save(&comment).Error; err != nil {
			return fmt.Errorf("failed to soft delete comment: %w", err)
		}

		// 최상위 댓글인 경우에만 게시물의 댓글 수 감소
		if isRootComment {
			// GORM 메서드 체이닝으로 업데이트
			if err := tx.Model(&postings.Posting{}).
				Where("id = ?", postID).
				UpdateColumn("comment_count", gorm.Expr("GREATEST(0, comment_count - ?)", 1)).Error; err != nil {
				return fmt.Errorf("failed to update posting comment count: %w", err)
			}
		}

		return nil
	})
}

// UpdatePostingCommentCount는 게시물의 댓글 수를 올바르게 업데이트합니다
func (r *repository) UpdatePostingCommentCount(ctx context.Context, postID int64) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// 최상위 댓글(parent_id가 null이고 삭제되지 않은)의 수를 계산
		var count int64
		if err := tx.Model(&Comment{}).
			Where("post_id = ? AND parent_id IS NULL AND deleted_at IS NULL", postID).
			Count(&count).Error; err != nil {
			return fmt.Errorf("failed to count root comments: %w", err)
		}

		// 게시물의 comment_count 업데이트
		if err := tx.Model(&postings.Posting{}).
			Where("id = ?", postID).
			Update("comment_count", count).Error; err != nil {
			return fmt.Errorf("failed to update posting comment count: %w", err)
		}

		return nil
	})
}

// GetThread는 댓글 목록을 페이지네이션하여 가져옴
func (r *repository) GetThread(ctx context.Context, postID int64, page, pageSize int) ([]Comment, error) {
	offset := (page - 1) * pageSize

	var rootComments []Comment
	// 최상위 댓글만 페이지네이션하여 가져옴
	if err := r.db.WithContext(ctx).
		Where("post_id = ? AND parent_id IS NULL", postID).
		Order("created_at DESC").
		Offset(offset).
		Limit(pageSize).
		Find(&rootComments).Error; err != nil {
		return nil, fmt.Errorf("failed to get thread: %w", err)
	}

	// 각 최상위 댓글에 대해 하위 댓글 갯수를 가져옴
	for i := range rootComments {
		var count int64
		if err := r.db.WithContext(ctx).
			Model(&Comment{}).
			Where("post_id = ? AND parent_id = ?", rootComments[i].PostID, rootComments[i].ID).
			Count(&count).Error; err != nil {
			return nil, fmt.Errorf("failed to get children count: %w", err)
		}
		rootComments[i].ChildrenCount = int(count)
	}

	return rootComments, nil
}

// GetChildComments는 특정 댓글의 자식 댓글들을 페이지네이션하여 가져옴
func (r *repository) GetChildComments(ctx context.Context, postID int64, parentID int64, page, pageSize int) ([]Comment, error) {
	offset := (page - 1) * pageSize

	var childComments []Comment
	if err := r.db.WithContext(ctx).
		Where("post_id = ? AND parent_id = ?", postID, parentID).
		Order("created_at ASC").
		Offset(offset).
		Limit(pageSize).
		Find(&childComments).Error; err != nil {
		return nil, fmt.Errorf("failed to get child comments: %w", err)
	}

	return childComments, nil
}
