package server

import (
	"context"
	"errors"
	"time"

	"github.com/kimseogyu/portfolio/backend/cmd/board/internal/comments"
	"github.com/kimseogyu/portfolio/backend/cmd/board/internal/postings"
	"github.com/kimseogyu/portfolio/backend/cmd/board/internal/utils"
	"github.com/kimseogyu/portfolio/backend/cmd/board/internal/viewcount"
	"github.com/kimseogyu/portfolio/backend/internal/cstore"
	"github.com/kimseogyu/portfolio/backend/internal/db"
	"github.com/kimseogyu/portfolio/backend/internal/dlock"
	boardServer "github.com/kimseogyu/portfolio/backend/internal/proto/board/v1"
	"github.com/kimseogyu/portfolio/backend/pkg/authenticator"
	"github.com/kimseogyu/portfolio/backend/pkg/pagination"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type Service struct {
	boardServer.UnimplementedBoardServiceServer

	commentRepository comments.Repository
	postingRepository postings.Repository
	authenticator     authenticator.UserAuthenticator
	cacheStore        cstore.CacheStore
	dlockerFactory    dlock.DLockerFactory
	viewCountManager  viewcount.ViewCountManager
}

// CreateComment implements boardv1.BoardServiceServer.
func (s *Service) CreateComment(ctx context.Context, req *boardServer.CreateCommentRequest) (*boardServer.Comment, error) {
	authUser, err := s.authenticator.FromGrpcContext(ctx)
	if err != nil {
		return nil, err
	}

	now := time.Now()
	comment := utils.ProtoToComment(&boardServer.Comment{
		PostingId:  req.PostingId,
		Content:    req.Content,
		AuthorId:   authUser.ID,
		AuthorName: authUser.Name,
		ParentId:   req.ParentId,
		CreatedAt:  timestamppb.New(now),
		UpdatedAt:  timestamppb.New(now),
		Status:     boardServer.CommentStatus_COMMENT_STATUS_PUBLISHED,
	})

	if err := s.commentRepository.Save(ctx, comment); err != nil {
		return nil, status.Errorf(codes.Internal, "commentRepository.Save: %v", err)
	}

	return utils.CommentToProto(comment), nil
}

// CreatePosting implements boardv1.BoardServiceServer.
func (s *Service) CreatePosting(ctx context.Context, req *boardServer.CreatePostingRequest) (*boardServer.Posting, error) {
	authUser, err := s.authenticator.FromGrpcContext(ctx)
	if err != nil {
		return nil, err
	}

	now := time.Now()
	posting := utils.ProtoToPosting(&boardServer.Posting{
		Title:      req.Title,
		Content:    req.Content,
		AuthorId:   authUser.ID,
		AuthorName: authUser.Name,
		CreatedAt:  timestamppb.New(now),
		UpdatedAt:  timestamppb.New(now),
		DeletedAt:  nil,
		Tags:       req.Tags,
		Status:     req.Status,
		Comments:   nil,
	})

	if err := s.postingRepository.Save(ctx, posting); err != nil {
		return nil, status.Errorf(codes.Internal, "postingRepository.Save: %v", err)
	}

	return utils.PostingToProto(posting), nil
}

// DeleteComment implements boardv1.BoardServiceServer.
func (s *Service) DeleteComment(ctx context.Context, req *boardServer.DeleteCommentRequest) (*emptypb.Empty, error) {
	authUser, err := s.authenticator.FromGrpcContext(ctx)
	if err != nil {
		return nil, err
	}

	comment, err := s.commentRepository.GetByID(ctx, req.CommentId)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "commentRepository.GetByID: %v", err)
	}

	if comment.AuthorID != authUser.ID {
		return nil, status.Errorf(codes.PermissionDenied, "you are not allowed to delete this comment")
	}

	if err := s.commentRepository.Delete(ctx, req.CommentId); err != nil {
		return nil, status.Errorf(codes.Internal, "commentRepository.Delete: %v", err)
	}

	return &emptypb.Empty{}, nil
}

// DeletePosting implements boardv1.BoardServiceServer.
func (s *Service) DeletePosting(ctx context.Context, req *boardServer.DeletePostingRequest) (*emptypb.Empty, error) {
	authUser, err := s.authenticator.FromGrpcContext(ctx)
	if err != nil {
		return nil, err
	}

	posting, err := s.postingRepository.FindOneByID(ctx, req.PostingId)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "postingRepository.FindOneByID: %v", err)
	}

	if posting.AuthorID != authUser.ID {
		return nil, status.Errorf(codes.PermissionDenied, "you are not allowed to delete this posting")
	}

	if err := s.postingRepository.Delete(ctx, req.PostingId); err != nil {
		return nil, status.Errorf(codes.Internal, "postingRepository.Delete: %v", err)
	}

	return &emptypb.Empty{}, nil
}

// GetPosting implements boardv1.BoardServiceServer.
func (s *Service) GetPosting(ctx context.Context, req *boardServer.GetPostingRequest) (*boardServer.Posting, error) {
	// 인증된 사용자 정보 가져오기 (익명 사용자도 허용)
	authUser, err := s.authenticator.FromGrpcContext(ctx)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "authenticator.FromGrpcContext: %v", err)
	}

	posting, err := s.postingRepository.FindOneByID(ctx, req.PostingId)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "postingRepository.FindOneByID: %v", err)
	}

	if err := s.viewCountManager.CheckAndIncrement(ctx, posting.ID, posting.AuthorID, authUser.ID); err != nil {
		return nil, status.Errorf(codes.Internal, "viewCountManager.CheckAndIncrement: %v", err)
	}

	comments, err := s.commentRepository.GetByPostID(ctx, req.PostingId)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "commentRepository.GetByPostID: %v", err)
	}

	boardComments := make([]*boardServer.Comment, len(comments))
	for i, comment := range comments {
		boardComments[i] = utils.CommentToProto(&comment)
	}

	postingForResponse := utils.PostingToProto(posting)
	postingForResponse.Comments = boardComments

	return postingForResponse, nil
}

// ListCommentsByPosting implements boardv1.BoardServiceServer.
func (s *Service) ListCommentsByPosting(ctx context.Context, req *boardServer.ListCommentsByPostingRequest) (*boardServer.ListCommentsResponse, error) {
	comments, err := s.commentRepository.GetByPostID(ctx, req.PostingId)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "commentRepository.GetByPostID: %v", err)
	}

	boardComments := make([]*boardServer.Comment, len(comments))
	for i, comment := range comments {
		boardComments[i] = utils.CommentToProto(&comment)
	}

	return &boardServer.ListCommentsResponse{Comments: boardComments}, nil
}

// ListPostings implements boardv1.BoardServiceServer.
func (s *Service) ListPostings(ctx context.Context, req *boardServer.ListPostingsRequest) (*boardServer.ListPostingsResponse, error) {
	token, err := pagination.FromEncodedString(req.PageToken)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "pagination.FromEncodedString: %v", err)
	}

	postings, err := s.postingRepository.FindAll(ctx, &db.CursorBasedPagination{
		Cursor: token.Cursor,
		Limit:  token.Limit,
	})
	if err != nil {
		return nil, status.Errorf(codes.Internal, "postingRepository.FindAll: %v", err)
	}

	boardPostings := make([]*boardServer.Posting, len(postings.Data))
	for i, posting := range postings.Data {
		boardPostings[i] = utils.PostingToProto(&posting)
	}

	nextCursor := int64(0)
	if len(postings.Data) > 0 {
		nextCursor = int64(postings.Data[len(postings.Data)-1].ID)
	}

	nextPageToken := pagination.NewToken(nextCursor, *token.Limit).Encode()

	return &boardServer.ListPostingsResponse{
		Postings:      boardPostings,
		NextPageToken: nextPageToken,
		TotalCount:    int32(postings.Total),
	}, nil
}

// SearchPostings implements boardv1.BoardServiceServer.
func (s *Service) SearchPostings(ctx context.Context, req *boardServer.SearchPostingsRequest) (*boardServer.ListPostingsResponse, error) {
	panic("unimplemented")
}

// UpdateComment implements boardv1.BoardServiceServer.
func (s *Service) UpdateComment(ctx context.Context, req *boardServer.UpdateCommentRequest) (*boardServer.Comment, error) {
	authUser, err := s.authenticator.FromGrpcContext(ctx)
	if err != nil {
		return nil, err
	}

	comment, err := s.commentRepository.GetByID(ctx, req.CommentId)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "commentRepository.GetByID: %v", err)
	}

	if comment.AuthorID != authUser.ID {
		return nil, status.Errorf(codes.PermissionDenied, "you are not allowed to update this comment")
	}

	comment.Content = req.Content
	comment.UpdatedAt = time.Now()

	if err := s.commentRepository.Save(ctx, comment); err != nil {
		return nil, status.Errorf(codes.Internal, "commentRepository.Save: %v", err)
	}

	return utils.CommentToProto(comment), nil
}

// UpdatePosting implements boardv1.BoardServiceServer.
func (s *Service) UpdatePosting(ctx context.Context, req *boardServer.UpdatePostingRequest) (*boardServer.Posting, error) {
	authUser, err := s.authenticator.FromGrpcContext(ctx)
	if err != nil {
		return nil, err
	}

	posting, err := s.postingRepository.FindOneByID(ctx, req.PostingId)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "postingRepository.FindOneByID: %v", err)
	}

	if posting.AuthorID != authUser.ID {
		return nil, status.Errorf(codes.PermissionDenied, "you are not allowed to update this posting")
	}

	posting.Title = req.Title
	posting.Content = req.Content
	posting.Tags = req.Tags
	posting.Status = req.Status
	posting.UpdatedAt = time.Now()

	if err := s.postingRepository.Save(ctx, posting); err != nil {
		return nil, status.Errorf(codes.Internal, "postingRepository.Save: %v", err)
	}

	return utils.PostingToProto(posting), nil
}

var _ boardServer.BoardServiceServer = &Service{}

type ServiceOption func(*Service)

func WithCommentRepository(commentRepository comments.Repository) ServiceOption {
	return func(s *Service) {
		s.commentRepository = commentRepository
	}
}

func WithPostingRepository(postingRepository postings.Repository) ServiceOption {
	return func(s *Service) {
		s.postingRepository = postingRepository
	}
}

func WithAuthenticator(authenticator authenticator.UserAuthenticator) ServiceOption {
	return func(s *Service) {
		s.authenticator = authenticator
	}
}

func WithCacheStore(cache cstore.CacheStore) ServiceOption {
	return func(s *Service) {
		s.cacheStore = cache
	}
}

func WithDLockerFactory(dlockerFactory dlock.DLockerFactory) ServiceOption {
	return func(s *Service) {
		s.dlockerFactory = dlockerFactory
	}
}

func WithViewCountManager(viewCountManager viewcount.ViewCountManager) ServiceOption {
	return func(s *Service) {
		s.viewCountManager = viewCountManager
	}
}

func NewService(opts ...ServiceOption) (*Service, error) {
	s := &Service{}
	for _, opt := range opts {
		opt(s)
	}

	if s.commentRepository == nil {
		return nil, errors.New("commentRepository is required")
	}

	if s.postingRepository == nil {
		return nil, errors.New("postingRepository is required")
	}

	if s.authenticator == nil {
		return nil, errors.New("authenticator is required")
	}

	if s.cacheStore == nil {
		return nil, errors.New("cache is required")
	}

	if s.dlockerFactory == nil {
		return nil, errors.New("dlockerFactory is required")
	}

	return s, nil
}
