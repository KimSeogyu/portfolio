package config

import (
	"errors"
	"os"

	"gopkg.in/yaml.v3"
)

type Config struct {
	DBConfig   DBConfig
	GRPCConfig GRPCConfig
}

type DB struct {
	Host     string
	Port     int
	User     string
	Password string
	DBName   string
	Schema   string
}

type GRPCConfig struct {
	GrpcPort    int
	GatewayPort int
}

type DBType string

const (
	DBTypePostgres     DBType = "postgres"
	DBTypeTestPostgres DBType = "test_postgres"
)

type DBConfig struct {
	DBType DBType
	DB     DB
}

func NewConfigFromFile(filename string) (*Config, error) {
	cfg := &Config{}

	// read config file
	content, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	// parse config file
	err = yaml.Unmarshal(content, cfg)
	if err != nil {
		return nil, err
	}

	return cfg, nil
}

func (c *Config) Validate() error {
	if c.DBConfig.DBType == "" {
		return errors.New("db type is required (postgres or sqlite)")
	}

	if c.DBConfig.DBType == DBTypePostgres {
		if c.DBConfig.DB.Host == "" {
			return errors.New("host is required")
		}
		if c.DBConfig.DB.Port == 0 {
			return errors.New("port is required")
		}
		if c.DBConfig.DB.User == "" {
			return errors.New("user is required")
		}
		if c.DBConfig.DB.Password == "" {
			return errors.New("password is required")
		}
		if c.DBConfig.DB.DBName == "" {
			return errors.New("db name is required")
		}
		if c.DBConfig.DB.Schema == "" {
			return errors.New("schema is required")
		}
	}

	if c.DBConfig.DBType == DBTypeTestPostgres {
		if c.DBConfig.DB.Host == "" {
			return errors.New("host is required")
		}
		if c.DBConfig.DB.Port == 0 {
			return errors.New("port is required")
		}
		if c.DBConfig.DB.User == "" {
			return errors.New("user is required")
		}
		if c.DBConfig.DB.Password == "" {
			return errors.New("password is required")
		}
		if c.DBConfig.DB.DBName == "" {
			return errors.New("db name is required")
		}
	}

	return nil
}

func Local() *Config {
	return &Config{
		DBConfig: DBConfig{
			DBType: DBTypePostgres,
			DB: DB{
				Host:     "docker.for.mac.localhost",
				Port:     5432,
				User:     "postgres",
				Password: "postgres",
				DBName:   "postgres",
				Schema:   "public",
			},
		},
		GRPCConfig: GRPCConfig{
			GrpcPort:    10010,
			GatewayPort: 8080,
		},
	}
}

func Test() *Config {
	return &Config{
		DBConfig: DBConfig{
			DBType: DBTypeTestPostgres,
			DB: DB{
				Host:     "localhost",
				Port:     5432,
				User:     "postgres",
				Password: "postgres",
				DBName:   "postgres",
				Schema:   "public",
			},
		},
		GRPCConfig: GRPCConfig{
			GrpcPort:    10010,
			GatewayPort: 8080,
		},
	}
}

func (c *Config) MarshalYAML() ([]byte, error) {
	return yaml.Marshal(c)
}
