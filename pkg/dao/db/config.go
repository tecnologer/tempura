package db

import (
	"fmt"

	"github.com/tecnologer/tempura/pkg/contants/envvarname"
	"github.com/tecnologer/tempura/pkg/utils/envvar"
)

const (
	DefaultHost    = "localhost"
	DefaultPort    = "5432"
	DefaultUser    = "postgres"
	DefaultDB      = "tempura"
	DefaultSSLMode = "require"
)

type Config struct {
	Host     string
	Port     string
	User     string
	Password string
	DBName   string
	SSLMode  string
}

func NewConfigFromEnvVars() *Config {
	return &Config{
		Host:     envvar.ValueWithDefault(envvarname.DBHost, DefaultHost),
		Port:     envvar.ValueWithDefault(envvarname.DBPort, DefaultPort),
		User:     envvar.ValueWithDefault(envvarname.DBUsername, DefaultUser),
		Password: envvar.ValueWithDefault(envvarname.DBPassword, ""),
		DBName:   envvar.ValueWithDefault(envvarname.DBName, DefaultDB),
		SSLMode:  envvar.ValueWithDefault(envvarname.DBSSLMode, DefaultSSLMode),
	}
}

func (c *Config) DSN() (string, error) {
	if err := c.OK(); err != nil {
		return "", fmt.Errorf("invalid db config: %w", err)
	}

	return fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=%s", c.Host, c.User, c.Password, c.DBName, c.Port, c.SSLMode), nil
}

func (c *Config) OK() error {
	if c.Host == "" {
		return fmt.Errorf("missing host")
	}

	if c.Port == "" {
		c.Port = DefaultPort
	}

	if c.User == "" {
		return fmt.Errorf("missing user")
	}

	if c.DBName == "" {
		return fmt.Errorf("missing db name")
	}

	return nil
}
