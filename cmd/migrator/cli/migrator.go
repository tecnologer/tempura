package cli

import (
	"errors"
	"fmt"
	"os"
	"path"
	"sort"
	"strings"

	"github.com/tecnologer/tempura/pkg/models"
	"github.com/tecnologer/tempura/pkg/utils/dir"
	"github.com/tecnologer/tempura/pkg/utils/log"
	"gorm.io/gorm"
)

const scriptsDirName = "../scripts"

type Migration struct {
	*gorm.Model
	Name     string `gorm:"unique"`
	Migrated bool
}

type Migrator struct {
	Models []any
}

func NewMigrator() *Migrator {
	return &Migrator{
		Models: []any{
			&Migration{},
			&models.Record{},
			&models.User{},
			&models.NotificationSetting{},
			&models.CasbinRule{},
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

func (m *Migrator) RunScripts(gormDB *gorm.DB) error {
	entries, err := m.scriptsEntries()
	if err != nil {
		return fmt.Errorf("scripts entries: %w", err)
	}

	scriptsCount := 0

	for _, filePath := range entries {
		isMigrated, err := m.isScriptMigrated(gormDB, filePath)
		if err != nil {
			return fmt.Errorf("check script %s: %w", filePath, err)
		}

		if isMigrated {
			continue
		}

		scriptsCount++

		err = m.runRawScript(gormDB, filePath)
		if err != nil {
			log.Warnf("script %s failed", filePath)

			if rErr := m.insertMigration(gormDB, filePath, false); rErr != nil {
				log.Warnf("insert migration %s failed", filePath)
			}

			return fmt.Errorf("run script %s: %w", filePath, err)
		}

		if err := m.insertMigration(gormDB, filePath, true); err != nil {
			log.Warnf("insert migration %s failed", filePath)
		}
	}

	if scriptsCount > 0 {
		log.Infof("%d scripts migrated", scriptsCount)
	}

	return nil
}

func (m *Migrator) scriptsEntries() ([]string, error) {
	dirPath := path.Join(dir.CallerDir(), scriptsDirName)

	entries, err := os.ReadDir(dirPath)
	if err != nil {
		return nil, fmt.Errorf("read scripts from %s: %w", scriptsDirName, err)
	}

	sort.Slice(entries, func(i, j int) bool {
		iTimeStamp := strings.Split(entries[i].Name(), "_")[0]
		jTimeStamp := strings.Split(entries[j].Name(), "_")[0]

		return iTimeStamp < jTimeStamp
	})

	files := make([]string, 0, 1)

	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}

		files = append(files, path.Join(dirPath, entry.Name()))
	}

	return files, nil
}

func (m *Migrator) runRawScript(gormDB *gorm.DB, scriptPath string) error {
	script, err := os.ReadFile(scriptPath)
	if err != nil {
		return fmt.Errorf("read script %s: %w", scriptPath, err)
	}

	tx := gormDB.Exec(string(script))
	if tx.Error != nil {
		return fmt.Errorf("exec script %s: %w", scriptPath, tx.Error)
	}

	return nil
}

func (m *Migrator) isScriptMigrated(gormDB *gorm.DB, scriptPath string) (bool, error) {
	var migration Migration

	err := gormDB.Where("name = ?", scriptPath).First(&migration).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return false, nil
		}

		return false, fmt.Errorf("find migration %s: %w", scriptPath, err)
	}

	return migration.Migrated, nil
}

func (m *Migrator) insertMigration(gormDB *gorm.DB, scriptPath string, success bool) error {
	migration := Migration{
		Name:     scriptPath,
		Migrated: success,
	}

	tx := gormDB.Create(&migration)
	if tx.Error != nil {
		return fmt.Errorf("create migration %s: %w", scriptPath, tx.Error)
	}

	return nil
}
