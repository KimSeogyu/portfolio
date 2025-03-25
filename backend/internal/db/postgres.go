package db

import (
	"errors"
	"fmt"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type PostgresConfig struct {
	host     string
	port     int
	user     string
	password string
	dbname   string
	schema   string
	dsn      string
	initSQL  string
}

type PostgresOption func(*PostgresConfig)

func WithDSN(dsn string) PostgresOption {
	return func(c *PostgresConfig) {
		c.dsn = dsn
	}
}

func WithHost(host string) PostgresOption {
	return func(c *PostgresConfig) {
		c.host = host
	}
}

func WithPort(port int) PostgresOption {
	return func(c *PostgresConfig) {
		c.port = port
	}
}

func WithUser(user string) PostgresOption {
	return func(c *PostgresConfig) {
		c.user = user
	}
}

func WithPassword(password string) PostgresOption {
	return func(c *PostgresConfig) {
		c.password = password
	}
}

func WithDBName(dbname string) PostgresOption {
	return func(c *PostgresConfig) {
		c.dbname = dbname
	}
}

func WithSchema(schema string) PostgresOption {
	return func(c *PostgresConfig) {
		c.schema = schema
	}
}

func WithInitSQL(initSQL string) PostgresOption {
	return func(c *PostgresConfig) {
		c.initSQL = initSQL
	}
}

var defaultPostgresConfig = PostgresConfig{
	host:     "localhost",
	port:     5432,
	user:     "postgres",
	password: "postgres",
	dbname:   "postgres",
	schema:   "public",
	initSQL:  "",
}

func NewPostgresDB(opts ...PostgresOption) (*gorm.DB, error) {
	cfg := defaultPostgresConfig
	for _, opt := range opts {
		opt(&cfg)
	}

	dsn := cfg.dsn
	if dsn == "" {
		if cfg.host == "" {
			return nil, errors.New("host is required")
		}
		if cfg.port == 0 {
			return nil, errors.New("port is required")
		}
		if cfg.user == "" {
			return nil, errors.New("user is required")
		}
		if cfg.password == "" {
			return nil, errors.New("password is required")
		}
		if cfg.dbname == "" {
			return nil, errors.New("dbname is required")
		}
		dsn = fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", cfg.host, cfg.port, cfg.user, cfg.password, cfg.dbname)
		if cfg.schema != "" {
			dsn = fmt.Sprintf("%s search_path=%s", dsn, cfg.schema)
		}
	} else {
		if cfg.schema != "" {
			dsn = fmt.Sprintf("%ssearch_path=%s", dsn, cfg.schema)
		}
	}

	db, err := gorm.Open(
		postgres.Open(dsn),
		&gorm.Config{},
	)
	if err != nil {
		return nil, fmt.Errorf("failed to open postgres db, dsn: %s, error: %w", dsn, err)
	}

	if cfg.initSQL != "" {
		err = db.Exec(cfg.initSQL).Error
		if err != nil {
			return nil, fmt.Errorf("failed to execute init sql, dsn: %s, error: %w", dsn, err)
		}
	}
	return db, nil
}
