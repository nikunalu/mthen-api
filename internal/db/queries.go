package db

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/nik/mthen-api/internal/models"
)

// ---------- Albums ----------

func GetTimelineYears(ctx context.Context) ([]models.TimelineYear, error) {
	rows, err := Pool.Query(ctx,
		`SELECT release_year, COUNT(*) AS album_count
		 FROM albums
		 GROUP BY release_year
		 ORDER BY release_year DESC`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var results []models.TimelineYear
	for rows.Next() {
		var ty models.TimelineYear
		if err := rows.Scan(&ty.Year, &ty.AlbumCount); err != nil {
			return nil, err
		}
		results = append(results, ty)
	}
	return results, rows.Err()
}

func GetTimelineMonths(ctx context.Context) ([]models.TimelineMonth, error) {
	rows, err := Pool.Query(ctx,
		`SELECT release_year, release_month, COUNT(*) AS album_count
		 FROM albums
		 GROUP BY release_year, release_month
		 ORDER BY release_year DESC, release_month DESC`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var results []models.TimelineMonth
	for rows.Next() {
		var tm models.TimelineMonth
		if err := rows.Scan(&tm.Year, &tm.Month, &tm.AlbumCount); err != nil {
			return nil, err
		}
		results = append(results, tm)
	}
	return results, rows.Err()
}

func GetAlbumsByYear(ctx context.Context, year int16) ([]models.AlbumWithArtist, error) {
	rows, err := Pool.Query(ctx,
		`SELECT a.id, a.title, a.artist_id, a.release_date, a.release_year, a.release_month,
		        a.release_day, a.release_date_precision, a.album_type, a.cover_art_url,
		        a.cover_art_thumbnail, a.label, a.track_count, a.is_outstanding,
		        a.outstanding_score, a.genres, a.created_at, a.updated_at,
		        ar.name AS artist_name, ar.image_url AS artist_image_url
		 FROM albums a
		 JOIN artists ar ON a.artist_id = ar.id
		 WHERE a.release_year = $1
		 ORDER BY a.is_outstanding DESC, a.release_month ASC, a.release_day ASC NULLS LAST`, year)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var albums []models.AlbumWithArtist
	for rows.Next() {
		var a models.AlbumWithArtist
		if err := rows.Scan(
			&a.ID, &a.Title, &a.ArtistID, &a.ReleaseDate, &a.ReleaseYear, &a.ReleaseMonth,
			&a.ReleaseDay, &a.ReleaseDatePrecision, &a.AlbumType, &a.CoverArtURL,
			&a.CoverArtThumbnail, &a.Label, &a.TrackCount, &a.IsOutstanding,
			&a.OutstandingScore, &a.Genres, &a.CreatedAt, &a.UpdatedAt,
			&a.ArtistName, &a.ArtistImageURL,
		); err != nil {
			return nil, err
		}
		albums = append(albums, a)
	}
	return albums, rows.Err()
}

func GetTimelineMonthDetail(ctx context.Context, year, month int16, limit, offset int32) ([]models.AlbumWithArtist, int64, error) {
	var total int64
	err := Pool.QueryRow(ctx,
		`SELECT COUNT(*) FROM albums WHERE release_year = $1 AND ($2 = 0 OR release_month = $2)`,
		year, month).Scan(&total)
	if err != nil {
		return nil, 0, err
	}

	rows, err := Pool.Query(ctx,
		`SELECT a.id, a.title, a.artist_id, a.release_date, a.release_year, a.release_month,
		        a.release_day, a.release_date_precision, a.album_type, a.cover_art_url,
		        a.cover_art_thumbnail, a.label, a.track_count, a.is_outstanding,
		        a.outstanding_score, a.created_at, a.updated_at,
		        ar.name AS artist_name, ar.image_url AS artist_image_url
		 FROM albums a
		 JOIN artists ar ON a.artist_id = ar.id
		 WHERE a.release_year = $1
		   AND ($2 = 0 OR a.release_month = $2)
		 ORDER BY a.is_outstanding DESC, a.outstanding_score DESC NULLS LAST,
		          a.release_day ASC NULLS LAST, a.title ASC
		 LIMIT $3 OFFSET $4`, year, month, limit, offset)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var albums []models.AlbumWithArtist
	for rows.Next() {
		var a models.AlbumWithArtist
		if err := rows.Scan(
			&a.ID, &a.Title, &a.ArtistID, &a.ReleaseDate, &a.ReleaseYear, &a.ReleaseMonth,
			&a.ReleaseDay, &a.ReleaseDatePrecision, &a.AlbumType, &a.CoverArtURL,
			&a.CoverArtThumbnail, &a.Label, &a.TrackCount, &a.IsOutstanding,
			&a.OutstandingScore, &a.CreatedAt, &a.UpdatedAt,
			&a.ArtistName, &a.ArtistImageURL,
		); err != nil {
			return nil, 0, err
		}
		albums = append(albums, a)
	}
	return albums, total, rows.Err()
}

func GetAlbumByID(ctx context.Context, id uuid.UUID) (*models.AlbumDetail, error) {
	var a models.AlbumDetail
	err := Pool.QueryRow(ctx,
		`SELECT a.id, a.title, a.artist_id, a.release_date, a.release_year, a.release_month,
		        a.release_day, a.release_date_precision, a.album_type, a.cover_art_url,
		        a.cover_art_thumbnail, a.mbid, a.spotify_id, a.label, a.track_count,
		        a.runtime_seconds, a.genres, a.wikipedia_url, a.description,
		        a.bandcamp_url, a.primary_format, a.original_label,
		        a.discogs_master_id, a.genius_pageviews, a.is_outstanding,
		        a.outstanding_score, a.created_at, a.updated_at,
		        ar.name AS artist_name, ar.image_url AS artist_image_url,
		        ar.bio AS artist_bio, ar.genres AS artist_genres, ar.country AS artist_country
		 FROM albums a
		 JOIN artists ar ON a.artist_id = ar.id
		 WHERE a.id = $1`, id).Scan(
		&a.ID, &a.Title, &a.ArtistID, &a.ReleaseDate, &a.ReleaseYear, &a.ReleaseMonth,
		&a.ReleaseDay, &a.ReleaseDatePrecision, &a.AlbumType, &a.CoverArtURL,
		&a.CoverArtThumbnail, &a.MBID, &a.SpotifyID, &a.Label, &a.TrackCount,
		&a.RuntimeSeconds, &a.Genres, &a.WikipediaURL, &a.Description,
		&a.BandcampURL, &a.PrimaryFormat, &a.OriginalLabel,
		&a.DiscogsMasterID, &a.GeniusPageviews, &a.IsOutstanding,
		&a.OutstandingScore, &a.CreatedAt, &a.UpdatedAt,
		&a.ArtistName, &a.ArtistImageURL,
		&a.ArtistBio, &a.ArtistGenres, &a.ArtistCountry,
	)
	if err != nil {
		return nil, err
	}
	return &a, nil
}

func GetAlbumSongs(ctx context.Context, albumID uuid.UUID) ([]models.SongWithArtist, error) {
	rows, err := Pool.Query(ctx,
		`SELECT s.id, s.title, s.album_id, s.artist_id, s.track_number, s.disc_number,
		        s.duration_seconds, s.spotify_id, s.isrc,
		        ar.name AS artist_name
		 FROM songs s
		 LEFT JOIN artists ar ON s.artist_id = ar.id
		 WHERE s.album_id = $1
		 ORDER BY s.disc_number ASC, s.track_number ASC`, albumID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var songs []models.SongWithArtist
	for rows.Next() {
		var s models.SongWithArtist
		if err := rows.Scan(&s.ID, &s.Title, &s.AlbumID, &s.ArtistID, &s.TrackNumber,
			&s.DiscNumber, &s.DurationSeconds, &s.SpotifyID, &s.ISRC, &s.ArtistName); err != nil {
			return nil, err
		}
		songs = append(songs, s)
	}
	return songs, rows.Err()
}

func GetAlbumArtists(ctx context.Context, albumID uuid.UUID) ([]models.AlbumArtistCredit, error) {
	rows, err := Pool.Query(ctx,
		`SELECT aa.artist_id, ar.name, aa.role, aa.position, ar.image_url
		 FROM album_artists aa
		 JOIN artists ar ON aa.artist_id = ar.id
		 WHERE aa.album_id = $1
		 ORDER BY aa.position ASC`, albumID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var credits []models.AlbumArtistCredit
	for rows.Next() {
		var c models.AlbumArtistCredit
		if err := rows.Scan(&c.ArtistID, &c.Name, &c.Role, &c.Position, &c.ImageURL); err != nil {
			return nil, err
		}
		credits = append(credits, c)
	}
	return credits, rows.Err()
}

func ListAlbums(ctx context.Context, year, month *int16, isOutstanding *bool, genre *string, sortBy, sortOrder string, limit, offset int32) ([]models.AlbumWithArtist, int64, error) {
	// Build dynamic query for filtering
	countQuery := `SELECT COUNT(*) FROM albums a WHERE 1=1`
	args := []interface{}{}
	argIdx := 1

	if year != nil {
		countQuery += fmt.Sprintf(` AND a.release_year = $%d`, argIdx)
		args = append(args, *year)
		argIdx++
	}
	if month != nil {
		countQuery += fmt.Sprintf(` AND a.release_month = $%d`, argIdx)
		args = append(args, *month)
		argIdx++
	}
	if isOutstanding != nil {
		countQuery += fmt.Sprintf(` AND a.is_outstanding = $%d`, argIdx)
		args = append(args, *isOutstanding)
		argIdx++
	}
	if genre != nil && *genre != "" {
		countQuery += fmt.Sprintf(` AND a.genres @> to_jsonb($%d::TEXT)`, argIdx)
		args = append(args, *genre)
		argIdx++
	}

	var total int64
	err := Pool.QueryRow(ctx, countQuery, args...).Scan(&total)
	if err != nil {
		return nil, 0, err
	}

	// Build data query
	dataQuery := `SELECT a.id, a.title, a.artist_id, a.release_date, a.release_year, a.release_month,
		a.release_day, a.release_date_precision, a.album_type, a.cover_art_url,
		a.cover_art_thumbnail, a.mbid, a.spotify_id, a.label, a.track_count,
		a.runtime_seconds, a.genres, a.wikipedia_url, a.description,
		a.bandcamp_url, a.primary_format, a.original_label,
		a.discogs_master_id, a.genius_pageviews, a.is_outstanding,
		a.outstanding_score, a.created_at, a.updated_at,
		ar.name AS artist_name, ar.image_url AS artist_image_url
	 FROM albums a
	 JOIN artists ar ON a.artist_id = ar.id
	 WHERE 1=1`

	dataArgs := []interface{}{}
	dataIdx := 1

	if year != nil {
		dataQuery += fmt.Sprintf(` AND a.release_year = $%d`, dataIdx)
		dataArgs = append(dataArgs, *year)
		dataIdx++
	}
	if month != nil {
		dataQuery += fmt.Sprintf(` AND a.release_month = $%d`, dataIdx)
		dataArgs = append(dataArgs, *month)
		dataIdx++
	}
	if isOutstanding != nil {
		dataQuery += fmt.Sprintf(` AND a.is_outstanding = $%d`, dataIdx)
		dataArgs = append(dataArgs, *isOutstanding)
		dataIdx++
	}
	if genre != nil && *genre != "" {
		dataQuery += fmt.Sprintf(` AND a.genres @> to_jsonb($%d::TEXT)`, dataIdx)
		dataArgs = append(dataArgs, *genre)
		dataIdx++
	}

	dataQuery += ` ORDER BY a.outstanding_score DESC NULLS LAST, a.title ASC`
	dataQuery += fmt.Sprintf(` LIMIT $%d OFFSET $%d`, dataIdx, dataIdx+1)
	dataArgs = append(dataArgs, limit, offset)

	rows, err := Pool.Query(ctx, dataQuery, dataArgs...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var albums []models.AlbumWithArtist
	for rows.Next() {
		var a models.AlbumWithArtist
		if err := rows.Scan(
			&a.ID, &a.Title, &a.ArtistID, &a.ReleaseDate, &a.ReleaseYear, &a.ReleaseMonth,
			&a.ReleaseDay, &a.ReleaseDatePrecision, &a.AlbumType, &a.CoverArtURL,
			&a.CoverArtThumbnail, &a.MBID, &a.SpotifyID, &a.Label, &a.TrackCount,
			&a.RuntimeSeconds, &a.Genres, &a.WikipediaURL, &a.Description,
			&a.BandcampURL, &a.PrimaryFormat, &a.OriginalLabel,
			&a.DiscogsMasterID, &a.GeniusPageviews, &a.IsOutstanding,
			&a.OutstandingScore, &a.CreatedAt, &a.UpdatedAt,
			&a.ArtistName, &a.ArtistImageURL,
		); err != nil {
			return nil, 0, err
		}
		albums = append(albums, a)
	}
	return albums, total, rows.Err()
}

// ---------- Artists ----------

func GetArtistByID(ctx context.Context, id uuid.UUID) (*models.Artist, error) {
	var a models.Artist
	err := Pool.QueryRow(ctx,
		`SELECT id, name, sort_name, mbid, spotify_id, image_url, bio, genres,
		        country, bandcamp_url, created_at, updated_at
		 FROM artists WHERE id = $1`, id).Scan(
		&a.ID, &a.Name, &a.SortName, &a.MBID, &a.SpotifyID, &a.ImageURL,
		&a.Bio, &a.Genres, &a.Country, &a.BandcampURL, &a.CreatedAt, &a.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	return &a, nil
}

func GetArtistDiscography(ctx context.Context, artistID uuid.UUID, limit, offset int32) ([]models.AlbumSummary, error) {
	rows, err := Pool.Query(ctx,
		`SELECT a.id, a.title, a.release_year, a.release_month, a.release_day,
		        a.album_type, a.cover_art_url, a.cover_art_thumbnail,
		        a.is_outstanding, a.outstanding_score,
		        ar.name AS artist_name
		 FROM albums a
		 JOIN artists ar ON a.artist_id = ar.id
		 WHERE a.artist_id = $1
		 ORDER BY a.release_year DESC, a.release_month DESC NULLS LAST, a.release_day DESC NULLS LAST
		 LIMIT $2 OFFSET $3`, artistID, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var albums []models.AlbumSummary
	for rows.Next() {
		var a models.AlbumSummary
		if err := rows.Scan(&a.ID, &a.Title, &a.ReleaseYear, &a.ReleaseMonth,
			&a.ReleaseDay, &a.AlbumType, &a.CoverArtURL, &a.CoverArtThumbnail,
			&a.IsOutstanding, &a.OutstandingScore, &a.ArtistName); err != nil {
			return nil, err
		}
		albums = append(albums, a)
	}
	return albums, rows.Err()
}

func GetArtistReleaseYears(ctx context.Context, artistID uuid.UUID) ([]int16, error) {
	rows, err := Pool.Query(ctx,
		`SELECT DISTINCT a.release_year
		 FROM albums a
		 WHERE a.artist_id = $1
		 ORDER BY a.release_year DESC`, artistID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var years []int16
	for rows.Next() {
		var y int16
		if err := rows.Scan(&y); err != nil {
			return nil, err
		}
		years = append(years, y)
	}
	return years, rows.Err()
}

func GetArtistRelatedArtists(ctx context.Context, artistID uuid.UUID) ([]models.RelatedArtist, error) {
	rows, err := Pool.Query(ctx,
		`SELECT ar.id, ar.name, ar.image_url, arr.relationship_type
		 FROM artist_relationships arr
		 JOIN artists ar ON
		   CASE WHEN arr.from_artist_id = $1 THEN arr.to_artist_id = ar.id
		        WHEN arr.to_artist_id = $1 THEN arr.from_artist_id = ar.id
		   END
		 WHERE arr.from_artist_id = $1 OR arr.to_artist_id = $1
		 ORDER BY ar.name`, artistID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var related []models.RelatedArtist
	for rows.Next() {
		var r models.RelatedArtist
		if err := rows.Scan(&r.ID, &r.Name, &r.ImageURL, &r.RelationshipType); err != nil {
			return nil, err
		}
		related = append(related, r)
	}
	return related, rows.Err()
}

func ListArtists(ctx context.Context, query string, limit, offset int32) ([]models.Artist, int64, error) {
	var total int64
	var totalErr error

	if query != "" {
		totalErr = Pool.QueryRow(ctx,
			`SELECT COUNT(*) FROM artists
			 WHERE name ILIKE '%' || $1 || '%'
			    OR name % $1
			    OR search_vector @@ plainto_tsquery('english', $1)`, query).Scan(&total)
	} else {
		totalErr = Pool.QueryRow(ctx, `SELECT COUNT(*) FROM artists`).Scan(&total)
	}
	if totalErr != nil {
		return nil, 0, totalErr
	}

	var pgRows pgx.Rows
	var scanErr error

	if query != "" {
		pgRows, scanErr = Pool.Query(ctx,
			`SELECT id, name, sort_name, mbid, spotify_id, image_url, bio, genres,
			        country, bandcamp_url, created_at, updated_at
			 FROM artists
			 WHERE name ILIKE '%' || $1 || '%'
			    OR name % $1
			    OR search_vector @@ plainto_tsquery('english', $1)
			 ORDER BY similarity(name, $1) DESC, name ASC
			 LIMIT $2 OFFSET $3`, query, limit, offset)
	} else {
		pgRows, scanErr = Pool.Query(ctx,
			`SELECT id, name, sort_name, mbid, spotify_id, image_url, bio, genres,
			        country, bandcamp_url, created_at, updated_at
			 FROM artists
			 ORDER BY name ASC
			 LIMIT $1 OFFSET $2`, limit, offset)
	}
	if scanErr != nil {
		return nil, 0, scanErr
	}
	defer pgRows.Close()

	var artists []models.Artist
	for pgRows.Next() {
		var a models.Artist
		if err := pgRows.Scan(&a.ID, &a.Name, &a.SortName, &a.MBID, &a.SpotifyID,
			&a.ImageURL, &a.Bio, &a.Genres, &a.Country, &a.BandcampURL,
			&a.CreatedAt, &a.UpdatedAt); err != nil {
			return nil, 0, err
		}
		artists = append(artists, a)
	}
	return artists, total, pgRows.Err()
}

// ---------- Users / Auth ----------

func CreateUser(ctx context.Context, email, passwordHash, displayName string) (*models.UserProfile, error) {
	var u models.UserProfile
	err := Pool.QueryRow(ctx,
		`INSERT INTO user_profiles (email, password_hash, display_name)
		 VALUES ($1, $2, $3)
		 RETURNING id, email, display_name, avatar_url, bio, joined_at,
		           top_10_songs, top_10_albums, top_10_artists, created_at, updated_at`,
		email, passwordHash, displayName,
	).Scan(&u.ID, &u.Email, &u.DisplayName, &u.AvatarURL, &u.Bio, &u.JoinedAt,
		&u.Top10Songs, &u.Top10Albums, &u.Top10Artists, &u.CreatedAt, &u.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return &u, nil
}

type UserWithPassword struct {
	models.UserProfile
	PasswordHash string
}

func GetUserByEmail(ctx context.Context, email string) (*UserWithPassword, error) {
	var u UserWithPassword
	err := Pool.QueryRow(ctx,
		`SELECT id, email, password_hash, display_name, avatar_url, bio,
		        joined_at, top_10_songs, top_10_albums, top_10_artists,
		        created_at, updated_at
		 FROM user_profiles WHERE email = $1`, email,
	).Scan(&u.ID, &u.Email, &u.PasswordHash, &u.DisplayName, &u.AvatarURL, &u.Bio,
		&u.JoinedAt, &u.Top10Songs, &u.Top10Albums, &u.Top10Artists,
		&u.CreatedAt, &u.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return &u, nil
}

func GetUserByID(ctx context.Context, id uuid.UUID) (*models.UserProfile, error) {
	var u models.UserProfile
	err := Pool.QueryRow(ctx,
		`SELECT id, email, display_name, avatar_url, bio,
		        joined_at, top_10_songs, top_10_albums, top_10_artists,
		        created_at, updated_at
		 FROM user_profiles WHERE id = $1`, id,
	).Scan(&u.ID, &u.Email, &u.DisplayName, &u.AvatarURL, &u.Bio,
		&u.JoinedAt, &u.Top10Songs, &u.Top10Albums, &u.Top10Artists,
		&u.CreatedAt, &u.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return &u, nil
}

func UpdateUserProfile(ctx context.Context, id uuid.UUID, displayName, avatarURL, bio *string) (*models.UserProfile, error) {
	var u models.UserProfile
	err := Pool.QueryRow(ctx,
		`UPDATE user_profiles
		 SET display_name = COALESCE($2, display_name),
		     avatar_url = COALESCE($3, avatar_url),
		     bio = COALESCE($4, bio),
		     updated_at = now()
		 WHERE id = $1
		 RETURNING id, email, display_name, avatar_url, bio, joined_at,
		           top_10_songs, top_10_albums, top_10_artists, created_at, updated_at`,
		id, displayName, avatarURL, bio,
	).Scan(&u.ID, &u.Email, &u.DisplayName, &u.AvatarURL, &u.Bio, &u.JoinedAt,
		&u.Top10Songs, &u.Top10Albums, &u.Top10Artists, &u.CreatedAt, &u.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return &u, nil
}

func UpdateUserTopAlbums(ctx context.Context, id uuid.UUID, albumIDs []uuid.UUID) error {
	_, err := Pool.Exec(ctx,
		`UPDATE user_profiles SET top_10_albums = $2, updated_at = now() WHERE id = $1`,
		id, albumIDs)
	return err
}

func UpdateUserTopSongs(ctx context.Context, id uuid.UUID, songIDs []uuid.UUID) error {
	_, err := Pool.Exec(ctx,
		`UPDATE user_profiles SET top_10_songs = $2, updated_at = now() WHERE id = $1`,
		id, songIDs)
	return err
}

func UpdateUserTopArtists(ctx context.Context, id uuid.UUID, artistIDs []uuid.UUID) error {
	_, err := Pool.Exec(ctx,
		`UPDATE user_profiles SET top_10_artists = $2, updated_at = now() WHERE id = $1`,
		id, artistIDs)
	return err
}

func GetTopAlbumsByIDs(ctx context.Context, ids []uuid.UUID) ([]models.AlbumSummary, error) {
	if len(ids) == 0 {
		return []models.AlbumSummary{}, nil
	}
	rows, err := Pool.Query(ctx,
		`SELECT a.id, a.title, a.release_year, a.cover_art_url,
		        a.cover_art_thumbnail, a.is_outstanding,
		        ar.name AS artist_name
		 FROM albums a
		 JOIN artists ar ON a.artist_id = ar.id
		 WHERE a.id = ANY($1)`, ids)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var albums []models.AlbumSummary
	for rows.Next() {
		var a models.AlbumSummary
		if err := rows.Scan(&a.ID, &a.Title, &a.ReleaseYear,
			&a.CoverArtURL, &a.CoverArtThumbnail, &a.IsOutstanding,
			&a.ArtistName); err != nil {
			return nil, err
		}
		albums = append(albums, a)
	}
	return albums, rows.Err()
}

func GetTopSongsByIDs(ctx context.Context, ids []uuid.UUID) ([]models.SongSummary, error) {
	if len(ids) == 0 {
		return []models.SongSummary{}, nil
	}
	rows, err := Pool.Query(ctx,
		`SELECT s.id, s.title, s.album_id, s.artist_id, s.track_number, s.duration_seconds,
		        a.title AS album_title, a.cover_art_url AS album_cover_url,
		        ar.name AS artist_name
		 FROM songs s
		 JOIN albums a ON s.album_id = a.id
		 LEFT JOIN artists ar ON s.artist_id = ar.id
		 WHERE s.id = ANY($1)`, ids)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var songs []models.SongSummary
	for rows.Next() {
		var s models.SongSummary
		if err := rows.Scan(&s.ID, &s.Title, &s.AlbumID, &s.ArtistID, &s.TrackNumber,
			&s.DurationSeconds, &s.AlbumTitle, &s.AlbumCoverURL, &s.ArtistName); err != nil {
			return nil, err
		}
		songs = append(songs, s)
	}
	return songs, rows.Err()
}

func GetTopArtistsByIDs(ctx context.Context, ids []uuid.UUID) ([]models.ArtistSummary, error) {
	if len(ids) == 0 {
		return []models.ArtistSummary{}, nil
	}
	rows, err := Pool.Query(ctx,
		`SELECT id, name, image_url, genres, country
		 FROM artists WHERE id = ANY($1)`, ids)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var artists []models.ArtistSummary
	for rows.Next() {
		var a models.ArtistSummary
		if err := rows.Scan(&a.ID, &a.Name, &a.ImageURL, &a.Genres, &a.Country); err != nil {
			return nil, err
		}
		artists = append(artists, a)
	}
	return artists, rows.Err()
}

func CreateListeningRecord(ctx context.Context, userID, songID, albumID, artistID uuid.UUID, listenedAt time.Time, listenedYear, listenedMonth int16, source *string, durationSeconds *int32) (*models.ListeningRecord, error) {
	var lr models.ListeningRecord
	err := Pool.QueryRow(ctx,
		`INSERT INTO listening_records (user_id, song_id, album_id, artist_id,
		                                listened_at, listened_year, listened_month,
		                                source, duration_seconds)
		 VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
		 RETURNING id, user_id, song_id, album_id, artist_id,
		           listened_at, listened_year, listened_month,
		           source, duration_seconds, created_at`,
		userID, songID, albumID, artistID, listenedAt, listenedYear, listenedMonth,
		source, durationSeconds,
	).Scan(&lr.ID, &lr.UserID, &lr.SongID, &lr.AlbumID, &lr.ArtistID,
		&lr.ListenedAt, &lr.ListenedYear, &lr.ListenedMonth,
		&lr.Source, &lr.DurationSeconds, &lr.CreatedAt)
	if err != nil {
		return nil, err
	}
	return &lr, nil
}

func ListListeningRecords(ctx context.Context, userID uuid.UUID, limit, offset int32) ([]models.ListeningRecordDetail, int64, error) {
	var total int64
	err := Pool.QueryRow(ctx,
		`SELECT COUNT(*) FROM listening_records WHERE user_id = $1`, userID).Scan(&total)
	if err != nil {
		return nil, 0, err
	}

	rows, err := Pool.Query(ctx,
		`SELECT lr.id, lr.user_id, lr.song_id, lr.album_id, lr.artist_id,
		        lr.listened_at, lr.listened_year, lr.listened_month,
		        lr.source, lr.duration_seconds, lr.created_at,
		        s.title AS song_title,
		        a.title AS album_title, a.cover_art_url AS album_cover_url,
		        ar.name AS artist_name
		 FROM listening_records lr
		 JOIN songs s ON lr.song_id = s.id
		 JOIN albums a ON lr.album_id = a.id
		 JOIN artists ar ON lr.artist_id = ar.id
		 WHERE lr.user_id = $1
		 ORDER BY lr.listened_at DESC
		 LIMIT $2 OFFSET $3`, userID, limit, offset)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var records []models.ListeningRecordDetail
	for rows.Next() {
		var d models.ListeningRecordDetail
		if err := rows.Scan(&d.ID, &d.UserID, &d.SongID, &d.AlbumID, &d.ArtistID,
			&d.ListenedAt, &d.ListenedYear, &d.ListenedMonth,
			&d.Source, &d.DurationSeconds, &d.CreatedAt,
			&d.SongTitle, &d.AlbumTitle, &d.AlbumCoverURL, &d.ArtistName); err != nil {
			return nil, 0, err
		}
		records = append(records, d)
	}
	return records, total, rows.Err()
}

func GetMonthlySet(ctx context.Context, userID uuid.UUID, year, month int16) (*models.MonthlyListeningSet, error) {
	var s models.MonthlyListeningSet
	err := Pool.QueryRow(ctx,
		`SELECT id, user_id, year, month, songs, notes, created_at, updated_at
		 FROM monthly_listening_sets
		 WHERE user_id = $1 AND year = $2 AND month = $3`,
		userID, year, month,
	).Scan(&s.ID, &s.UserID, &s.Year, &s.Month, &s.Songs, &s.Notes, &s.CreatedAt, &s.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return &s, nil
}

func UpsertMonthlySet(ctx context.Context, userID uuid.UUID, year, month int16, songs []uuid.UUID, notes *string) (*models.MonthlyListeningSet, error) {
	var s models.MonthlyListeningSet
	err := Pool.QueryRow(ctx,
		`INSERT INTO monthly_listening_sets (user_id, year, month, songs, notes)
		 VALUES ($1, $2, $3, $4, $5)
		 ON CONFLICT (user_id, year, month)
		 DO UPDATE SET songs = EXCLUDED.songs, notes = EXCLUDED.notes, updated_at = now()
		 RETURNING id, user_id, year, month, songs, notes, created_at, updated_at`,
		userID, year, month, songs, notes,
	).Scan(&s.ID, &s.UserID, &s.Year, &s.Month, &s.Songs, &s.Notes, &s.CreatedAt, &s.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return &s, nil
}

// ---------- Search ----------

func UnifiedSearch(ctx context.Context, query string, limit, offset int32) ([]models.SearchResult, error) {
	rows, err := Pool.Query(ctx,
		`SELECT 'artist' AS result_type,
		        a.id::TEXT AS id,
		        a.name AS title,
		        a.name AS subtitle,
		        a.image_url,
		        NULL::SMALLINT AS release_year,
		        NULL::VARCHAR(20) AS album_type,
		        NULL::BOOLEAN AS is_outstanding
		 FROM artists a
		 WHERE a.search_vector @@ plainto_tsquery('english', $1)
		    OR a.name % $1

		 UNION ALL

		 SELECT 'album' AS result_type,
		        al.id::TEXT AS id,
		        al.title AS title,
		        ar.name AS subtitle,
		        al.cover_art_url AS image_url,
		        al.release_year,
		        al.album_type::VARCHAR(20),
		        al.is_outstanding
		 FROM albums al
		 JOIN artists ar ON al.artist_id = ar.id
		 WHERE al.search_vector @@ plainto_tsquery('english', $1)
		    OR al.title % $1

		 ORDER BY
		   CASE WHEN result_type = 'album' AND is_outstanding THEN 1 ELSE 2 END,
		   title ASC
		 LIMIT $2 OFFSET $3`, query, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var results []models.SearchResult
	for rows.Next() {
		var r models.SearchResult
		if err := rows.Scan(&r.ResultType, &r.ID, &r.Title, &r.Subtitle, &r.ImageURL,
			&r.ReleaseYear, &r.AlbumType, &r.IsOutstanding); err != nil {
			return nil, err
		}
		results = append(results, r)
	}
	return results, rows.Err()
}

// ---------- Genres ----------

func ListGenres(ctx context.Context) ([]models.GenreCount, error) {
	rows, err := Pool.Query(ctx,
		`SELECT genre, COUNT(*) AS album_count
		 FROM (
		   SELECT jsonb_array_elements_text(genres) AS genre
		   FROM albums
		   WHERE genres IS NOT NULL AND jsonb_array_length(genres) > 0
		 ) genre_rows
		 GROUP BY genre
		 ORDER BY album_count DESC, genre ASC`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var genres []models.GenreCount
	for rows.Next() {
		var g models.GenreCount
		if err := rows.Scan(&g.Genre, &g.AlbumCount); err != nil {
			return nil, err
		}
		genres = append(genres, g)
	}
	return genres, rows.Err()
}
