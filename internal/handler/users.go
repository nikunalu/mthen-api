package handler

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/nik/mthen-api/internal/middleware"
	"github.com/nik/mthen-api/internal/models"
	"github.com/nik/mthen-api/internal/service"
)

type UserHandler struct {
	svc *service.UserService
}

func NewUserHandler(svc *service.UserService) *UserHandler {
	return &UserHandler{svc: svc}
}

// GetProfile handles GET /api/me/profile
func (h *UserHandler) GetProfile(w http.ResponseWriter, r *http.Request) {
	userID := middleware.GetUserID(r)

	profile, err := h.svc.GetProfile(r.Context(), userID)
	if err != nil {
		middleware.WriteError(w, http.StatusInternalServerError, "failed to fetch profile")
		return
	}

	middleware.WriteJSON(w, http.StatusOK, profile)
}

// UpdateProfile handles PUT /api/me/profile
func (h *UserHandler) UpdateProfile(w http.ResponseWriter, r *http.Request) {
	userID := middleware.GetUserID(r)

	var req models.UpdateProfileRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		middleware.WriteError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	profile, err := h.svc.UpdateProfile(r.Context(), userID, req)
	if err != nil {
		middleware.WriteError(w, http.StatusInternalServerError, "failed to update profile")
		return
	}

	middleware.WriteJSON(w, http.StatusOK, profile)
}

// GetTopAlbums handles GET /api/me/top-albums
func (h *UserHandler) GetTopAlbums(w http.ResponseWriter, r *http.Request) {
	userID := middleware.GetUserID(r)

	resp, err := h.svc.GetTopAlbums(r.Context(), userID)
	if err != nil {
		middleware.WriteError(w, http.StatusInternalServerError, "failed to fetch top albums")
		return
	}

	middleware.WriteJSON(w, http.StatusOK, resp)
}

// UpdateTopAlbums handles PUT /api/me/top-albums
func (h *UserHandler) UpdateTopAlbums(w http.ResponseWriter, r *http.Request) {
	userID := middleware.GetUserID(r)

	var req models.UpdateTopItemsRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		middleware.WriteError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	resp, err := h.svc.UpdateTopAlbums(r.Context(), userID, req)
	if err != nil {
		if errResp, ok := err.(*models.ErrorResponse); ok {
			middleware.WriteError(w, http.StatusBadRequest, errResp.Message)
			return
		}
		middleware.WriteError(w, http.StatusInternalServerError, "failed to update top albums")
		return
	}

	middleware.WriteJSON(w, http.StatusOK, resp)
}

// GetTopSongs handles GET /api/me/top-songs
func (h *UserHandler) GetTopSongs(w http.ResponseWriter, r *http.Request) {
	userID := middleware.GetUserID(r)

	resp, err := h.svc.GetTopSongs(r.Context(), userID)
	if err != nil {
		middleware.WriteError(w, http.StatusInternalServerError, "failed to fetch top songs")
		return
	}

	middleware.WriteJSON(w, http.StatusOK, resp)
}

// UpdateTopSongs handles PUT /api/me/top-songs
func (h *UserHandler) UpdateTopSongs(w http.ResponseWriter, r *http.Request) {
	userID := middleware.GetUserID(r)

	var req models.UpdateTopItemsRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		middleware.WriteError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	resp, err := h.svc.UpdateTopSongs(r.Context(), userID, req)
	if err != nil {
		if errResp, ok := err.(*models.ErrorResponse); ok {
			middleware.WriteError(w, http.StatusBadRequest, errResp.Message)
			return
		}
		middleware.WriteError(w, http.StatusInternalServerError, "failed to update top songs")
		return
	}

	middleware.WriteJSON(w, http.StatusOK, resp)
}

// GetTopArtists handles GET /api/me/top-artists
func (h *UserHandler) GetTopArtists(w http.ResponseWriter, r *http.Request) {
	userID := middleware.GetUserID(r)

	resp, err := h.svc.GetTopArtists(r.Context(), userID)
	if err != nil {
		middleware.WriteError(w, http.StatusInternalServerError, "failed to fetch top artists")
		return
	}

	middleware.WriteJSON(w, http.StatusOK, resp)
}

// UpdateTopArtists handles PUT /api/me/top-artists
func (h *UserHandler) UpdateTopArtists(w http.ResponseWriter, r *http.Request) {
	userID := middleware.GetUserID(r)

	var req models.UpdateTopItemsRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		middleware.WriteError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	resp, err := h.svc.UpdateTopArtists(r.Context(), userID, req)
	if err != nil {
		if errResp, ok := err.(*models.ErrorResponse); ok {
			middleware.WriteError(w, http.StatusBadRequest, errResp.Message)
			return
		}
		middleware.WriteError(w, http.StatusInternalServerError, "failed to update top artists")
		return
	}

	middleware.WriteJSON(w, http.StatusOK, resp)
}

// CreateListening handles POST /api/me/listening
func (h *UserHandler) CreateListening(w http.ResponseWriter, r *http.Request) {
	userID := middleware.GetUserID(r)

	var req models.CreateListeningRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		middleware.WriteError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	record, err := h.svc.CreateListening(r.Context(), userID, req)
	if err != nil {
		middleware.WriteError(w, http.StatusInternalServerError, "failed to create listening record")
		return
	}

	middleware.WriteJSON(w, http.StatusCreated, record)
}

// ListListening handles GET /api/me/listening
func (h *UserHandler) ListListening(w http.ResponseWriter, r *http.Request) {
	userID := middleware.GetUserID(r)
	page, limit := parsePagination(r)

	resp, err := h.svc.ListListening(r.Context(), userID, page, limit)
	if err != nil {
		middleware.WriteError(w, http.StatusInternalServerError, "failed to fetch listening records")
		return
	}

	middleware.WriteJSON(w, http.StatusOK, resp)
}

// GetMonthlySet handles GET /api/me/monthly-set/{year}/{month}
func (h *UserHandler) GetMonthlySet(w http.ResponseWriter, r *http.Request) {
	userID := middleware.GetUserID(r)

	yearStr := chi.URLParam(r, "year")
	monthStr := chi.URLParam(r, "month")

	year, err := strconv.ParseInt(yearStr, 10, 16)
	if err != nil {
		middleware.WriteError(w, http.StatusBadRequest, "invalid year")
		return
	}
	month, err := strconv.ParseInt(monthStr, 10, 16)
	if err != nil || month < 1 || month > 12 {
		middleware.WriteError(w, http.StatusBadRequest, "invalid month")
		return
	}

	set, err := h.svc.GetMonthlySet(r.Context(), userID, int16(year), int16(month))
	if err != nil {
		middleware.WriteError(w, http.StatusNotFound, "monthly set not found")
		return
	}

	middleware.WriteJSON(w, http.StatusOK, set)
}

// UpsertMonthlySet handles PUT /api/me/monthly-set/{year}/{month}
func (h *UserHandler) UpsertMonthlySet(w http.ResponseWriter, r *http.Request) {
	userID := middleware.GetUserID(r)

	yearStr := chi.URLParam(r, "year")
	monthStr := chi.URLParam(r, "month")

	year, err := strconv.ParseInt(yearStr, 10, 16)
	if err != nil {
		middleware.WriteError(w, http.StatusBadRequest, "invalid year")
		return
	}
	month, err := strconv.ParseInt(monthStr, 10, 16)
	if err != nil || month < 1 || month > 12 {
		middleware.WriteError(w, http.StatusBadRequest, "invalid month")
		return
	}

	var req models.UpdateMonthlySetRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		middleware.WriteError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	set, err := h.svc.UpsertMonthlySet(r.Context(), userID, int16(year), int16(month), req)
	if err != nil {
		middleware.WriteError(w, http.StatusInternalServerError, "failed to save monthly set")
		return
	}

	middleware.WriteJSON(w, http.StatusOK, set)
}
