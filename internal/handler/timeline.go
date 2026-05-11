package handler

import (
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/nik/mthen-api/internal/middleware"
	"github.com/nik/mthen-api/internal/models"
	"github.com/nik/mthen-api/internal/service"
)

type TimelineHandler struct {
	svc *service.TimelineService
}

func NewTimelineHandler(svc *service.TimelineService) *TimelineHandler {
	return &TimelineHandler{svc: svc}
}

// GetYears handles GET /api/timeline
func (h *TimelineHandler) GetYears(w http.ResponseWriter, r *http.Request) {
	years, err := h.svc.GetYears(r.Context())
	if err != nil {
		middleware.WriteError(w, http.StatusInternalServerError, "failed to fetch timeline years")
		return
	}
	if years == nil {
		years = []models.TimelineYear{}
	}
	middleware.WriteJSON(w, http.StatusOK, map[string]any{"years": years})
}

// GetYear handles GET /api/timeline/{year}
func (h *TimelineHandler) GetYear(w http.ResponseWriter, r *http.Request) {
	yearStr := chi.URLParam(r, "year")
	year, err := strconv.ParseInt(yearStr, 10, 16)
	if err != nil {
		middleware.WriteError(w, http.StatusBadRequest, "invalid year")
		return
	}

	albums, err := h.svc.GetAlbumsByYear(r.Context(), int16(year))
	if err != nil {
		middleware.WriteError(w, http.StatusInternalServerError, "failed to fetch year timeline")
		return
	}

	// Group by month for the frontend's YearDetail format
	type albumSummary struct {
		ID            string   `json:"id"`
		Title         string   `json:"title"`
		ArtistName    string   `json:"artist_name"`
		ArtistID      string   `json:"artist_id"`
		ReleaseDate   string   `json:"release_date"`
		ReleaseMonth  int16    `json:"release_month"`
		IsOutstanding bool     `json:"is_outstanding"`
		Genres        []string `json:"genres"`
		CoverArtURL   string   `json:"cover_art_url"`
	}

	type monthGroup struct {
		Month      int16          `json:"month"`
		AlbumCount int64          `json:"album_count"`
		Albums     []albumSummary `json:"albums"`
	}

	months := make([]monthGroup, 12)
	monthMap := make(map[int16][]albumSummary)
	outstandingCount := 0

	for _, a := range albums {
		m := int16(0)
		if a.ReleaseMonth != nil {
			m = *a.ReleaseMonth
		}
		rd := ""
		if a.ReleaseDate != nil {
			rd = a.ReleaseDate.Format("2006-01-02")
		}
		cu := ""
		if a.CoverArtURL != nil {
			cu = *a.CoverArtURL
		}
		aid := a.ArtistID.String()
		monthMap[m] = append(monthMap[m], albumSummary{
			ID: a.ID.String(), Title: a.Title, ArtistName: a.ArtistName,
			ArtistID: aid, ReleaseDate: rd, ReleaseMonth: m,
			IsOutstanding: a.IsOutstanding, Genres: a.Genres, CoverArtURL: cu,
		})
		if a.IsOutstanding {
			outstandingCount++
		}
	}

	for m := int16(1); m <= 12; m++ {
		albs := monthMap[m]
		if albs == nil {
			albs = []albumSummary{}
		}
		months[m-1] = monthGroup{Month: m, AlbumCount: int64(len(albs)), Albums: albs}
	}

	middleware.WriteJSON(w, http.StatusOK, map[string]any{
		"year":              year,
		"album_count":       len(albums),
		"outstanding_count": outstandingCount,
		"months":            months,
	})
}

// GetMonth handles GET /api/timeline/{year}/{month}
func (h *TimelineHandler) GetMonth(w http.ResponseWriter, r *http.Request) {
	yearStr := chi.URLParam(r, "year")
	monthStr := chi.URLParam(r, "month")

	year, err := strconv.ParseInt(yearStr, 10, 16)
	if err != nil {
		middleware.WriteError(w, http.StatusBadRequest, "invalid year")
		return
	}
	month, err := strconv.ParseInt(monthStr, 10, 16)
	if err != nil || month < 1 || month > 12 {
		middleware.WriteError(w, http.StatusBadRequest, "invalid month (must be 1-12)")
		return
	}

	page, limit := parsePagination(r)

	resp, err := h.svc.GetMonthDetail(r.Context(), int16(year), int16(month), page, limit)
	if err != nil {
		middleware.WriteError(w, http.StatusInternalServerError, "failed to fetch month timeline")
		return
	}
	middleware.WriteJSON(w, http.StatusOK, resp)
}
