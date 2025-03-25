package db

import (
	"errors"

	"gorm.io/gorm"
)

type DBType string

const (
	DBTypePostgres DBType = "postgres"
)

type DbOption func(*DBConfig)

type DBConfig struct {
	dbtype          DBType
	postgresOptions []PostgresOption
}

func WithDBType(dbtype DBType) DbOption {
	return func(c *DBConfig) {
		c.dbtype = dbtype
	}
}

func WithPostgresOptions(opts ...PostgresOption) DbOption {
	return func(c *DBConfig) {
		c.postgresOptions = append(c.postgresOptions, opts...)
	}
}

var defaultDBConfig = DBConfig{
	dbtype: DBTypePostgres,
}

func NewDB(opts ...DbOption) (*gorm.DB, error) {
	cfg := defaultDBConfig
	for _, opt := range opts {
		opt(&cfg)
	}

	switch cfg.dbtype {
	case DBTypePostgres:
		db, err := NewPostgresDB(cfg.postgresOptions...)
		if err != nil {
			return nil, err
		}
		return db, nil
	default:
		return nil, errors.New("invalid db type")
	}
}
