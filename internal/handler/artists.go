package handler

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/nik/mthen-api/internal/middleware"
	"github.com/nik/mthen-api/internal/service"
)

type ArtistHandler struct {
	svc *service.ArtistService
}

func NewArtistHandler(svc *service.ArtistService) *ArtistHandler {
	return &ArtistHandler{svc: svc}
}

// List handles GET /api/artists
func (h *ArtistHandler) List(w http.ResponseWriter, r *http.Request) {
	page, limit := parsePagination(r)
	query := r.URL.Query().Get("q")

	resp, err := h.svc.List(r.Context(), query, page, limit)
	if err != nil {
		middleware.WriteError(w, http.StatusInternalServerError, "failed to list artists")
		return
	}
	middleware.WriteJSON(w, http.StatusOK, resp)
}

// GetByID handles GET /api/artists/{id}
func (h *ArtistHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		middleware.WriteError(w, http.StatusBadRequest, "invalid artist ID")
		return
	}

	artist, err := h.svc.GetByID(r.Context(), id)
	if err != nil {
		middleware.WriteError(w, http.StatusInternalServerError, "failed to fetch artist")
		return
	}
	if artist == nil {
		middleware.WriteError(w, http.StatusNotFound, "artist not found")
		return
	}

	middleware.WriteJSON(w, http.StatusOK, artist)
}

// GetReleaseYears handles GET /api/artists/{id}/years
func (h *ArtistHandler) GetReleaseYears(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		middleware.WriteError(w, http.StatusBadRequest, "invalid artist ID")
		return
	}

	years, err := h.svc.GetReleaseYears(r.Context(), id)
	if err != nil {
		middleware.WriteError(w, http.StatusInternalServerError, "failed to fetch release years")
		return
	}

	middleware.WriteJSON(w, http.StatusOK, years)
}
