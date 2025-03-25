package grpcutils

import (
	"context"
	"net"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc/health/grpc_health_v1"
)

func TestClient(t *testing.T) {
	serverAddr := ":50051"
	listener, err := net.Listen("tcp", serverAddr)
	if err != nil {
		t.Fatalf("Failed to create listener: %v", err)
	}
	server, err := NewGrpcServer(WithListener(listener), WithGrpcServerOptions(), WithServices())
	if err != nil {
		t.Fatalf("Failed to create server: %v", err)
	}
	ctx := context.Background()
	go server.Start(ctx)

	t.Cleanup(func() {
		server.Stop(ctx)
	})

	time.Sleep(1 * time.Second)

	client, err := NewGrpcClient(serverAddr)
	if err != nil {
		t.Fatalf("Failed to create client: %v", err)
	}

	t.Cleanup(func() {
		client.Close()
	})

	cl := grpc_health_v1.NewHealthClient(client)

	resp, err := cl.Check(context.Background(), &grpc_health_v1.HealthCheckRequest{
		Service: "",
	})
	if err != nil {
		t.Fatalf("Failed to check health: %v", err)
	}

	assert.Equal(t, grpc_health_v1.HealthCheckResponse_SERVING, resp.Status)
}
