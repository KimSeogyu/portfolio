package server

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/kimseogyu/portfolio/backend/cmd/board/internal/comments"
	"github.com/kimseogyu/portfolio/backend/cmd/board/internal/postings"
	"github.com/kimseogyu/portfolio/backend/internal/cache"
	"github.com/kimseogyu/portfolio/backend/internal/db"
	boardServer "github.com/kimseogyu/portfolio/backend/internal/proto/board/v1"
	"github.com/kimseogyu/portfolio/backend/pkg/authenticator"
	"github.com/kimseogyu/portfolio/backend/pkg/pagination"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type Service struct {
	boardServer.UnimplementedBoardServiceServer

	commentRepository comments.Repository
	postingRepository postings.Repository
	authenticator     authenticator.UserAuthenticator
	cache             cache.Cache
}

// CreateComment implements boardv1.BoardServiceServer.
func (s *Service) CreateComment(ctx context.Context, req *boardServer.CreateCommentRequest) (*boardServer.Comment, error) {
	authUser, err := s.authenticator.FromGrpcContext(ctx)
	if err != nil {
		return nil, err
	}

	var parentID *int64
	if req.ParentId != 0 {
		parentID = &req.ParentId
	}

	comment := &comments.Comment{
		PostID:     req.PostingId,
		Content:    req.Content,
		AuthorID:   authUser.ID,
		AuthorName: authUser.Name,
		ParentID:   parentID,
		CreatedAt:  time.Now(),
		UpdatedAt:  time.Now(),
	}

	if err := s.commentRepository.Save(ctx, comment); err != nil {
		return nil, status.Errorf(codes.Internal, "failed to create comment: %v", err)
	}

	depth := 0
	if req.ParentId != 0 {
		depth = 1
	}

	return &boardServer.Comment{
		CommentId:  comment.ID,
		PostingId:  comment.PostID,
		Content:    comment.Content,
		AuthorId:   authUser.ID,
		AuthorName: authUser.Name,
		CreatedAt:  timestamppb.New(comment.CreatedAt),
		UpdatedAt:  timestamppb.New(comment.UpdatedAt),
		ParentId:   req.ParentId,
		Depth:      int32(depth),
		Status:     boardServer.CommentStatus_COMMENT_STATUS_PUBLISHED,
	}, nil
}

// CreatePosting implements boardv1.BoardServiceServer.
func (s *Service) CreatePosting(ctx context.Context, req *boardServer.CreatePostingRequest) (*boardServer.Posting, error) {
	authUser, err := s.authenticator.FromGrpcContext(ctx)
	if err != nil {
		return nil, err
	}

	posting := postings.Posting{
		Title:     req.Title,
		Content:   req.Content,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	if err := s.postingRepository.Save(ctx, &posting); err != nil {
		return nil, status.Errorf(codes.Internal, "failed to create posting: %v", err)
	}

	return &boardServer.Posting{
		PostingId:    int64(posting.ID),
		Title:        posting.Title,
		Content:      posting.Content,
		AuthorId:     authUser.ID,
		AuthorName:   authUser.Name,
		CreatedAt:    timestamppb.New(posting.CreatedAt),
		UpdatedAt:    timestamppb.New(posting.UpdatedAt),
		ViewCount:    int32(posting.ViewCount),
		CommentCount: int32(posting.CommentCount),
		Tags:         posting.Tags,
		Status:       boardServer.PostingStatus_POSTING_STATUS_PUBLISHED,
	}, nil
}

// DeleteComment implements boardv1.BoardServiceServer.
func (s *Service) DeleteComment(ctx context.Context, req *boardServer.DeleteCommentRequest) (*emptypb.Empty, error) {
	authUser, err := s.authenticator.FromGrpcContext(ctx)
	if err != nil {
		return nil, err
	}

	comment, err := s.commentRepository.GetByID(ctx, req.CommentId)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to get comment: %v", err)
	}

	if comment.AuthorID != authUser.ID {
		return nil, status.Errorf(codes.PermissionDenied, "you are not allowed to delete this comment")
	}

	if err := s.commentRepository.Delete(ctx, req.CommentId); err != nil {
		return nil, status.Errorf(codes.Internal, "failed to delete comment: %v", err)
	}

	return &emptypb.Empty{}, nil
}

// DeletePosting implements boardv1.BoardServiceServer.
func (s *Service) DeletePosting(ctx context.Context, req *boardServer.DeletePostingRequest) (*emptypb.Empty, error) {
	authUser, err := s.authenticator.FromGrpcContext(ctx)
	if err != nil {
		return nil, err
	}

	posting, err := s.postingRepository.FindOneByID(ctx, int(req.PostingId))
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to get posting: %v", err)
	}

	if posting.AuthorID != authUser.ID {
		return nil, status.Errorf(codes.PermissionDenied, "you are not allowed to delete this posting")
	}

	if err := s.postingRepository.Delete(ctx, int(req.PostingId)); err != nil {
		return nil, status.Errorf(codes.Internal, "failed to delete posting: %v", err)
	}

	return &emptypb.Empty{}, nil
}

// GetPosting implements boardv1.BoardServiceServer.
func (s *Service) GetPosting(ctx context.Context, req *boardServer.GetPostingRequest) (*boardServer.Posting, error) {
	// 인증된 사용자 정보 가져오기 (익명 사용자도 허용)
	var userID string
	authUser, err := s.authenticator.FromGrpcContext(ctx)
	if err == nil && authUser != nil {
		userID = authUser.ID
	} else {
		// 익명 사용자의 경우 IP 주소 활용 (HTTP 헤더에서 가져오기)
		md, ok := metadata.FromIncomingContext(ctx)
		if ok {
			if values := md.Get("x-forwarded-for"); len(values) > 0 {
				userID = "anon:" + values[0]
			} else if values := md.Get("x-real-ip"); len(values) > 0 {
				userID = "anon:" + values[0]
			}
		}
	}

	// 캐시 키 생성 (Redis 또는 메모리 캐시 사용)
	viewCountKey := fmt.Sprintf("viewcount:%d:%s", req.PostingId, userID)

	// 캐시 확인
	exists, err := s.cache.Exists(ctx, viewCountKey)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to check cache: %v", err)
	}

	if !exists {
		// 조회수 증가
		if err := s.postingRepository.IncrementViewCount(ctx, int(req.PostingId)); err != nil {
			zap.S().Errorf("failed to increment view count: %v", err)
		}

		// 캐시에 기록 (24시간 유효)
		s.cache.Set(ctx, viewCountKey, "1", 24*time.Hour)
	}

	// 원래 코드: 게시글 조회
	posting, err := s.postingRepository.FindOneByID(ctx, int(req.PostingId))
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to get posting: %v", err)
	}

	comments, err := s.commentRepository.GetByPostID(ctx, int64(req.PostingId))
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to get comments: %v", err)
	}

	boardComments := make([]*boardServer.Comment, len(comments))
	for i, comment := range comments {
		depth := 0
		if comment.ParentID != nil {
			depth = 1
		}
		parentID := int64(0)
		if comment.ParentID != nil {
			parentID = *comment.ParentID
		}
		boardComments[i] = &boardServer.Comment{
			CommentId:  comment.ID,
			PostingId:  comment.PostID,
			Content:    comment.Content,
			AuthorId:   comment.AuthorID,
			AuthorName: comment.AuthorName,
			CreatedAt:  timestamppb.New(comment.CreatedAt),
			UpdatedAt:  timestamppb.New(comment.UpdatedAt),
			ParentId:   parentID,
			Depth:      int32(depth),
			Status:     comment.Status,
		}
	}

	return &boardServer.Posting{
		PostingId:    int64(posting.ID),
		Title:        posting.Title,
		Content:      posting.Content,
		AuthorId:     posting.AuthorID,
		AuthorName:   posting.AuthorName,
		CreatedAt:    timestamppb.New(posting.CreatedAt),
		UpdatedAt:    timestamppb.New(posting.UpdatedAt),
		ViewCount:    int32(posting.ViewCount),
		CommentCount: int32(posting.CommentCount),
		Tags:         posting.Tags,
		Status:       boardServer.PostingStatus_POSTING_STATUS_PUBLISHED,
		Comments:     boardComments,
	}, nil
}

// ListCommentsByPosting implements boardv1.BoardServiceServer.
func (s *Service) ListCommentsByPosting(ctx context.Context, req *boardServer.ListCommentsByPostingRequest) (*boardServer.ListCommentsResponse, error) {
	comments, err := s.commentRepository.GetByPostID(ctx, req.PostingId)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to get comments: %v", err)
	}

	boardComments := make([]*boardServer.Comment, len(comments))
	for i, comment := range comments {
		depth := 0
		if comment.ParentID != nil {
			depth = 1
		}
		boardComments[i] = &boardServer.Comment{
			CommentId:  comment.ID,
			PostingId:  comment.PostID,
			Content:    comment.Content,
			AuthorId:   comment.AuthorID,
			AuthorName: comment.AuthorName,
			CreatedAt:  timestamppb.New(comment.CreatedAt),
			UpdatedAt:  timestamppb.New(comment.UpdatedAt),
			ParentId:   *comment.ParentID,
			Depth:      int32(depth),
			Status:     comment.Status,
		}
	}
	return &boardServer.ListCommentsResponse{Comments: boardComments}, nil
}

// ListPostings implements boardv1.BoardServiceServer.
func (s *Service) ListPostings(ctx context.Context, req *boardServer.ListPostingsRequest) (*boardServer.ListPostingsResponse, error) {
	token, err := pagination.FromEncodedString(req.PageToken)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid token: %v", err)
	}

	postings, err := s.postingRepository.FindAll(ctx, &db.CursorBasedPagination{
		Cursor: token.Cursor,
		Limit:  token.Limit,
	})
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to get postings: %v", err)
	}

	boardPostings := make([]*boardServer.Posting, len(postings.Data))
	for i, posting := range postings.Data {
		boardPostings[i] = &boardServer.Posting{
			PostingId:    int64(posting.ID),
			Title:        posting.Title,
			Content:      posting.Content,
			AuthorId:     posting.AuthorID,
			AuthorName:   posting.AuthorName,
			CreatedAt:    timestamppb.New(posting.CreatedAt),
			UpdatedAt:    timestamppb.New(posting.UpdatedAt),
			ViewCount:    int32(posting.ViewCount),
			CommentCount: int32(posting.CommentCount),
			Tags:         posting.Tags,
			Status:       boardServer.PostingStatus_POSTING_STATUS_PUBLISHED,
		}
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
		return nil, status.Errorf(codes.Internal, "failed to get comment: %v", err)
	}

	if comment.AuthorID != authUser.ID {
		return nil, status.Errorf(codes.PermissionDenied, "you are not allowed to update this comment")
	}

	comment.Content = req.Content
	comment.UpdatedAt = time.Now()

	if err := s.commentRepository.Save(ctx, comment); err != nil {
		return nil, status.Errorf(codes.Internal, "failed to update comment: %v", err)
	}

	depth := 0
	if comment.ParentID != nil {
		depth = 1
	}

	parentID := int64(0)
	if comment.ParentID != nil {
		parentID = *comment.ParentID
	}

	return &boardServer.Comment{
		CommentId:  comment.ID,
		PostingId:  comment.PostID,
		Content:    comment.Content,
		AuthorId:   comment.AuthorID,
		AuthorName: comment.AuthorName,
		CreatedAt:  timestamppb.New(comment.CreatedAt),
		UpdatedAt:  timestamppb.New(comment.UpdatedAt),
		ParentId:   parentID,
		Depth:      int32(depth),
		Status:     comment.Status,
	}, nil
}

// UpdatePosting implements boardv1.BoardServiceServer.
func (s *Service) UpdatePosting(ctx context.Context, req *boardServer.UpdatePostingRequest) (*boardServer.Posting, error) {
	authUser, err := s.authenticator.FromGrpcContext(ctx)
	if err != nil {
		return nil, err
	}

	posting, err := s.postingRepository.FindOneByID(ctx, int(req.PostingId))
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to get posting: %v", err)
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
		return nil, status.Errorf(codes.Internal, "failed to update posting: %v", err)
	}

	return &boardServer.Posting{
		PostingId:    int64(posting.ID),
		Title:        posting.Title,
		Content:      posting.Content,
		AuthorId:     posting.AuthorID,
		AuthorName:   posting.AuthorName,
		CreatedAt:    timestamppb.New(posting.CreatedAt),
		UpdatedAt:    timestamppb.New(posting.UpdatedAt),
		ViewCount:    int32(posting.ViewCount),
		CommentCount: int32(posting.CommentCount),
		Tags:         posting.Tags,
		Status:       posting.Status,
	}, nil
}

// mustEmbedUnimplementedBoardServiceServer implements boardv1.BoardServiceServer.
func (s *Service) mustEmbedUnimplementedBoardServiceServer() {
	panic("unimplemented")
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

func WithCache(cache cache.Cache) ServiceOption {
	return func(s *Service) {
		s.cache = cache
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

	if s.cache == nil {
		return nil, errors.New("cache is required")
	}

	return s, nil
}
