package comments

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/kimseogyu/portfolio/backend/cmd/board/internal/postings"
	"github.com/kimseogyu/portfolio/backend/internal/cstore"
	"gorm.io/gorm"
)

type Repository interface {
	Save(ctx context.Context, comment *Comment) error
	GetByID(ctx context.Context, comment_id int64) (*Comment, error)
	GetByPostID(ctx context.Context, postID int64) ([]Comment, error)
	Delete(ctx context.Context, comment_id int64) error
	GetThread(ctx context.Context, postID int64) ([]Comment, error)
	GetChildComments(ctx context.Context, postID int64, parentID int64) ([]Comment, error)
}

type repository struct {
	db         *gorm.DB
	cacheStore cstore.CacheStore
}

func NewRepository(db *gorm.DB, cacheStore cstore.CacheStore) Repository {
	return &repository{db: db, cacheStore: cacheStore}
}

func (r *repository) GetByID(ctx context.Context, comment_id int64) (*Comment, error) {
	var comment Comment
	if err := r.db.WithContext(ctx).Where("id = ?", comment_id).First(&comment).Error; err != nil {
		return nil, fmt.Errorf("db.WithContext.Where.First: %w", err)
	}
	return &comment, nil
}

func (r *repository) Save(ctx context.Context, comment *Comment) error {
	isNew := comment.ID == 0

	if err := r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Save(comment).Error; err != nil {
			return fmt.Errorf("tx.Save: %w", err)
		}

		if isNew && comment.ParentID == nil {
			if err := tx.Model(&postings.Posting{}).
				Where("id = ?", comment.PostID).
				UpdateColumn("comment_count", gorm.Expr("comment_count + ?", 1)).Error; err != nil {
				return fmt.Errorf("tx.Model.Where.UpdateColumn: %w", err)
			}
		}

		return nil
	}); err != nil {
		return fmt.Errorf("tx.Transaction: %w", err)
	}

	return nil
}

func (r *repository) GetByPostID(ctx context.Context, postID int64) ([]Comment, error) {
	cacheKey := cstore.CommentsListCacheKeyFunc(postID)
	cacheValue, err := r.cacheStore.Get(ctx, cacheKey)
	if err != nil {
		return nil, fmt.Errorf("cacheStore.Get: %w", err)
	}

	if cacheValue != "" {
		var comments []Comment
		if err := json.Unmarshal([]byte(cacheValue), &comments); err != nil {
			return nil, fmt.Errorf("json.Unmarshal: %w", err)
		}
		return comments, nil
	}

	var comments []Comment
	if err := r.db.WithContext(ctx).Where("post_id = ? AND deleted_at IS NULL", postID).Find(&comments).Error; err != nil {
		return nil, fmt.Errorf("db.WithContext.Where: %w", err)
	}

	commentsBytes, err := json.Marshal(comments)
	if err != nil {
		return nil, fmt.Errorf("json.Marshal: %w", err)
	}
	r.cacheStore.Set(ctx, cacheKey, string(commentsBytes), 24*time.Hour)

	return comments, nil
}

func (r *repository) Delete(ctx context.Context, comment_id int64) error {
	now := time.Now()
	err := r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Model(&Comment{}).Where("id = ?", comment_id).Update("deleted_at", now).Error; err != nil {
			return fmt.Errorf("tx.Model.Where.Update: %w", err)
		}
		return nil
	})

	if err != nil {
		return fmt.Errorf("tx.Transaction: %w", err)
	}

	return nil
}

// GetThread는 댓글 목록을 페이지네이션하여 가져옴
func (r *repository) GetThread(ctx context.Context, postID int64) ([]Comment, error) {
	cacheKey := cstore.CommentsListCacheKeyFunc(postID)
	cacheValue, err := r.cacheStore.Get(ctx, cacheKey)
	if err != nil {
		return nil, fmt.Errorf("cacheStore.Get: %w", err)
	}

	if cacheValue != "" {
		var rootComments []Comment
		if err := json.Unmarshal([]byte(cacheValue), &rootComments); err != nil {
			return nil, fmt.Errorf("json.Unmarshal: %w", err)
		}
		return rootComments, nil
	}

	var rootComments []Comment
	if err := r.db.WithContext(ctx).
		Where("post_id = ? and parent_id is null", postID).
		Order("created_at desc").
		Find(&rootComments).Error; err != nil {
		return nil, fmt.Errorf("db.WithContext.Select.Where.Order.Find: %w", err)
	}

	for idx, _ := range rootComments {
		children, err := r.GetChildComments(ctx, postID, rootComments[idx].ID)
		if err != nil {
			return nil, fmt.Errorf("GetChildComments: %w", err)
		}
		rootComments[idx].ChildrenCount = len(children)
	}

	rootCommentsBytes, err := json.Marshal(rootComments)
	if err != nil {
		return nil, fmt.Errorf("json.Marshal: %w", err)
	}
	r.cacheStore.Set(ctx, cacheKey, string(rootCommentsBytes), 24*time.Hour)

	return rootComments, nil
}

// GetChildComments는 특정 댓글의 자식 댓글들을 페이지네이션하여 가져옴
func (r *repository) GetChildComments(ctx context.Context, postID int64, parentID int64) ([]Comment, error) {
	var childComments []Comment
	if err := r.db.WithContext(ctx).
		Where("post_id = ? AND parent_id = ?", postID, parentID).
		Order("created_at ASC").
		Find(&childComments).Error; err != nil {
		return nil, fmt.Errorf("db.WithContext.Where.Find: %w", err)
	}

	return childComments, nil
}
