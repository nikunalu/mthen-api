package handler

import (
	"net/http"

	"github.com/nik/mthen-api/internal/middleware"
	"github.com/nik/mthen-api/internal/models"
	"github.com/nik/mthen-api/internal/service"
)

type GenreHandler struct {
	svc *service.GenreService
}

func NewGenreHandler(svc *service.GenreService) *GenreHandler {
	return &GenreHandler{svc: svc}
}

// List handles GET /api/genres
func (h *GenreHandler) List(w http.ResponseWriter, r *http.Request) {
	genres, err := h.svc.List(r.Context())
	if err != nil {
		middleware.WriteError(w, http.StatusInternalServerError, "failed to fetch genres")
		return
	}

	if genres == nil {
		genres = []models.GenreCount{}
	}

	middleware.WriteJSON(w, http.StatusOK, genres)
}
