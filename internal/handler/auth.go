package handler

import (
	"encoding/json"
	"net/http"

	"github.com/nik/mthen-api/internal/middleware"
	"github.com/nik/mthen-api/internal/models"
	"github.com/nik/mthen-api/internal/service"
)

type AuthHandler struct {
	svc *service.AuthService
}

func NewAuthHandler(svc *service.AuthService) *AuthHandler {
	return &AuthHandler{svc: svc}
}

// Register handles POST /api/auth/register
func (h *AuthHandler) Register(w http.ResponseWriter, r *http.Request) {
	var req models.RegisterRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		middleware.WriteError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	resp, err := h.svc.Register(r.Context(), req)
	if err != nil {
		middleware.WriteError(w, http.StatusConflict, err.Error())
		return
	}

	middleware.WriteJSON(w, http.StatusCreated, resp)
}

// Login handles POST /api/auth/login
func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	var req models.LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		middleware.WriteError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	resp, err := h.svc.Login(r.Context(), req)
	if err != nil {
		middleware.WriteError(w, http.StatusUnauthorized, err.Error())
		return
	}

	middleware.WriteJSON(w, http.StatusOK, resp)
}
