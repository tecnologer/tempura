package handler

import (
	"encoding/json"
	"net/http"

	"github.com/tecnologer/tempura/pkg/dao"
	"github.com/tecnologer/tempura/pkg/models"
	"github.com/tecnologer/tempura/pkg/utils/log"
)

func (h *Handler) CreateRecord(writer http.ResponseWriter, request *http.Request) {
	var record *models.Record

	log.Debug("api: creating record")

	err := json.NewDecoder(request.Body).Decode(&record)
	if err != nil {
		log.Errorf("decoding request: %v", err)

		http.Error(writer, "invalid request", http.StatusBadRequest)

		return
	}

	log.Debugf("record: %v", *record)

	record, err = h.records.InsertRecord(request.Context(), record)
	if err != nil {
		InternalServerError(writer, "inserting record", err)

		return
	}

	OK(writer, record)
}

func (h *Handler) GetRecords(writer http.ResponseWriter, request *http.Request) {
	filters := dao.Filter{
		Limit: 1000,
	}

	records, err := h.records.GetRecords(request.Context(), filters)
	if err != nil {
		InternalServerError(writer, "getting records", err)

		return
	}

	OK(writer, records)
}

func (h *Handler) GetRecord(writer http.ResponseWriter, request *http.Request) {
	id := request.URL.Query().Get("id")

	record, err := h.records.GetRecord(request.Context(), id)
	if err != nil {
		InternalServerError(writer, "getting record", err)

		return
	}

	OK(writer, record)
}
