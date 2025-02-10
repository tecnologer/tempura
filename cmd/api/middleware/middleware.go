package middleware

import (
	"net/http"

	"github.com/tecnologer/tempura/pkg/utils/log"
)

//nolint:gofumpt // TBD
type Middleware struct {
}

func NewMiddleware() *Middleware {
	return &Middleware{}
}

func (m *Middleware) ServeHTTP(rw http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	log.Info("middleware")

	next(rw, r)
}
