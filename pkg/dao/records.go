package dao

import (
	"context"
	"fmt"

	"github.com/tecnologer/tempura/pkg/dao/db"
	"github.com/tecnologer/tempura/pkg/models"
	"github.com/tecnologer/tempura/pkg/utils/log"
)

type Records struct {
	cnn *db.Connection
}

func NewRecords(cnn *db.Connection) *Records {
	return &Records{cnn: cnn}
}

func (r *Records) InsertRecord(_ context.Context, record *models.Record) (*models.Record, error) {
	log.Debug("inserting record")

	tx := r.cnn.DB.Create(record)
	if tx.Error != nil {
		return record, fmt.Errorf("inserting record: %w", tx.Error)
	}

	log.Debug("record inserted")

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
