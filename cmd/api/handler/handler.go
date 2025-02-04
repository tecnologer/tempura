package handler

import "github.com/tecnologer/tempura/pkg/dao/db"

type Handler struct {
	Connection *db.Connection
}

func NewHandler(connection *db.Connection) *Handler {
	return &Handler{Connection: connection}
}
