package api

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"pet-project/internal/config"
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

func (h *Handler) StartServer(ctx context.Context) error {
	server := &http.Server{
		Addr:         fmt.Sprintf("%s:%d", h.config.HTTPServer.Address, h.config.HTTPServer.Port),
		Handler:      h.loggingMiddleware(h.mux),
		ReadTimeout:  h.config.HTTPServer.Timeout,
		WriteTimeout: h.config.HTTPServer.Timeout,
		IdleTimeout:  h.config.HTTPServer.IdleTimeout,
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
		return fmt.Errorf("failed to shutdown server: %w", err)
	}

	h.logger.Info("HTTP server stopped")
	return nil
}

func (h *Handler) setupRoutes() {
	h.mux.HandleFunc("/users", h.CreateUser)
	h.mux.HandleFunc("/users/", h.GetUserByID)
}
