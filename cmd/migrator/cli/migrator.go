package cli

import (
	"fmt"

	"github.com/tecnologer/tempura/pkg/models"
	"github.com/tecnologer/tempura/pkg/utils/log"
	"gorm.io/gorm"
)

type Migrator struct {
	Models []any
}

func NewMigrator() *Migrator {
	return &Migrator{
		Models: []any{
			&models.Record{},
			&models.User{},
		},
	}
}

func (m *Migrator) Run(gormDB *gorm.DB) error {
	log.Info("running auto migration")

	err := gormDB.AutoMigrate(m.Models...)
	if err != nil {
		return fmt.Errorf("migrate: %w", err)
	}

	log.Infof("migration completed for %d models", len(m.Models))

	return nil
}
