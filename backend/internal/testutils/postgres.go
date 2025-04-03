package testutils

import (
	"context"
	"os"
	"time"

	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
	"github.com/testcontainers/testcontainers-go/wait"
)

type PostgresTestContainer struct {
	container *postgres.PostgresContainer
}

type PostgresTestContainerConfig struct {
	dbname   string
	user     string
	password string
}

type PostgresTestContainerOption func(*PostgresTestContainerConfig)

func WithDBName(dbname string) PostgresTestContainerOption {
	return func(c *PostgresTestContainerConfig) {
		c.dbname = dbname
	}
}

func WithUser(user string) PostgresTestContainerOption {
	return func(c *PostgresTestContainerConfig) {
		c.user = user
	}
}

func WithPassword(password string) PostgresTestContainerOption {
	return func(c *PostgresTestContainerConfig) {
		c.password = password
	}
}

var defaultPostgresTestContainerConfig = PostgresTestContainerConfig{
	dbname:   "postgres",
	user:     "postgres",
	password: "postgres",
}

func NewPostgresTestContainer(ctx context.Context, opts ...PostgresTestContainerOption) (*PostgresTestContainer, error) {
	os.Setenv("TESTCONTAINERS_RYUK_DISABLED", "true")
	cfg := defaultPostgresTestContainerConfig
	for _, opt := range opts {
		opt(&cfg)
	}

	ctr, err := postgres.Run(ctx,
		"postgres:16-alpine",
		postgres.WithDatabase(cfg.dbname),
		postgres.WithUsername(cfg.user),
		postgres.WithPassword(cfg.password),
		testcontainers.WithWaitStrategy(
			wait.ForLog("database system is ready to accept connections").
				WithOccurrence(2).
				WithStartupTimeout(5*time.Second)))
	if err != nil {
		return nil, err
	}

	return &PostgresTestContainer{container: ctr}, nil
}

func (c *PostgresTestContainer) Close() error {
	return c.container.Terminate(context.Background())
}

func (c *PostgresTestContainer) Endpoint() (string, error) {
	connectionString, err := c.container.ConnectionString(context.Background(), "")
	if err != nil {
		return "", err
	}

	return connectionString, nil
}
