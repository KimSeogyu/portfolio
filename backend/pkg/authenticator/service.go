package authenticator

import (
	"context"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

type AuthUser struct {
	ID    string
	Email string
	Name  string
}

// fromGrpcContext는 gRPC 컨텍스트에서 인증 사용자 정보를 추출합니다.
// 컨텍스트에 인증 사용자 정보가 없으면 오류를 반환합니다.
// TODO: grpc server에 authuser를 세팅하는 작업 해야 함
func fromGrpcContext(ctx context.Context) (*AuthUser, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return nil, status.Errorf(codes.Unauthenticated, "no auth user in context")
	}

	userID := md.Get("auth_user_id")
	if len(userID) == 0 {
		return nil, status.Errorf(codes.Unauthenticated, "no auth user ID in context")
	}

	userEmail := md.Get("auth_user_email")
	if len(userEmail) == 0 {
		return nil, status.Errorf(codes.Unauthenticated, "no auth user email in context")
	}

	userName := md.Get("auth_user_name")
	if len(userName) == 0 {
		return nil, status.Errorf(codes.Unauthenticated, "no auth user name in context")
	}

	return &AuthUser{
		ID:    userID[0],
		Email: userEmail[0],
		Name:  userName[0],
	}, nil
}

// 인터페이스 정의
type UserAuthenticator interface {
	FromGrpcContext(ctx context.Context) (*AuthUser, error)
}

// 실제 구현
type RealAuthenticator struct{}

func NewRealAuthenticator() *RealAuthenticator {
	return &RealAuthenticator{}
}

func (a *RealAuthenticator) FromGrpcContext(ctx context.Context) (*AuthUser, error) {
	// 실제 로직
	return fromGrpcContext(ctx)
}

// 테스트용 구현
type TestAuthenticator struct{}

func (a *TestAuthenticator) FromGrpcContext(ctx context.Context) (*AuthUser, error) {
	return &AuthUser{ID: "test-user-id", Name: "Test User"}, nil
}

func NewTestAuthenticator() *TestAuthenticator {
	return &TestAuthenticator{}
}
