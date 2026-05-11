package service

import (
	"context"
	"math"

	"github.com/nik/mthen-api/internal/db"
	"github.com/nik/mthen-api/internal/models"
)

type SearchService struct{}

func NewSearchService() *SearchService {
	return &SearchService{}
}

func (s *SearchService) Search(ctx context.Context, query string, page, limit int) (*models.PaginatedResponse, error) {
	offset := (page - 1) * limit

	results, err := db.UnifiedSearch(ctx, query, int32(limit), int32(offset))
	if err != nil {
		return nil, err
	}

	total := int64(len(results))
	if len(results) == int(limit) {
		// Rough estimate for pagination — could be improved with a count query
		total = int64(offset+int(limit)) + 1
	}

	totalPages := int(math.Ceil(float64(total) / float64(limit)))
	return &models.PaginatedResponse{
		Data:       results,
		Page:       page,
		Limit:      limit,
		TotalCount: total,
		TotalPages: totalPages,
	}, nil
}
