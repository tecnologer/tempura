package handler

import "net/http"

func (*Handler) Ping(w http.ResponseWriter, _ *http.Request) {
	_, err := w.Write([]byte("pong"))
	if err != nil {
		http.Error(w, "internal server error", http.StatusInternalServerError)
	}
}
