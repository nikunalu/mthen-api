package handler

import (
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/nik/mthen-api/internal/middleware"
	"github.com/nik/mthen-api/internal/service"
)

type AlbumHandler struct {
	svc *service.AlbumService
}

func NewAlbumHandler(svc *service.AlbumService) *AlbumHandler {
	return &AlbumHandler{svc: svc}
}

// List handles GET /api/albums
func (h *AlbumHandler) List(w http.ResponseWriter, r *http.Request) {
	page, limit := parsePagination(r)

	var year, month *int16
	if y := r.URL.Query().Get("year"); y != "" {
		v, err := strconv.ParseInt(y, 10, 16)
		if err == nil {
			yv := int16(v)
			year = &yv
		}
	}
	if m := r.URL.Query().Get("month"); m != "" {
		v, err := strconv.ParseInt(m, 10, 16)
		if err == nil {
			mv := int16(v)
			month = &mv
		}
	}

	var isOutstanding *bool
	if o := r.URL.Query().Get("outstanding"); o == "true" {
		t := true
		isOutstanding = &t
	} else if o == "false" {
		f := false
		isOutstanding = &f
	}

	genre := r.URL.Query().Get("genre")
	var genrePtr *string
	if genre != "" {
		genrePtr = &genre
	}

	sortBy := r.URL.Query().Get("sort_by")
	if sortBy == "" {
		sortBy = "release_date"
	}
	sortOrder := r.URL.Query().Get("sort_order")
	if sortOrder == "" {
		sortOrder = "desc"
	}

	resp, err := h.svc.List(r.Context(), year, month, isOutstanding, genrePtr, sortBy, sortOrder, page, limit)
	if err != nil {
		middleware.WriteError(w, http.StatusInternalServerError, "failed to list albums")
		return
	}
	middleware.WriteJSON(w, http.StatusOK, resp)
}

// GetByID handles GET /api/albums/{id}
func (h *AlbumHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		middleware.WriteError(w, http.StatusBadRequest, "invalid album ID")
		return
	}

	album, err := h.svc.GetByID(r.Context(), id)
	if err != nil {
		middleware.WriteError(w, http.StatusInternalServerError, "failed to fetch album")
		return
	}
	if album == nil {
		middleware.WriteError(w, http.StatusNotFound, "album not found")
		return
	}

	middleware.WriteJSON(w, http.StatusOK, album)
}
