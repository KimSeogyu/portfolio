package postings

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/kimseogyu/portfolio/backend/internal/cstore"
	"github.com/kimseogyu/portfolio/backend/internal/db"
	"gorm.io/gorm"
)

type Repository interface {
	Save(ctx context.Context, data *Posting) error
	Delete(ctx context.Context, id int64) error
	FindAll(ctx context.Context, pagination *db.CursorBasedPagination) (*db.CursorBasedPaginationResponse[Posting], error)
	FindOneByID(ctx context.Context, id int64) (*Posting, error)
	Update(ctx context.Context, id int64, data Posting) error
	IncrementViewCount(ctx context.Context, id int64) error
	IncrementCommentCount(ctx context.Context, id int64) error
	DecrementCommentCount(ctx context.Context, id int64) error
}

type postingsGormRepository struct {
	db *gorm.DB

	cacheStore cstore.CacheStore
}

var _ Repository = &postingsGormRepository{}

// Save 메서드는 게시물을 저장합니다
func (p *postingsGormRepository) Save(ctx context.Context, data *Posting) error {
	tx := p.db.WithContext(ctx)
	if err := tx.Save(data).Error; err != nil {
		return fmt.Errorf("tx.Save: %w", err)
	}

	postingBytes, err := json.Marshal(data)
	if err != nil {
		return fmt.Errorf("json.Marshal: %w", err)
	}

	// 캐시 저장
	p.cacheStore.Set(ctx, cstore.PostingCacheKeyFunc(data.ID), string(postingBytes), 24*time.Hour)

	return nil
}

// Delete 메서드는 게시물을 삭제합니다
func (p *postingsGormRepository) Delete(ctx context.Context, id int64) error {
	tx := p.db.WithContext(ctx)
	posting := &Posting{}
	if err := tx.First(posting, id).Error; err != nil {
		return fmt.Errorf("tx.First: %w", err)
	}

	if err := tx.Model(posting).Update("deleted_at", gorm.Expr("NOW()")).Error; err != nil {
		return fmt.Errorf("tx.Model.Update: %w", err)
	}

	return nil
}

// FindAll 메서드는 게시물 목록을 조회합니다
func (p *postingsGormRepository) FindAll(ctx context.Context, pagination *db.CursorBasedPagination) (*db.CursorBasedPaginationResponse[Posting], error) {
	cacheKey := cstore.PostingListCacheKeyFunc(*pagination.Cursor, *pagination.Limit)
	cacheValue, err := p.cacheStore.Get(ctx, cacheKey)
	if err != nil {
		return nil, fmt.Errorf("cacheStore.Get: %w", err)
	}

	if cacheValue != "" {
		var postings *db.CursorBasedPaginationResponse[Posting]
		if err := json.Unmarshal([]byte(cacheValue), &postings); err != nil {
			return nil, fmt.Errorf("json.Unmarshal: %w", err)
		}

		return postings, nil
	}

	tx := p.db.WithContext(ctx)

	var total int64
	err = tx.Model(&Posting{}).Count(&total).Error
	if err != nil {
		return nil, fmt.Errorf("tx.Model.Count: %w", err)
	}

	tx = tx.Model(&Posting{})

	if pagination.Cursor != nil {
		tx = tx.Where("id > ?", *pagination.Cursor)
	}

	if pagination.Limit != nil {
		tx = tx.Limit(int(*pagination.Limit))
	}

	var postings []Posting
	err = tx.Find(&postings).Error
	if err != nil {
		return nil, fmt.Errorf("tx.Find: %w", err)
	}

	nextCursor := int64(0)
	if len(postings) > 0 {
		nextCursor = int64(postings[len(postings)-1].ID)
	}

	result := &db.CursorBasedPaginationResponse[Posting]{
		Data:       postings,
		NextCursor: &nextCursor,
		HasNext:    len(postings) == int(*pagination.Limit),
		Total:      int(total),
	}

	postingsBytes, err := json.Marshal(result)
	if err != nil {
		return nil, fmt.Errorf("json.Marshal: %w", err)
	}

	p.cacheStore.Set(ctx, cacheKey, string(postingsBytes), 24*time.Hour)

	return result, nil
}

// FindOneByID 메서드는 게시물을 하나 조회합니다
func (p *postingsGormRepository) FindOneByID(ctx context.Context, id int64) (*Posting, error) {
	cacheKey := cstore.PostingCacheKeyFunc(id)
	cacheValue, err := p.cacheStore.Get(ctx, cacheKey)
	if err != nil {
		return nil, fmt.Errorf("cacheStore.Get: %w", err)
	}

	if cacheValue != "" {
		var posting Posting
		if err := json.Unmarshal([]byte(cacheValue), &posting); err != nil {
			return nil, fmt.Errorf("json.Unmarshal: %w", err)
		}

		return &posting, nil
	}

	tx := p.db.WithContext(ctx)

	var posting Posting
	err = tx.Where("id = ?", id).First(&posting).Error
	if err != nil {
		return nil, fmt.Errorf("tx.Where.First: %w", err)
	}

	postingBytes, err := json.Marshal(&posting)
	if err != nil {
		return nil, fmt.Errorf("json.Marshal: %w", err)
	}

	p.cacheStore.Set(ctx, cstore.PostingCacheKeyFunc(id), string(postingBytes), 24*time.Hour)

	return &posting, nil
}

// Update 메서드는 게시물을 업데이트합니다
func (p *postingsGormRepository) Update(ctx context.Context, id int64, data Posting) error {
	tx := p.db.WithContext(ctx)
	var posting Posting
	if err := tx.First(&posting, id).Error; err != nil {
		return fmt.Errorf("tx.First: %w", err)
	}

	if err := tx.Model(&posting).Updates(data).Error; err != nil {
		return fmt.Errorf("tx.Model.Updates: %w", err)
	}

	postingBytes, err := json.Marshal(posting)
	if err != nil {
		return fmt.Errorf("json.Marshal: %w", err)
	}

	p.cacheStore.Set(ctx, cstore.PostingCacheKeyFunc(id), string(postingBytes), 24*time.Hour)

	return nil
}

// IncrementViewCount 메서드는 게시물의 조회수를 증가시킵니다
func (p *postingsGormRepository) IncrementViewCount(ctx context.Context, id int64) error {
	tx := p.db.WithContext(ctx)

	var posting Posting
	err := tx.Transaction(func(tx *gorm.DB) error {
		if err := tx.First(&posting, id).Error; err != nil {
			return fmt.Errorf("tx.First: %w", err)
		}

		posting.ViewCount++
		if err := tx.Model(&Posting{}).Where("id = ?", id).
			Update("view_count", gorm.Expr("view_count + ?", 1)).Error; err != nil {
			return fmt.Errorf("tx.Model.Where.Update: %w", err)
		}
		return nil
	})
	if err != nil {
		return fmt.Errorf("tx.Transaction: %w", err)
	}

	return nil
}

// IncrementCommentCount 메서드는 게시물의 댓글 수를 증가시킵니다
func (p *postingsGormRepository) IncrementCommentCount(ctx context.Context, id int64) error {
	tx := p.db.WithContext(ctx)
	err := tx.Transaction(func(tx *gorm.DB) error {
		var posting Posting
		if err := tx.First(&posting, id).Error; err != nil {
			return fmt.Errorf("tx.First: %w", err)
		}

		posting.CommentCount++
		if err := tx.Model(&Posting{}).Where("id = ?", id).
			Update("comment_count", gorm.Expr("comment_count + ?", 1)).Error; err != nil {
			return fmt.Errorf("tx.Model.Where.Update: %w", err)
		}
		return nil
	})
	if err != nil {
		return fmt.Errorf("tx.Transaction: %w", err)
	}

	return nil
}

// DecrementCommentCount 메서드는 게시물의 댓글 수를 감소시킵니다
func (p *postingsGormRepository) DecrementCommentCount(ctx context.Context, id int64) error {
	tx := p.db.WithContext(ctx)
	err := tx.Transaction(func(tx *gorm.DB) error {
		var posting Posting
		if err := tx.First(&posting, id).Error; err != nil {
			return fmt.Errorf("tx.First: %w", err)
		}

		posting.CommentCount--
		if err := tx.Model(&Posting{}).Where("id = ?", id).
			Update("comment_count", gorm.Expr("comment_count - ?", 1)).Error; err != nil {
			return fmt.Errorf("tx.Model.Where.Update: %w", err)
		}
		return nil
	})
	if err != nil {
		return fmt.Errorf("tx.Transaction: %w", err)
	}

	return nil
}

func NewRepository(db *gorm.DB, cacheStore cstore.CacheStore) *postingsGormRepository {
	return &postingsGormRepository{
		db:         db,
		cacheStore: cacheStore,
	}
}
