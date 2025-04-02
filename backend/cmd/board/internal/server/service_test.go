package server

import (
	"context"
	"testing"
	"time"

	"github.com/kimseogyu/portfolio/backend/cmd/board/internal/comments"
	mockComments "github.com/kimseogyu/portfolio/backend/cmd/board/internal/comments/mocks"
	"github.com/kimseogyu/portfolio/backend/cmd/board/internal/postings"
	mockPostings "github.com/kimseogyu/portfolio/backend/cmd/board/internal/postings/mocks"
	mockViewCountManager "github.com/kimseogyu/portfolio/backend/cmd/board/internal/viewcount/mocks"
	mockCStore "github.com/kimseogyu/portfolio/backend/internal/cstore/mocks"
	"github.com/kimseogyu/portfolio/backend/internal/db"
	mockDLockerFactory "github.com/kimseogyu/portfolio/backend/internal/dlock/mocks"
	boardServer "github.com/kimseogyu/portfolio/backend/internal/proto/board/v1"
	"github.com/kimseogyu/portfolio/backend/pkg/authenticator"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
	"google.golang.org/grpc/metadata"
	"google.golang.org/protobuf/types/known/emptypb"
)

// 테스트를 위한 모의 객체와 서비스를 생성하는 구조체
type testFixture struct {
	ctrl                 *gomock.Controller
	mockCommentRepo      *mockComments.MockRepository
	mockPostingRepo      *mockPostings.MockRepository
	mockCache            *mockCStore.MockCacheStore
	mockViewCountManager *mockViewCountManager.MockViewCountManager
	service              *Service
}

// 테스트 픽스처 생성 함수
func newTestFixture(t *testing.T) *testFixture {
	ctrl := gomock.NewController(t)

	mockCommentRepo := mockComments.NewMockRepository(ctrl)
	mockPostingRepo := mockPostings.NewMockRepository(ctrl)
	mockCache := mockCStore.NewMockCacheStore(ctrl)
	mockDLockerFactory := mockDLockerFactory.NewMockDLockerFactory(ctrl)
	mockViewCountManager := mockViewCountManager.NewMockViewCountManager(ctrl)
	service, err := NewService(
		WithCommentRepository(mockCommentRepo),
		WithPostingRepository(mockPostingRepo),
		WithAuthenticator(&authenticator.TestAuthenticator{}),
		WithCacheStore(mockCache),
		WithDLockerFactory(mockDLockerFactory),
		WithViewCountManager(mockViewCountManager),
	)
	require.NoError(t, err)

	return &testFixture{
		ctrl:                 ctrl,
		mockCommentRepo:      mockCommentRepo,
		mockPostingRepo:      mockPostingRepo,
		mockCache:            mockCache,
		mockViewCountManager: mockViewCountManager,
		service:              service,
	}
}

// 테스트용 컨텍스트 생성 함수
func createAuthContext() context.Context {
	ctx := context.Background()
	md := metadata.New(map[string]string{
		"Authorization": "Bearer test-token",
	})
	return metadata.NewIncomingContext(ctx, md)
}

func TestCreateComment(t *testing.T) {
	// 테스트 픽스처 생성
	fixture := newTestFixture(t)
	defer fixture.ctrl.Finish()

	// Save 메서드 호출 시 성공 반환 설정
	fixture.mockCommentRepo.EXPECT().
		Save(gomock.Any(), gomock.AssignableToTypeOf(&comments.Comment{})).
		DoAndReturn(func(ctx context.Context, comment *comments.Comment) error {
			comment.ID = 123 // ID 할당
			return nil
		})

	// 테스트 요청 생성
	req := &boardServer.CreateCommentRequest{
		PostingId: 456,
		Content:   "Test comment content",
		ParentId:  0, // 최상위 댓글
	}

	// 인증된 컨텍스트 생성
	ctx := createAuthContext()

	// 메서드 호출
	response, err := fixture.service.CreateComment(ctx, req)

	// 검증
	assert.NoError(t, err)
	assert.NotNil(t, response)
	assert.Equal(t, int64(123), response.CommentId)
	assert.Equal(t, req.Content, response.Content)
	assert.Equal(t, req.PostingId, response.PostingId)
	assert.Equal(t, "test-user-id", response.AuthorId)
	assert.Equal(t, "Test User", response.AuthorName)
	assert.Equal(t, boardServer.CommentStatus_COMMENT_STATUS_PUBLISHED, response.Status)
}

func TestCreatePosting(t *testing.T) {
	// 테스트 픽스처 생성
	fixture := newTestFixture(t)
	defer fixture.ctrl.Finish()

	// Save 메서드 호출 시 성공 반환 설정
	fixture.mockPostingRepo.EXPECT().
		Save(gomock.Any(), gomock.AssignableToTypeOf(&postings.Posting{})).
		DoAndReturn(func(ctx context.Context, posting *postings.Posting) error {
			posting.ID = 789 // ID 할당
			return nil
		})

	// 테스트 요청 생성
	req := &boardServer.CreatePostingRequest{
		Title:   "Test Title",
		Content: "Test content for posting",
		Tags:    []string{"test", "example"},
		Status:  boardServer.PostingStatus_POSTING_STATUS_PUBLISHED,
	}

	// 인증된 컨텍스트 생성
	ctx := createAuthContext()

	// 메서드 호출
	response, err := fixture.service.CreatePosting(ctx, req)

	// 검증
	assert.NoError(t, err)
	assert.NotNil(t, response)
	assert.Equal(t, int64(789), response.PostingId)
	assert.Equal(t, req.Title, response.Title)
	assert.Equal(t, req.Content, response.Content)
	assert.Equal(t, "test-user-id", response.AuthorId)
	assert.Equal(t, boardServer.PostingStatus_POSTING_STATUS_PUBLISHED, response.Status)
}

func TestGetPosting(t *testing.T) {
	// 테스트 픽스처 생성
	fixture := newTestFixture(t)
	defer fixture.ctrl.Finish()

	// 테스트용 시간
	now := time.Now()

	fixture.mockViewCountManager.EXPECT().CheckAndIncrement(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
		Return(nil)

	// FindOneByID 메서드 호출 시 게시물 반환 설정
	fixture.mockPostingRepo.EXPECT().
		FindOneByID(gomock.Any(), int64(789)).
		Return(&postings.Posting{
			ID:           789,
			Title:        "Test Posting",
			Content:      "Test content",
			AuthorID:     "test-user-id",
			AuthorName:   "Test User",
			CommentCount: 2,
			ViewCount:    10,
			CreatedAt:    now,
			UpdatedAt:    now,
			Status:       boardServer.PostingStatus_POSTING_STATUS_PUBLISHED,
		}, nil)

	// GetByPostID 메서드 호출 시 댓글 반환 설정
	parentID := int64(100)
	fixture.mockCommentRepo.EXPECT().
		GetByPostID(gomock.Any(), int64(789)).
		Return([]comments.Comment{
			{
				ID:         101,
				PostID:     789,
				Content:    "Comment 1",
				AuthorID:   "user1",
				AuthorName: "User 1",
				CreatedAt:  now,
				UpdatedAt:  now,
				ParentID:   nil,
				Status:     boardServer.CommentStatus_COMMENT_STATUS_PUBLISHED,
			},
			{
				ID:         102,
				PostID:     789,
				Content:    "Reply to comment 1",
				AuthorID:   "user2",
				AuthorName: "User 2",
				CreatedAt:  now,
				UpdatedAt:  now,
				ParentID:   &parentID,
				Status:     boardServer.CommentStatus_COMMENT_STATUS_PUBLISHED,
			},
		}, nil)

	// 테스트 요청 생성
	req := &boardServer.GetPostingRequest{
		PostingId: 789,
	}

	fixture.mockPostingRepo.EXPECT().
		IncrementViewCount(gomock.Any(), int64(789)).
		Return(nil).
		AnyTimes()

	// 컨텍스트 생성
	ctx := context.Background()

	// 메서드 호출
	response, err := fixture.service.GetPosting(ctx, req)

	// 검증
	assert.NoError(t, err)
	assert.NotNil(t, response)
	assert.Equal(t, int64(789), response.PostingId)
	assert.Equal(t, "Test Posting", response.Title)
	assert.Equal(t, int32(2), response.CommentCount)
	assert.Equal(t, 2, len(response.Comments))
}

func TestDeletePosting(t *testing.T) {
	// 테스트 픽스처 생성
	fixture := newTestFixture(t)
	defer fixture.ctrl.Finish()

	// FindOneByID 메서드 호출 시 게시물 반환 설정
	fixture.mockPostingRepo.EXPECT().
		FindOneByID(gomock.Any(), int64(789)).
		Return(&postings.Posting{
			ID:         789,
			Title:      "Test Posting",
			Content:    "Test content",
			AuthorID:   "test-user-id",
			AuthorName: "Test User",
		}, nil)

	// Delete 메서드 호출 시 성공 반환 설정
	fixture.mockPostingRepo.EXPECT().
		Delete(gomock.Any(), int64(789)).
		Return(nil)

	// 테스트 요청 생성
	req := &boardServer.DeletePostingRequest{
		PostingId: 789,
	}

	// 인증된 컨텍스트 생성
	ctx := createAuthContext()

	// 메서드 호출
	response, err := fixture.service.DeletePosting(ctx, req)

	// 검증
	assert.NoError(t, err)
	assert.NotNil(t, response)
	assert.IsType(t, &emptypb.Empty{}, response)
}

func TestDeleteComment(t *testing.T) {
	// 테스트 픽스처 생성
	fixture := newTestFixture(t)
	defer fixture.ctrl.Finish()

	// GetByID 메서드 호출 시 댓글 반환 설정
	fixture.mockCommentRepo.EXPECT().
		GetByID(gomock.Any(), int64(123)).
		Return(&comments.Comment{
			ID:         123,
			PostID:     789,
			Content:    "Test Comment",
			AuthorID:   "test-user-id",
			AuthorName: "Test User",
		}, nil)

	// Delete 메서드 호출 시 성공 반환 설정
	fixture.mockCommentRepo.EXPECT().
		Delete(gomock.Any(), int64(123)).
		Return(nil)

	// 테스트 요청 생성
	req := &boardServer.DeleteCommentRequest{
		CommentId: 123,
		PostingId: 789,
	}

	// 인증된 컨텍스트 생성
	ctx := createAuthContext()

	// 메서드 호출
	response, err := fixture.service.DeleteComment(ctx, req)

	// 검증
	assert.NoError(t, err)
	assert.NotNil(t, response)
	assert.IsType(t, &emptypb.Empty{}, response)
}

func TestUpdateComment(t *testing.T) {
	// 테스트 픽스처 생성
	fixture := newTestFixture(t)
	defer fixture.ctrl.Finish()

	// 테스트용 시간
	now := time.Now()

	// GetByID 메서드 호출 시 댓글 반환 설정
	fixture.mockCommentRepo.EXPECT().
		GetByID(gomock.Any(), int64(123)).
		Return(&comments.Comment{
			ID:         123,
			PostID:     789,
			Content:    "Original comment",
			AuthorID:   "test-user-id",
			AuthorName: "Test User",
			CreatedAt:  now,
			UpdatedAt:  now,
		}, nil)

	// Save 메서드 호출 시 성공 반환 설정
	fixture.mockCommentRepo.EXPECT().
		Save(gomock.Any(), gomock.Any()).
		DoAndReturn(func(ctx context.Context, comment *comments.Comment) error {
			assert.Equal(t, int64(123), comment.ID)
			assert.Equal(t, "Updated comment", comment.Content)
			return nil
		})

	// 테스트 요청 생성
	req := &boardServer.UpdateCommentRequest{
		CommentId: 123,
		PostingId: 789,
		Content:   "Updated comment",
	}

	// 인증된 컨텍스트 생성
	ctx := createAuthContext()

	// 메서드 호출
	response, err := fixture.service.UpdateComment(ctx, req)

	// 검증
	assert.NoError(t, err)
	assert.NotNil(t, response)
	assert.Equal(t, int64(123), response.CommentId)
	assert.Equal(t, "Updated comment", response.Content)
}

func TestUpdatePosting(t *testing.T) {
	// 테스트 픽스처 생성
	fixture := newTestFixture(t)
	defer fixture.ctrl.Finish()

	// 테스트용 시간
	now := time.Now()

	// FindOneByID 메서드 호출 시 게시물 반환 설정
	fixture.mockPostingRepo.EXPECT().
		FindOneByID(gomock.Any(), int64(789)).
		Return(&postings.Posting{
			ID:         789,
			Title:      "Original title",
			Content:    "Original content",
			AuthorID:   "test-user-id",
			AuthorName: "Test User",
			CreatedAt:  now,
			UpdatedAt:  now,
		}, nil)

	// Save 메서드 호출 시 성공 반환 설정
	fixture.mockPostingRepo.EXPECT().
		Save(gomock.Any(), gomock.Any()).
		DoAndReturn(func(ctx context.Context, posting *postings.Posting) error {
			assert.Equal(t, int64(789), posting.ID)
			assert.Equal(t, "Updated title", posting.Title)
			assert.Equal(t, "Updated content", posting.Content)
			return nil
		})

	// 테스트 요청 생성
	req := &boardServer.UpdatePostingRequest{
		PostingId: 789,
		Title:     "Updated title",
		Content:   "Updated content",
		Tags:      []string{"updated", "test"},
		Status:    boardServer.PostingStatus_POSTING_STATUS_PUBLISHED,
	}

	// 인증된 컨텍스트 생성
	ctx := createAuthContext()

	// 메서드 호출
	response, err := fixture.service.UpdatePosting(ctx, req)

	// 검증
	assert.NoError(t, err)
	assert.NotNil(t, response)
	assert.Equal(t, int64(789), response.PostingId)
	assert.Equal(t, "Updated title", response.Title)
	assert.Equal(t, "Updated content", response.Content)
	assert.Equal(t, []string{"updated", "test"}, response.Tags)
}

func TestListPostings(t *testing.T) {
	// 테스트 픽스처 생성
	fixture := newTestFixture(t)
	defer fixture.ctrl.Finish()

	// 테스트용 시간
	now := time.Now()

	// FindAll 메서드 호출 시 게시물 목록 반환 설정
	fixture.mockPostingRepo.EXPECT().
		FindAll(gomock.Any(), gomock.Any()).
		Return(&db.CursorBasedPaginationResponse[postings.Posting]{
			Data: []postings.Posting{
				{
					ID:         1,
					Title:      "Posting 1",
					Content:    "Content 1",
					AuthorID:   "user1",
					AuthorName: "User 1",
					CreatedAt:  now,
					UpdatedAt:  now,
				},
				{
					ID:         2,
					Title:      "Posting 2",
					Content:    "Content 2",
					AuthorID:   "user2",
					AuthorName: "User 2",
					CreatedAt:  now,
					UpdatedAt:  now,
				},
			},
			NextCursor: nil,
			HasNext:    false,
			Total:      2,
		}, nil)

	// 테스트 요청 생성 - 페이지 토큰을 nil로 설정
	req := &boardServer.ListPostingsRequest{
		PageSize:  10,
		PageToken: "",
	}

	// 인증된 컨텍스트 생성 - 일반 컨텍스트 대신 인증 컨텍스트 사용
	ctx := createAuthContext()

	// 메서드 호출
	response, err := fixture.service.ListPostings(ctx, req)

	// 검증
	assert.NoError(t, err)
	assert.NotNil(t, response)
	assert.Equal(t, 2, len(response.Postings))
	assert.Equal(t, int32(2), response.TotalCount)
}

// IP 주소가 있는 컨텍스트 생성 헬퍼 함수
func createContextWithIP(ip string) context.Context {
	ctx := context.Background()
	md := metadata.New(map[string]string{
		"x-forwarded-for": ip,
		"x-real-ip":       ip, // 실제 IP 주소도 추가
	})
	return metadata.NewIncomingContext(ctx, md)
}

// 익명 사용자용 모의 인증기
type AnonymousAuthenticator struct{}

func (a *AnonymousAuthenticator) FromGrpcContext(ctx context.Context) (*authenticator.AuthUser, error) {
	// 컨텍스트에서 IP 주소 추출
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return nil, nil
	}

	var ip string
	if ips := md.Get("x-real-ip"); len(ips) > 0 {
		ip = ips[0]
	} else if ips := md.Get("x-forwarded-for"); len(ips) > 0 {
		ip = ips[0]
	} else {
		ip = "127.0.0.1"
	}

	return &authenticator.AuthUser{ID: "anon:" + ip, Name: "Anonymous"}, nil
}

// 익명 사용자용 테스트 픽스처 생성
func newAnonymousTestFixture(t *testing.T) *testFixture {
	ctrl := gomock.NewController(t)

	mockCommentRepo := mockComments.NewMockRepository(ctrl)
	mockPostingRepo := mockPostings.NewMockRepository(ctrl)
	mockCache := mockCStore.NewMockCacheStore(ctrl)
	mockDLockerFactory := mockDLockerFactory.NewMockDLockerFactory(ctrl)
	mockViewCountManager := mockViewCountManager.NewMockViewCountManager(ctrl)
	service, err := NewService(
		WithCommentRepository(mockCommentRepo),
		WithPostingRepository(mockPostingRepo),
		WithAuthenticator(&AnonymousAuthenticator{}),
		WithCacheStore(mockCache),
		WithDLockerFactory(mockDLockerFactory),
		WithViewCountManager(mockViewCountManager),
	)
	require.NoError(t, err)

	return &testFixture{
		ctrl:                 ctrl,
		mockCommentRepo:      mockCommentRepo,
		mockPostingRepo:      mockPostingRepo,
		mockCache:            mockCache,
		mockViewCountManager: mockViewCountManager,
		service:              service,
	}
}
