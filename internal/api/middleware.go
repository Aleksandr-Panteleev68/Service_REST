package api

import (
	"net/http"
	"time"
)

func (h *Handler) loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		h.logger.Info("Received request", "method", r.Method, "path", r.URL.Path)
		next.ServeHTTP(w, r)
		duration := time.Since(start)
		h.logger.Info("Request completed", "method", r.Method, "path", r.URL.Path, "duration", duration)
	})
}
