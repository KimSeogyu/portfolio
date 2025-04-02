package comments

import (
	"time"

	boardServer "github.com/kimseogyu/portfolio/backend/internal/proto/board/v1"
)

type Comment struct {
	ID            int64                     `json:"id" gorm:"primaryKey"`
	PostID        int64                     `json:"post_id"`
	ParentID      *int64                    `json:"parent_id"`
	AuthorID      string                    `json:"author_id"`
	AuthorName    string                    `json:"author_name"`
	ChildrenCount int                       `json:"children_count"`
	Content       string                    `json:"content"`
	CreatedAt     time.Time                 `json:"created_at"`
	UpdatedAt     time.Time                 `json:"updated_at"`
	DeletedAt     *time.Time                `json:"deleted_at"`
	Status        boardServer.CommentStatus `json:"status"`
}
