package dao

import (
	"context"
	"fmt"
	"os"

	"github.com/tecnologer/tempura/pkg/contants/envvarname"
	"github.com/tecnologer/tempura/pkg/dao/db"
	"github.com/tecnologer/tempura/pkg/models"
	"github.com/tecnologer/tempura/pkg/utils/log"
)

type Records struct {
	cnn           *db.Connection
	notifications *Notification
}

func NewRecords(cnn *db.Connection) *Records {
	notifications, err := NewNotification(cnn, os.Getenv(envvarname.TelegramBotToken))
	if err != nil {
		log.Errorf("creating notification service: %v", err)
	}

	return &Records{
		cnn:           cnn,
		notifications: notifications,
	}
}

func (r *Records) InsertRecord(_ context.Context, record *models.Record) (*models.Record, error) {
	log.Debugf("inserting record: %s", record)

	tx := r.cnn.DB.Create(record)
	if tx.Error != nil {
		return record, fmt.Errorf("inserting record: %w", tx.Error)
	}

	log.Debug("record inserted")

	err := r.NotifyNewRecord(record)
	if err != nil {
		return record, fmt.Errorf("notifying new record: %w", err)
	}

	return record, nil
}

type Filter struct {
	Limit int
}

func (r *Records) GetRecords(_ context.Context, filters Filter) ([]models.Record, error) {
	var records []models.Record

	tx := r.cnn.DB.Order("created_at desc").Limit(filters.Limit).Find(&records)
	if tx.Error != nil {
		return records, fmt.Errorf("getting records: %w", tx.Error)
	}

	return records, nil
}

func (r *Records) GetRecord(_ context.Context, id string) (*models.Record, error) {
	var record models.Record

	tx := r.cnn.DB.First(&record, id)
	if tx.Error != nil {
		return &record, fmt.Errorf("getting record: %w", tx.Error)
	}

	return &record, nil
}

func (r *Records) NotifyNewRecord(record *models.Record) error {
	if r.notifications == nil {
		return nil
	}

	err := r.notifications.NotifyNewRecord(record)
	if err != nil {
		return fmt.Errorf("notifying new record: %w", err)
	}

	return nil
}
