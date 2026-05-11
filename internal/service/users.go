package service

import (
	"context"
	"math"
	"time"

	"github.com/google/uuid"
	"github.com/nik/mthen-api/internal/db"
	"github.com/nik/mthen-api/internal/models"
)

type UserService struct{}

func NewUserService() *UserService {
	return &UserService{}
}

func (s *UserService) GetProfile(ctx context.Context, userID uuid.UUID) (*models.UserProfileResponse, error) {
	user, err := db.GetUserByID(ctx, userID)
	if err != nil {
		return nil, err
	}

	profile := models.UserProfileResponse{
		ID:          user.ID,
		Email:       user.Email,
		DisplayName: user.DisplayName,
		AvatarURL:   user.AvatarURL,
		Bio:         user.Bio,
		JoinedAt:    user.JoinedAt,
	}

	// Load top items
	if len(user.Top10Albums) > 0 {
		albums, err := db.GetTopAlbumsByIDs(ctx, user.Top10Albums)
		if err == nil {
			profile.Top10Albums = albums
		}
	}

	if len(user.Top10Songs) > 0 {
		songs, err := db.GetTopSongsByIDs(ctx, user.Top10Songs)
		if err == nil {
			profile.Top10Songs = songs
		}
	}

	if len(user.Top10Artists) > 0 {
		artists, err := db.GetTopArtistsByIDs(ctx, user.Top10Artists)
		if err == nil {
			profile.Top10Artists = artists
		}
	}

	return &profile, nil
}

func (s *UserService) UpdateProfile(ctx context.Context, userID uuid.UUID, req models.UpdateProfileRequest) (*models.UserProfileResponse, error) {
	user, err := db.UpdateUserProfile(ctx, userID, req.DisplayName, req.AvatarURL, req.Bio)
	if err != nil {
		return nil, err
	}

	profile := models.UserProfileResponse{
		ID:          user.ID,
		Email:       user.Email,
		DisplayName: user.DisplayName,
		AvatarURL:   user.AvatarURL,
		Bio:         user.Bio,
		JoinedAt:    user.JoinedAt,
	}
	return &profile, nil
}

func (s *UserService) GetTopAlbums(ctx context.Context, userID uuid.UUID) (*models.PaginatedResponse, error) {
	user, err := db.GetUserByID(ctx, userID)
	if err != nil {
		return nil, err
	}

	albums, err := db.GetTopAlbumsByIDs(ctx, user.Top10Albums)
	if err != nil {
		return nil, err
	}

	return &models.PaginatedResponse{
		Data:       albums,
		Page:       1,
		Limit:      10,
		TotalCount: int64(len(albums)),
		TotalPages: 1,
	}, nil
}

func (s *UserService) UpdateTopAlbums(ctx context.Context, userID uuid.UUID, req models.UpdateTopItemsRequest) (*models.PaginatedResponse, error) {
	if len(req.IDs) > 10 {
		return nil, &models.ErrorResponse{Message: "maximum 10 albums allowed"}
	}

	if err := db.UpdateUserTopAlbums(ctx, userID, req.IDs); err != nil {
		return nil, err
	}

	albums, err := db.GetTopAlbumsByIDs(ctx, req.IDs)
	if err != nil {
		return nil, err
	}

	return &models.PaginatedResponse{
		Data:       albums,
		Page:       1,
		Limit:      10,
		TotalCount: int64(len(albums)),
		TotalPages: 1,
	}, nil
}

func (s *UserService) GetTopSongs(ctx context.Context, userID uuid.UUID) (*models.PaginatedResponse, error) {
	user, err := db.GetUserByID(ctx, userID)
	if err != nil {
		return nil, err
	}

	songs, err := db.GetTopSongsByIDs(ctx, user.Top10Songs)
	if err != nil {
		return nil, err
	}

	return &models.PaginatedResponse{
		Data:       songs,
		Page:       1,
		Limit:      10,
		TotalCount: int64(len(songs)),
		TotalPages: 1,
	}, nil
}

func (s *UserService) UpdateTopSongs(ctx context.Context, userID uuid.UUID, req models.UpdateTopItemsRequest) (*models.PaginatedResponse, error) {
	if len(req.IDs) > 10 {
		return nil, &models.ErrorResponse{Message: "maximum 10 songs allowed"}
	}

	if err := db.UpdateUserTopSongs(ctx, userID, req.IDs); err != nil {
		return nil, err
	}

	songs, err := db.GetTopSongsByIDs(ctx, req.IDs)
	if err != nil {
		return nil, err
	}

	return &models.PaginatedResponse{
		Data:       songs,
		Page:       1,
		Limit:      10,
		TotalCount: int64(len(songs)),
		TotalPages: 1,
	}, nil
}

func (s *UserService) GetTopArtists(ctx context.Context, userID uuid.UUID) (*models.PaginatedResponse, error) {
	user, err := db.GetUserByID(ctx, userID)
	if err != nil {
		return nil, err
	}

	artists, err := db.GetTopArtistsByIDs(ctx, user.Top10Artists)
	if err != nil {
		return nil, err
	}

	return &models.PaginatedResponse{
		Data:       artists,
		Page:       1,
		Limit:      10,
		TotalCount: int64(len(artists)),
		TotalPages: 1,
	}, nil
}

func (s *UserService) UpdateTopArtists(ctx context.Context, userID uuid.UUID, req models.UpdateTopItemsRequest) (*models.PaginatedResponse, error) {
	if len(req.IDs) > 10 {
		return nil, &models.ErrorResponse{Message: "maximum 10 artists allowed"}
	}

	if err := db.UpdateUserTopArtists(ctx, userID, req.IDs); err != nil {
		return nil, err
	}

	artists, err := db.GetTopArtistsByIDs(ctx, req.IDs)
	if err != nil {
		return nil, err
	}

	return &models.PaginatedResponse{
		Data:       artists,
		Page:       1,
		Limit:      10,
		TotalCount: int64(len(artists)),
		TotalPages: 1,
	}, nil
}

func (s *UserService) CreateListening(ctx context.Context, userID uuid.UUID, req models.CreateListeningRequest) (*models.ListeningRecord, error) {
	now := time.Now()
	listenedAt := now
	if req.ListenedAt != nil {
		listenedAt = *req.ListenedAt
	}

	listenedYear := int16(listenedAt.Year())
	listenedMonth := int16(listenedAt.Month())

	return db.CreateListeningRecord(ctx, userID, req.SongID, req.AlbumID, req.ArtistID,
		listenedAt, listenedYear, listenedMonth, req.Source, req.DurationSeconds)
}

func (s *UserService) ListListening(ctx context.Context, userID uuid.UUID, page, limit int) (*models.PaginatedResponse, error) {
	offset := (page - 1) * limit

	records, total, err := db.ListListeningRecords(ctx, userID, int32(limit), int32(offset))
	if err != nil {
		return nil, err
	}

	totalPages := int(math.Ceil(float64(total) / float64(limit)))
	return &models.PaginatedResponse{
		Data:       records,
		Page:       page,
		Limit:      limit,
		TotalCount: total,
		TotalPages: totalPages,
	}, nil
}

func (s *UserService) GetMonthlySet(ctx context.Context, userID uuid.UUID, year, month int16) (*models.MonthlyListeningSet, error) {
	return db.GetMonthlySet(ctx, userID, year, month)
}

func (s *UserService) UpsertMonthlySet(ctx context.Context, userID uuid.UUID, year, month int16, req models.UpdateMonthlySetRequest) (*models.MonthlyListeningSet, error) {
	return db.UpsertMonthlySet(ctx, userID, year, month, req.Songs, req.Notes)
}
