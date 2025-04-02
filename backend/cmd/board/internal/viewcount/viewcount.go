package viewcount

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/kimseogyu/portfolio/backend/cmd/board/internal/postings"
	"github.com/kimseogyu/portfolio/backend/internal/cstore"
	"go.uber.org/zap"
)

type ViewCountManager interface {
	CheckAndIncrement(ctx context.Context, postingID int64, authorID string, authUserID string) error
}

type viewCountManagerImpl struct {
	cacheStore        cstore.CacheStore
	postingRepository postings.Repository
}

type options struct {
	cacheStore        cstore.CacheStore
	postingRepository postings.Repository
}

type OptionFunc func(*options)

func WithCacheStore(cacheStore cstore.CacheStore) OptionFunc {
	return func(o *options) {
		o.cacheStore = cacheStore
	}
}

func WithPostingRepository(postingRepository postings.Repository) OptionFunc {
	return func(o *options) {
		o.postingRepository = postingRepository
	}
}

func NewViewCountManager(opts ...OptionFunc) (ViewCountManager, error) {
	options := &options{}
	for _, opt := range opts {
		opt(options)
	}

	if options.cacheStore == nil {
		return nil, errors.New("cacheStore is required")
	}
	if options.postingRepository == nil {
		return nil, errors.New("postingRepository is required")
	}

	return &viewCountManagerImpl{
		cacheStore:        options.cacheStore,
		postingRepository: options.postingRepository,
	}, nil
}

func (v *viewCountManagerImpl) CheckAndIncrement(ctx context.Context, postingID int64, authorID string, authUserID string) error {
	// 캐시 키 생성 (Redis 또는 메모리 캐시 사용)
	viewCountKey := fmt.Sprintf("viewcount:%d:%s", postingID, authUserID)

	// 캐시 확인
	exists, err := v.cacheStore.Exists(ctx, viewCountKey)
	if err != nil {
		return fmt.Errorf("failed to check cache: %w", err)
	}

	if !exists {
		// 캐시에 기록 (24시간 유효)
		v.cacheStore.Set(ctx, viewCountKey, "1", 24*time.Hour)
		go func() {
			// 조회수 증가
			if err := v.postingRepository.IncrementViewCount(ctx, postingID); err != nil {
				zap.S().Errorf("failed to increment view count: %v", err)
			}
		}()
	}

	return nil
}
