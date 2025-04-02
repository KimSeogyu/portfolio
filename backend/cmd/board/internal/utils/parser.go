package utils

import (
	"github.com/kimseogyu/portfolio/backend/cmd/board/internal/comments"
	"github.com/kimseogyu/portfolio/backend/cmd/board/internal/postings"
	boardServer "github.com/kimseogyu/portfolio/backend/internal/proto/board/v1"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func CommentToProto(comment *comments.Comment) *boardServer.Comment {
	var deletedAt *timestamppb.Timestamp
	if comment.DeletedAt != nil {
		deletedAt = timestamppb.New(*comment.DeletedAt)
	}
	var parentId int64
	if comment.ParentID != nil {
		parentId = *comment.ParentID
	}

	return &boardServer.Comment{
		CommentId:  comment.ID,
		PostingId:  comment.PostID,
		Content:    comment.Content,
		AuthorId:   comment.AuthorID,
		AuthorName: comment.AuthorName,
		CreatedAt:  timestamppb.New(comment.CreatedAt),
		UpdatedAt:  timestamppb.New(comment.UpdatedAt),
		DeletedAt:  deletedAt,
		ParentId:   parentId,
		Status:     comment.Status,
		Children:   nil,
	}
}

func PostingToProto(posting *postings.Posting) *boardServer.Posting {
	var deletedAt *timestamppb.Timestamp
	if posting.DeletedAt != nil {
		deletedAt = timestamppb.New(*posting.DeletedAt)
	}

	return &boardServer.Posting{
		PostingId:    posting.ID,
		Title:        posting.Title,
		Content:      posting.Content,
		AuthorId:     posting.AuthorID,
		AuthorName:   posting.AuthorName,
		CreatedAt:    timestamppb.New(posting.CreatedAt),
		UpdatedAt:    timestamppb.New(posting.UpdatedAt),
		DeletedAt:    deletedAt,
		ViewCount:    int32(posting.ViewCount),
		CommentCount: int32(posting.CommentCount),
		Tags:         posting.Tags,
		Status:       posting.Status,
		Comments:     nil,
	}
}

func ProtoToComment(proto *boardServer.Comment) *comments.Comment {
	deletedAt := proto.DeletedAt.AsTime()
	return &comments.Comment{
		ID:            proto.CommentId,
		PostID:        proto.PostingId,
		ParentID:      &proto.ParentId,
		AuthorID:      proto.AuthorId,
		AuthorName:    proto.AuthorName,
		ChildrenCount: 0,
		Content:       proto.Content,
		CreatedAt:     proto.CreatedAt.AsTime(),
		UpdatedAt:     proto.UpdatedAt.AsTime(),
		DeletedAt:     &deletedAt,
		Status:        proto.Status,
	}
}

func ProtoToPosting(proto *boardServer.Posting) *postings.Posting {
	deletedAt := proto.DeletedAt.AsTime()
	return &postings.Posting{
		ID:           proto.PostingId,
		Title:        proto.Title,
		Content:      proto.Content,
		AuthorID:     proto.AuthorId,
		AuthorName:   proto.AuthorName,
		CommentCount: int(proto.CommentCount),
		ViewCount:    int(proto.ViewCount),
		Tags:         proto.Tags,
		Status:       proto.Status,
		CreatedAt:    proto.CreatedAt.AsTime(),
		UpdatedAt:    proto.UpdatedAt.AsTime(),
		DeletedAt:    &deletedAt,
	}
}
