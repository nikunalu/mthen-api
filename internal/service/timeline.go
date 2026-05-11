package service

import (
	"context"

	"github.com/nik/mthen-api/internal/db"
	"github.com/nik/mthen-api/internal/models"
)

type TimelineService struct{}

func NewTimelineService() *TimelineService {
	return &TimelineService{}
}

func (s *TimelineService) GetYears(ctx context.Context) ([]models.TimelineYear, error) {
	return db.GetTimelineYears(ctx)
}

func (s *TimelineService) GetMonths(ctx context.Context) ([]models.TimelineMonth, error) {
	return db.GetTimelineMonths(ctx)
}

func (s *TimelineService) GetMonthDetail(ctx context.Context, year, month int16, page, limit int) (*models.TimelineMonthDetail, error) {
	offset := (page - 1) * limit
	albums, total, err := db.GetTimelineMonthDetail(ctx, year, month, int32(limit), int32(offset))
	if err != nil {
		return nil, err
	}

	return &models.TimelineMonthDetail{
		Year:       year,
		Month:      month,
		TotalCount: total,
		Albums:     albums,
	}, nil
}

func (s *TimelineService) GetAlbumsByYear(ctx context.Context, year int16) ([]models.AlbumWithArtist, error) {
	return db.GetAlbumsByYear(ctx, year)
}
