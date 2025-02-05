package db

import (
	"fmt"

	"github.com/tecnologer/tempura/pkg/utils/log"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type Connection struct {
	*gorm.DB
	hasTransaction bool
}

func NewConnection(config *Config) (*Connection, error) {
	log.Infof("connecting to the DB at %s:%s as user %s", config.Host, config.Port, config.User)

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

	log.Infof("DB connection established")

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

	log.Debug("transaction started")

	return nil
}

func (c *Connection) Commit() error {
	if !c.hasTransaction {
		return nil
	}

	c.DB = c.DB.Commit()
	c.hasTransaction = false

	log.Debug("transaction committed")

	return nil
}

func (c *Connection) Rollback() error {
	if !c.hasTransaction {
		return nil
	}

	c.DB = c.DB.Rollback()
	c.hasTransaction = false

	log.Debug("transaction rolled back")

	return nil
}
