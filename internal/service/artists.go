package service

import (
	"context"
	"math"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/nik/mthen-api/internal/db"
	"github.com/nik/mthen-api/internal/models"
)

type ArtistService struct{}

func NewArtistService() *ArtistService {
	return &ArtistService{}
}

func (s *ArtistService) GetByID(ctx context.Context, id uuid.UUID) (*models.ArtistDetail, error) {
	artist, err := db.GetArtistByID(ctx, id)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	detail := &models.ArtistDetail{
		Artist: *artist,
	}

	discography, err := db.GetArtistDiscography(ctx, id, 100, 0)
	if err == nil {
		detail.Discography = discography
	}

	years, err := db.GetArtistReleaseYears(ctx, id)
	if err == nil {
		detail.ReleaseYears = years
	}

	related, err := db.GetArtistRelatedArtists(ctx, id)
	if err == nil {
		detail.RelatedArtists = related
	}

	return detail, nil
}

func (s *ArtistService) GetReleaseYears(ctx context.Context, id uuid.UUID) (*models.ActiveYearsResult, error) {
	years, err := db.GetArtistReleaseYears(ctx, id)
	if err != nil {
		return nil, err
	}

	return &models.ActiveYearsResult{
		ArtistID: id,
		Years:    years,
	}, nil
}

func (s *ArtistService) List(ctx context.Context, query string, page, limit int) (*models.PaginatedResponse, error) {
	offset := (page - 1) * limit

	artists, total, err := db.ListArtists(ctx, query, int32(limit), int32(offset))
	if err != nil {
		return nil, err
	}

	totalPages := int(math.Ceil(float64(total) / float64(limit)))
	return &models.PaginatedResponse{
		Data:       artists,
		Page:       page,
		Limit:      limit,
		TotalCount: total,
		TotalPages: totalPages,
	}, nil
}
