package service

import (
	"context"
	"math"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/nik/mthen-api/internal/db"
	"github.com/nik/mthen-api/internal/models"
)

type AlbumService struct{}

func NewAlbumService() *AlbumService {
	return &AlbumService{}
}

func (s *AlbumService) GetByID(ctx context.Context, id uuid.UUID) (*models.AlbumDetail, error) {
	album, err := db.GetAlbumByID(ctx, id)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	songs, err := db.GetAlbumSongs(ctx, id)
	if err == nil {
		album.Songs = songs
	}

	credits, err := db.GetAlbumArtists(ctx, id)
	if err == nil {
		album.Credits = credits
	}

	return album, nil
}

func (s *AlbumService) List(ctx context.Context, year, month *int16, isOutstanding *bool, genre *string, sortBy, sortOrder string, page, limit int) (*models.PaginatedResponse, error) {
	offset := (page - 1) * limit

	albums, total, err := db.ListAlbums(ctx, year, month, isOutstanding, genre, sortBy, sortOrder, int32(limit), int32(offset))
	if err != nil {
		return nil, err
	}

	totalPages := int(math.Ceil(float64(total) / float64(limit)))
	return &models.PaginatedResponse{
		Data:       albums,
		Page:       page,
		Limit:      limit,
		TotalCount: total,
		TotalPages: totalPages,
	}, nil
}
