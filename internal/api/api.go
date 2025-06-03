package api

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"pet-project/internal/config"
	"pet-project/internal/domain"
	"pet-project/internal/logger"
	"pet-project/internal/service"
)

type Handler struct {
	service service.UserService
	logger  *logger.Logger
	config  *config.Config
	mux     *http.ServeMux
}

func NewHandler(service service.UserService, logger *logger.Logger, config *config.Config) *Handler {
	h := &Handler{
		service: service,
		logger:  logger,
		config:  config,
		mux:     http.NewServeMux(),
	}
	h.setupRoutes()
	return h
}

func (h *Handler) setupRoutes() {
	h.mux.HandleFunc("/users", h.handleUsers)
	h.mux.HandleFunc("/users", h.handleUserByID)
}

func (h *Handler) loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		h.logger.Info("Received request", "method", r.Method, "path", r.URL.Path)

		next.ServeHTTP(w, r)

		duration := time.Since(start)
		h.logger.Info("Request completed", "method", r.Method, "path", r.URL.Path, "duration", duration)
	})
}

func (h *Handler) handleUsers(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		h.writeError(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}

	var user domain.User
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		h.writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	if user.FirstName == "" || user.LastName == "" {
		h.writeError(w, http.StatusBadRequest, "first_name and last_name are required")
		return
	}

	user.FullName = user.FirstName + " " + user.LastName

	id, err := h.service.CreateUser(r.Context(), user)
	if err != nil {
		h.handleServiceError(w, err)
		return
	}

	response := map[string]int64{"id": id}
	h.writeJSON(w, http.StatusCreated, response)
}

func (h *Handler) handleUserByID(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		h.writeError(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}

	parts := strings.Split(r.URL.Path, "/")
	if len(parts) != 3 || parts[1] != "users" {
		h.writeError(w, http.StatusBadRequest, "invalid URL path")
		return
	}

	idStr := parts[2]
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		h.writeError(w, http.StatusBadRequest, "invalid user ID")
		return
	}

	user, err := h.service.GetUserByID(r.Context(), id)
	if err != nil {
		h.handleServiceError(w, err)
		return
	}

	user.FullName = user.FirstName + " " + user.LastName
	h.writeJSON(w, http.StatusOK, user)
}

func (h *Handler) handleServiceError(w http.ResponseWriter, err error) {
	switch {
	case strings.Contains(err.Error(), "first_name and last_name cannot be empty"),
		strings.Contains(err.Error(), "age must be at least 18"),
		strings.Contains(err.Error(), "password must be at least 8 characters"):
		h.writeError(w, http.StatusBadRequest, err.Error())
	case strings.Contains(err.Error(), "user with name"):
		h.writeError(w, http.StatusConflict, err.Error())
	case strings.Contains(err.Error(), "user not found"):
		h.writeError(w, http.StatusNotFound, err.Error())
	default:
		h.writeError(w, http.StatusInternalServerError, "internal server error")
	}
}

func (h *Handler) writeJSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Contenr-Type", "application/json")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(data); err != nil {
		h.logger.Error(err, "Failed to encode JSON response")
	}
}

func (h *Handler) writeError(w http.ResponseWriter, status int, message string) {
	h.logger.Error(nil, message)
	h.writeJSON(w, status, map[string]string{"error": message})
}

func (h *Handler) StartServer(ctx context.Context) error {
	server := &http.Server{
		Addr: fmt.Sprintf("%s:%d", h.config.HTTPServer.Address, h.config.HTTPServer.Port),
		Handler: h.loggingMiddleware(h.mux),
		ReadTimeout: h.config.HTTPServer.Timeout,
		WriteTimeout: h.config.HTTPServer.Timeout,
		IdleTimeout: h.config.HTTPServer.IdleTimeout,
	}

	go func() {
		h.logger.Info("Starting HTTP Server", "address", server.Addr)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			h.logger.Error(err, "Failed to start HTTP server")
		}
	}()

	<-ctx.Done()
	h.logger.Info("Shutting down HTTP server")

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := server.Shutdown(shutdownCtx); err != nil {
		h.logger.Error(err, "Failed to shutdown HTTP server gracefully")
		return fmt.Errorf("failed to shotdown server: %w", err)
	}

	h.logger.Info("HTTP server stopped")
	return nil
}
