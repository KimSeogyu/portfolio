package testutils

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
)

func NewRedisTestContainer(ctx context.Context) (testcontainers.Container, error) {
	os.Setenv("TESTCONTAINERS_RYUK_DISABLED", "true")
	req := testcontainers.ContainerRequest{
		Image:        "redis:latest",
		ExposedPorts: []string{"6379/tcp"},
		WaitingFor:   wait.ForLog("Ready to accept connections"),
	}
	redisC, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
	})
	if err != nil {
		return nil, err
	}

	// Clean up the container
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT, syscall.SIGKILL)
	go func() {
		<-sigs
		err := redisC.Terminate(context.Background())
		if err != nil {
			fmt.Println(err)
		}
	}()

	return redisC, nil
}
