package models

import (
	"time"

	"github.com/google/uuid"
)

// --- Enums ---

type DatePrecision string

const (
	DatePrecisionYear  DatePrecision = "year"
	DatePrecisionMonth DatePrecision = "month"
	DatePrecisionDay   DatePrecision = "day"
)

type AlbumType string

const (
	AlbumTypeAlbum       AlbumType = "album"
	AlbumTypeEP          AlbumType = "ep"
	AlbumTypeSingle      AlbumType = "single"
	AlbumTypeCompilation AlbumType = "compilation"
	AlbumTypeLive        AlbumType = "live"
	AlbumTypeSoundtrack  AlbumType = "soundtrack"
	AlbumTypeOther       AlbumType = "other"
)

type ArtistRole string

const (
	ArtistRolePrimary   ArtistRole = "primary"
	ArtistRoleFeatured  ArtistRole = "featured"
	ArtistRoleProducer  ArtistRole = "producer"
	ArtistRoleRemixer   ArtistRole = "remixer"
	ArtistRoleComposer  ArtistRole = "composer"
)

// --- Core Entities ---

type Artist struct {
	ID          uuid.UUID              `json:"id"`
	Name        string                 `json:"name"`
	SortName    *string                `json:"sort_name,omitempty"`
	MBID        *string                `json:"mbid,omitempty"`
	SpotifyID   *string                `json:"spotify_id,omitempty"`
	ImageURL    *string                `json:"image_url,omitempty"`
	Bio         *string                `json:"bio,omitempty"`
	Genres      []string               `json:"genres"`
	Country     *string                `json:"country,omitempty"`
	BandcampURL *string                `json:"bandcamp_url,omitempty"`
	CreatedAt   time.Time              `json:"created_at"`
	UpdatedAt   time.Time              `json:"updated_at"`
}

type Album struct {
	ID                  uuid.UUID      `json:"id"`
	Title               string         `json:"title"`
	ArtistID            uuid.UUID      `json:"artist_id"`
	ReleaseDate         *time.Time     `json:"release_date,omitempty"`
	ReleaseYear         int16          `json:"release_year"`
	ReleaseMonth        *int16         `json:"release_month,omitempty"`
	ReleaseDay          *int16         `json:"release_day,omitempty"`
	ReleaseDatePrecision DatePrecision  `json:"release_date_precision"`
	AlbumType           AlbumType      `json:"album_type"`
	CoverArtURL         *string        `json:"cover_art_url,omitempty"`
	CoverArtThumbnail   *string        `json:"cover_art_thumbnail,omitempty"`
	MBID                *string        `json:"mbid,omitempty"`
	SpotifyID           *string        `json:"spotify_id,omitempty"`
	Label               *string        `json:"label,omitempty"`
	TrackCount          *int16         `json:"track_count,omitempty"`
	RuntimeSeconds      *int32         `json:"runtime_seconds,omitempty"`
	Genres              []string       `json:"genres,omitempty"`
	WikipediaURL        *string        `json:"wikipedia_url,omitempty"`
	Description         *string        `json:"description,omitempty"`
	BandcampURL         *string        `json:"bandcamp_url,omitempty"`
	PrimaryFormat       *string        `json:"primary_format,omitempty"`
	OriginalLabel       *string        `json:"original_label,omitempty"`
	DiscogsMasterID     *int32         `json:"discogs_master_id,omitempty"`
	GeniusPageviews     *int32         `json:"genius_pageviews,omitempty"`
	IsOutstanding       bool           `json:"is_outstanding"`
	OutstandingScore    *float32       `json:"outstanding_score,omitempty"`
	CreatedAt           time.Time      `json:"created_at"`
	UpdatedAt           time.Time      `json:"updated_at"`
}

type Song struct {
	ID             uuid.UUID  `json:"id"`
	Title          string     `json:"title"`
	AlbumID        uuid.UUID  `json:"album_id"`
	ArtistID       *uuid.UUID `json:"artist_id,omitempty"`
	TrackNumber    *int16     `json:"track_number,omitempty"`
	DiscNumber     *int16     `json:"disc_number,omitempty"`
	DurationSeconds *int32    `json:"duration_seconds,omitempty"`
	SpotifyID      *string    `json:"spotify_id,omitempty"`
	MBRecordingID  *string    `json:"mb_recording_id,omitempty"`
	ISRC           *string    `json:"isrc,omitempty"`
	CreatedAt      time.Time  `json:"created_at"`
}

type AlbumArtist struct {
	AlbumID     uuid.UUID  `json:"album_id"`
	ArtistID    uuid.UUID  `json:"artist_id"`
	Role        ArtistRole `json:"role"`
	Position    int16      `json:"position"`
}

// --- View / Enriched types ---

type AlbumWithArtist struct {
	Album
	ArtistName      string  `json:"artist_name"`
	ArtistImageURL  *string `json:"artist_image_url,omitempty"`
}

type AlbumDetail struct {
	AlbumWithArtist
	ArtistBio     *string  `json:"artist_bio,omitempty"`
	ArtistGenres  []string `json:"artist_genres,omitempty"`
	ArtistCountry *string  `json:"artist_country,omitempty"`
	Songs         []SongWithArtist `json:"songs,omitempty"`
	Credits       []AlbumArtistCredit `json:"credits,omitempty"`
}

type AlbumArtistCredit struct {
	ArtistID    uuid.UUID  `json:"artist_id"`
	Name        string     `json:"name"`
	Role        ArtistRole `json:"role"`
	Position    int16      `json:"position"`
	ImageURL    *string    `json:"image_url,omitempty"`
}

type SongWithArtist struct {
	Song
	ArtistName *string `json:"artist_name,omitempty"`
}

type ArtistDetail struct {
	Artist
	Discography    []AlbumSummary     `json:"discography,omitempty"`
	ReleaseYears   []int16            `json:"release_years,omitempty"`
	RelatedArtists []RelatedArtist    `json:"related_artists,omitempty"`
}

type AlbumSummary struct {
	ID                uuid.UUID  `json:"id"`
	Title             string     `json:"title"`
	ArtistName        string     `json:"artist_name"`
	ReleaseYear       int16      `json:"release_year"`
	ReleaseMonth      *int16     `json:"release_month,omitempty"`
	ReleaseDay        *int16     `json:"release_day,omitempty"`
	AlbumType         AlbumType  `json:"album_type"`
	CoverArtURL       *string    `json:"cover_art_url,omitempty"`
	CoverArtThumbnail *string    `json:"cover_art_thumbnail,omitempty"`
	IsOutstanding     bool       `json:"is_outstanding"`
	OutstandingScore  *float32   `json:"outstanding_score,omitempty"`
}

type RelatedArtist struct {
	ID               uuid.UUID `json:"id"`
	Name             string    `json:"name"`
	ImageURL         *string   `json:"image_url,omitempty"`
	RelationshipType string    `json:"relationship_type"`
}

type ActiveYearsResult struct {
	ArtistID uuid.UUID `json:"artist_id"`
	Years    []int16   `json:"years"`
}

// --- User types ---

type UserProfile struct {
	ID           uuid.UUID   `json:"id"`
	Email        string      `json:"email"`
	DisplayName  string      `json:"display_name"`
	AvatarURL    *string     `json:"avatar_url,omitempty"`
	Bio          *string     `json:"bio,omitempty"`
	JoinedAt     time.Time   `json:"joined_at"`
	Top10Songs   []uuid.UUID `json:"top_10_songs,omitempty"`
	Top10Albums  []uuid.UUID `json:"top_10_albums,omitempty"`
	Top10Artists []uuid.UUID `json:"top_10_artists,omitempty"`
	CreatedAt    time.Time   `json:"created_at"`
	UpdatedAt    time.Time   `json:"updated_at"`
}

type UserProfileResponse struct {
	ID           uuid.UUID        `json:"id"`
	Email        string           `json:"email"`
	DisplayName  string           `json:"display_name"`
	AvatarURL    *string          `json:"avatar_url,omitempty"`
	Bio          *string          `json:"bio,omitempty"`
	JoinedAt     time.Time        `json:"joined_at"`
	Top10Albums  []AlbumSummary   `json:"top_10_albums,omitempty"`
	Top10Songs   []SongSummary    `json:"top_10_songs,omitempty"`
	Top10Artists []ArtistSummary  `json:"top_10_artists,omitempty"`
}

type SongSummary struct {
	ID             uuid.UUID `json:"id"`
	Title          string    `json:"title"`
	AlbumID        uuid.UUID `json:"album_id"`
	AlbumTitle     string    `json:"album_title"`
	AlbumCoverURL  *string   `json:"album_cover_url,omitempty"`
	ArtistID       *uuid.UUID `json:"artist_id,omitempty"`
	ArtistName     *string   `json:"artist_name,omitempty"`
	TrackNumber    *int16    `json:"track_number,omitempty"`
	DurationSeconds *int32   `json:"duration_seconds,omitempty"`
}

type ArtistSummary struct {
	ID       uuid.UUID `json:"id"`
	Name     string    `json:"name"`
	ImageURL *string   `json:"image_url,omitempty"`
	Genres   []string  `json:"genres,omitempty"`
	Country  *string   `json:"country,omitempty"`
}

type ListeningRecord struct {
	ID              uuid.UUID  `json:"id"`
	UserID          uuid.UUID  `json:"user_id"`
	SongID          uuid.UUID  `json:"song_id"`
	AlbumID         uuid.UUID  `json:"album_id"`
	ArtistID        uuid.UUID  `json:"artist_id"`
	ListenedAt      time.Time  `json:"listened_at"`
	ListenedYear    int16      `json:"listened_year"`
	ListenedMonth   int16      `json:"listened_month"`
	Source          *string    `json:"source,omitempty"`
	DurationSeconds *int32     `json:"duration_seconds,omitempty"`
	CreatedAt       time.Time  `json:"created_at"`
}

type ListeningRecordDetail struct {
	ListeningRecord
	SongTitle    string  `json:"song_title"`
	AlbumTitle   string  `json:"album_title"`
	AlbumCoverURL *string `json:"album_cover_url,omitempty"`
	ArtistName   string  `json:"artist_name"`
}

type MonthlyListeningSet struct {
	ID        uuid.UUID   `json:"id"`
	UserID    uuid.UUID   `json:"user_id"`
	Year      int16       `json:"year"`
	Month     int16       `json:"month"`
	Songs     []uuid.UUID `json:"songs"`
	Notes     *string     `json:"notes,omitempty"`
	CreatedAt time.Time   `json:"created_at"`
	UpdatedAt time.Time   `json:"updated_at"`
}

// --- Timeline types ---

type TimelineYear struct {
	Year       int16 `json:"year"`
	AlbumCount int64 `json:"album_count"`
}

type TimelineMonth struct {
	Year       int16 `json:"year"`
	Month      int16 `json:"month"`
	AlbumCount int64 `json:"album_count"`
}

type TimelineMonthDetail struct {
	Year       int16          `json:"year"`
	Month      int16          `json:"month"`
	TotalCount int64          `json:"total_count"`
	Albums     []AlbumWithArtist `json:"albums"`
}

// --- Search types ---

type SearchResult struct {
	ResultType    string  `json:"result_type"` // "artist" or "album"
	ID            string  `json:"id"`
	Title         string  `json:"title"`
	Subtitle      string  `json:"subtitle"`
	ImageURL      *string `json:"image_url,omitempty"`
	ReleaseYear   *int16  `json:"release_year,omitempty"`
	AlbumType     *string `json:"album_type,omitempty"`
	IsOutstanding *bool   `json:"is_outstanding,omitempty"`
}

// --- Genre types ---

type GenreCount struct {
	Genre      string `json:"genre"`
	AlbumCount int64  `json:"album_count"`
}

type ArtistGenreCount struct {
	Genre       string `json:"genre"`
	ArtistCount int64  `json:"artist_count"`
}

// --- Request / Response types ---

type RegisterRequest struct {
	Email       string `json:"email"`
	Password    string `json:"password"`
	DisplayName string `json:"display_name"`
}

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type AuthResponse struct {
	Token       string              `json:"token"`
	User        UserProfileResponse `json:"user"`
}

type UpdateProfileRequest struct {
	DisplayName *string `json:"display_name,omitempty"`
	AvatarURL   *string `json:"avatar_url,omitempty"`
	Bio         *string `json:"bio,omitempty"`
}

type UpdateTopItemsRequest struct {
	IDs []uuid.UUID `json:"ids"`
}

type CreateListeningRequest struct {
	SongID          uuid.UUID `json:"song_id"`
	AlbumID         uuid.UUID `json:"album_id"`
	ArtistID        uuid.UUID `json:"artist_id"`
	ListenedAt      *time.Time `json:"listened_at,omitempty"`
	Source          *string   `json:"source,omitempty"`
	DurationSeconds *int32    `json:"duration_seconds,omitempty"`
}

type UpdateMonthlySetRequest struct {
	Songs []uuid.UUID `json:"songs"`
	Notes *string     `json:"notes,omitempty"`
}

type PaginatedResponse struct {
	Data       interface{} `json:"data"`
	Page       int         `json:"page"`
	Limit      int         `json:"limit"`
	TotalCount int64       `json:"total_count"`
	TotalPages int         `json:"total_pages"`
}

type ErrorResponse struct {
	Message string `json:"error"`
}

// Error implements the error interface.
func (e *ErrorResponse) Error() string {
	return e.Message
}
