package router

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/tecnologer/tempura/cmd/api/handler"
)

const (
	APIVersionV1 = "/v1"
	PingPath     = "/ping"
	RecordsPath  = "/records"
	RecordPath   = "/records/{id}"
)

func New(apiHandler *handler.Handler) http.Handler {
	router := mux.NewRouter()

	router.HandleFunc(APIVersionV1+PingPath, apiHandler.Ping).Methods(http.MethodGet)
	router.HandleFunc(APIVersionV1+RecordsPath, apiHandler.CreateRecord).Methods(http.MethodPost)
	router.HandleFunc(APIVersionV1+RecordsPath, apiHandler.GetRecords).Methods(http.MethodGet)
	router.HandleFunc(APIVersionV1+RecordPath, apiHandler.GetRecord).Methods(http.MethodGet)

	return router
}
