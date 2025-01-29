package db

import (
	"fmt"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type Connection struct {
	*gorm.DB
	hasTransaction bool
}

func NewConnection(config *Config) (*Connection, error) {
	dsn, err := config.DSN()
	if err != nil {
		return nil, fmt.Errorf("create DSN: %w", err)
	}

	postgresConfig := postgres.Config{
		DSN:                  dsn,
		PreferSimpleProtocol: true,
	}

	gormDB, err := gorm.Open(postgres.New(postgresConfig), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("open connection: %w", err)
	}

	return &Connection{
		DB: gormDB,
	}, nil
}

func (c *Connection) BeginTransaction() error {
	if c.hasTransaction {
		return nil
	}

	c.DB = c.DB.Begin()
	c.hasTransaction = true

	return nil
}

func (c *Connection) Commit() error {
	if !c.hasTransaction {
		return nil
	}

	c.DB = c.DB.Commit()
	c.hasTransaction = false

	return nil
}

func (c *Connection) Rollback() error {
	if !c.hasTransaction {
		return nil
	}

	c.DB = c.DB.Rollback()
	c.hasTransaction = false

	return nil
}
