package postings

import (
	"time"

	boardServer "github.com/kimseogyu/portfolio/backend/internal/proto/board/v1"
)

type Posting struct {
	ID           int                       `json:"id" gorm:"primaryKey"`
	Title        string                    `json:"title" gorm:"not null;size:200;unique"`
	Content      string                    `json:"content" gorm:"not null;size:50000"`
	AuthorID     string                    `json:"author_id" gorm:"not null"`
	AuthorName   string                    `json:"author_name" gorm:"not null"`
	CommentCount int                       `json:"comment_count" gorm:"default:0"`
	ViewCount    int                       `json:"view_count" gorm:"default:0"`
	Tags         []string                  `json:"tags" gorm:"-"`
	Status       boardServer.PostingStatus `json:"status" gorm:"not null;size:20"`
	CreatedAt    time.Time                 `json:"created_at"`
	UpdatedAt    time.Time                 `json:"updated_at"`
	DeletedAt    *time.Time                `json:"deleted_at"`
}
