package handler

import (
	"encoding/json"
	"net/http"

	"github.com/tecnologer/tempura/pkg/utils/log"
)

func OK(writer http.ResponseWriter, data any) {
	writer.Header().Set("Content-Type", "application/json")

	encodedData, err := json.Marshal(data)
	if err != nil {
		log.Errorf("encoding response: %v", err)
		http.Error(writer, "internal server error", http.StatusInternalServerError)

		return
	}

	writer.WriteHeader(http.StatusOK)

	_, err = writer.Write(encodedData)
	if err != nil {
		log.Errorf("writing response: %v", err)
		http.Error(writer, "internal server error", http.StatusInternalServerError)
	}
}

func InternalServerError(writer http.ResponseWriter, msg string, err error) {
	Error(writer, msg, err, http.StatusInternalServerError)
}

type ErrorResponse struct {
	Error   error  `json:"error,omitempty"`
	Message string `json:"message,omitempty"`
	Code    int    `json:"code,omitempty"`
}

func Error(writer http.ResponseWriter, msg string, errResp error, statusCode int) {
	writer.Header().Set("Content-Type", "application/json")

	encodedData, err := json.Marshal(ErrorResponse{
		Error:   errResp,
		Message: msg,
		Code:    statusCode,
	})
	if err != nil {
		log.Errorf("encoding response: %v", err)
		http.Error(writer, err.Error(), http.StatusInternalServerError)

		return
	}

	writer.WriteHeader(statusCode)

	_, err = writer.Write(encodedData)
	if err != nil {
		log.Errorf("writing response: %v", err)
		http.Error(writer, err.Error(), http.StatusInternalServerError)
	}

	log.Errorf("%s: %v", msg, errResp)
}
