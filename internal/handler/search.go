package handler

import (
	"net/http"

	"github.com/nik/mthen-api/internal/middleware"
	"github.com/nik/mthen-api/internal/service"
)

type SearchHandler struct {
	svc *service.SearchService
}

func NewSearchHandler(svc *service.SearchService) *SearchHandler {
	return &SearchHandler{svc: svc}
}

// Search handles GET /api/search?q=
func (h *SearchHandler) Search(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query().Get("q")
	if query == "" {
		middleware.WriteError(w, http.StatusBadRequest, "search query 'q' is required")
		return
	}

	page, limit := parsePagination(r)

	resp, err := h.svc.Search(r.Context(), query, page, limit)
	if err != nil {
		middleware.WriteError(w, http.StatusInternalServerError, "search failed")
		return
	}

	middleware.WriteJSON(w, http.StatusOK, resp)
}
