package api

import (
	"encoding/json"
	"errors"
	"net/http"

	"pet-project/internal/domain"
	"pet-project/internal/service"
	"pet-project/internal/validation"
)

func (h *Handler) CreateUser(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		h.writeError(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}

	var user domain.User
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		h.writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	if err := validation.ValidateCreateUser(user); err != nil {
		h.writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	id, err := h.service.CreateUser(r.Context(), user)
	if err != nil {
		h.ServiceError(w, err)
		return
	}

	response := CreateUserResponse{ID: id}
	h.writeJSON(w, http.StatusCreated, response)
}

func (h *Handler) GetUserByID(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		h.writeError(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}

	id, err := validation.ValidateUserID(r.URL.Path)
	if err != nil {
		h.writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	user, err := h.service.GetUserByID(r.Context(), id)
	if err != nil {
		h.ServiceError(w, err)
		return
	}

	h.writeJSON(w, http.StatusOK, user)
}

func (h *Handler) ServiceError(w http.ResponseWriter, err error) {
	switch {
	case errors.Is(err, service.ErrValidation):
		h.writeError(w, http.StatusBadRequest, err.Error())
	case errors.Is(err, service.ErrConflict):
		h.writeError(w, http.StatusConflict, err.Error())
	case errors.Is(err, service.ErrNotFound):
		h.writeError(w, http.StatusNotFound, err.Error())
	default:
		h.writeError(w, http.StatusInternalServerError, "internal server error")
	}
}

func (h *Handler) writeJSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(data); err != nil {
		h.logger.Error(err, "Failed to encode JSON response")
	}
}

func (h *Handler) writeError(w http.ResponseWriter, status int, message string) {
	h.logger.Error(nil, message)
	h.writeJSON(w, status, map[string]string{"error": message})
}

type CreateUserResponse struct {
	ID int64 `json:"id"`
}
