package grpcutils

import (
	"context"
	"net"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/health"
	"google.golang.org/grpc/health/grpc_health_v1"
	"google.golang.org/grpc/keepalive"
	"google.golang.org/grpc/reflection"
)

// GrpcServer is a gRPC server.
//
// Fields:
//   - grpcServer: the gRPC server
//   - listener: the listener to listen on
type GrpcServer struct {
	grpcServer        *grpc.Server
	grpcServerOptions []grpc.ServerOption
	services          []Service
	listener          net.Listener
}

func (s *GrpcServer) RegisterService(service *grpc.ServiceDesc, impl any) {
	s.grpcServer.RegisterService(service, impl)
}

// Service is a service to register on the server.
//
// Fields:
//   - service: the service to register
//   - impl: the implementation of the service
type Service struct {
	Desc *grpc.ServiceDesc
	Impl any
}

func NewService(service *grpc.ServiceDesc, impl any) Service {
	return Service{Desc: service, Impl: impl}
}

// defaultServerOptions are the default server options.
var defaultServerOptions = []grpc.ServerOption{
	grpc.ChainUnaryInterceptor(
	// TODO: add interceptors here
	),
	grpc.KeepaliveParams(keepalive.ServerParameters{
		MaxConnectionIdle: 180 * time.Second,
		Time:              10 * time.Second,
		Timeout:           5 * time.Second,
	}),
	grpc.KeepaliveEnforcementPolicy(keepalive.EnforcementPolicy{
		MinTime: 10 * time.Second,
	}),
}

type ServerOption func(*GrpcServer)

func WithListener(listener net.Listener) ServerOption {
	return func(s *GrpcServer) {
		s.listener = listener
	}
}

func WithGrpcServerOptions(opts ...grpc.ServerOption) ServerOption {
	return func(s *GrpcServer) {
		s.grpcServerOptions = append(s.grpcServerOptions, opts...)
	}
}

func WithServices(srvs ...Service) ServerOption {
	return func(s *GrpcServer) {
		s.services = append(s.services, srvs...)
	}
}

// NewGrpcServer creates a new gRPC server.
//
// Parameters:
//   - opts: the options to pass to the server
//
// Returns:
//   - a new Server instance
//   - an error if the server fails to start
//
// Details:
//   - The server will register the health service and the reflection service
//   - The server will use the default server options
//   - The server will register the services passed in the srvs parameter
//   - The caller is responsible for calling the Start method to start the server
func NewGrpcServer(opts ...ServerOption) (*GrpcServer, error) {
	serverOptions := defaultServerOptions

	server := &GrpcServer{}
	for _, opt := range opts {
		opt(server)
	}

	grpcServer := grpc.NewServer(serverOptions...)
	grpc_health_v1.RegisterHealthServer(grpcServer, health.NewServer())
	reflection.Register(grpcServer)

	for _, service := range server.services {
		grpcServer.RegisterService(service.Desc, service.Impl)
	}

	server.grpcServer = grpcServer
	return server, nil
}

// Start starts the server.
//
// Returns:
//   - an error if the server fails to start
//
// Details:
//   - The server will start serving on the listener passed in the NewServer function
func (s *GrpcServer) Start(ctx context.Context) error {
	return s.grpcServer.Serve(s.listener)
}

// Stop stops the server.
//
// Returns:
//   - an error if the server fails to stop
//
// Details:
//   - The server will stop serving on the listener passed in the NewServer function
//   - The caller is responsible for calling the Close method to close the listener
func (s *GrpcServer) Stop(ctx context.Context) error {
	s.grpcServer.Stop()
	return nil
}
