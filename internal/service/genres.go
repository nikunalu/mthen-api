package service

import (
	"context"

	"github.com/nik/mthen-api/internal/db"
	"github.com/nik/mthen-api/internal/models"
)

type GenreService struct{}

func NewGenreService() *GenreService {
	return &GenreService{}
}

func (s *GenreService) List(ctx context.Context) ([]models.GenreCount, error) {
	return db.ListGenres(ctx)
}
