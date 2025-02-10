package handler

import (
	"github.com/tecnologer/tempura/pkg/dao"
	"github.com/tecnologer/tempura/pkg/dao/db"
)

type Handler struct {
	Connection *db.Connection
	records    *dao.Records
}

func NewHandler(connection *db.Connection) *Handler {
	return &Handler{
		Connection: connection,
		records:    dao.NewRecords(connection),
	}
}
