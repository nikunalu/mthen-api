-- name: CreateUser :one
INSERT INTO user_profiles (email, password_hash, display_name)
VALUES ($1, $2, $3)
RETURNING id, email, display_name, avatar_url, bio, joined_at, created_at, updated_at;

-- name: GetUserByEmail :one
SELECT id, email, password_hash, display_name, avatar_url, bio,
       joined_at, top_10_songs, top_10_albums, top_10_artists,
       created_at, updated_at
FROM user_profiles
WHERE email = $1;

-- name: GetUserByID :one
SELECT id, email, display_name, avatar_url, bio,
       joined_at, top_10_songs, top_10_albums, top_10_artists,
       created_at, updated_at
FROM user_profiles
WHERE id = $1;

-- name: UpdateUserProfile :one
UPDATE user_profiles
SET display_name = COALESCE($2, display_name),
    avatar_url = COALESCE($3, avatar_url),
    bio = COALESCE($4, bio),
    updated_at = now()
WHERE id = $1
RETURNING id, email, display_name, avatar_url, bio, joined_at,
          top_10_songs, top_10_albums, top_10_artists,
          created_at, updated_at;

-- name: UpdateUserTopAlbums :one
UPDATE user_profiles
SET top_10_albums = $2, updated_at = now()
WHERE id = $1
RETURNING id, top_10_albums;

-- name: UpdateUserTopSongs :one
UPDATE user_profiles
SET top_10_songs = $2, updated_at = now()
WHERE id = $1
RETURNING id, top_10_songs;

-- name: UpdateUserTopArtists :one
UPDATE user_profiles
SET top_10_artists = $2, updated_at = now()
WHERE id = $1
RETURNING id, top_10_artists;

-- name: GetTopAlbumsByIDs :many
SELECT a.id, a.title, a.artist_id, a.release_year, a.cover_art_url,
       a.cover_art_thumbnail, a.is_outstanding,
       ar.name AS artist_name
FROM albums a
JOIN artists ar ON a.artist_id = ar.id
WHERE a.id = ANY($1::UUID[]);

-- name: GetTopSongsByIDs :many
SELECT s.id, s.title, s.album_id, s.artist_id, s.track_number, s.duration_seconds,
       a.title AS album_title, a.cover_art_url AS album_cover_url,
       ar.name AS artist_name
FROM songs s
JOIN albums a ON s.album_id = a.id
LEFT JOIN artists ar ON s.artist_id = ar.id
WHERE s.id = ANY($1::UUID[]);

-- name: GetTopArtistsByIDs :many
SELECT id, name, image_url, genres, country
FROM artists
WHERE id = ANY($1::UUID[]);

-- name: CreateListeningRecord :one
INSERT INTO listening_records (user_id, song_id, album_id, artist_id,
                               listened_at, listened_year, listened_month,
                               source, duration_seconds)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
RETURNING id, user_id, song_id, album_id, artist_id, listened_at,
          listened_year, listened_month, source, duration_seconds, created_at;

-- name: ListListeningRecords :many
SELECT lr.id, lr.user_id, lr.song_id, lr.album_id, lr.artist_id,
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
LIMIT $2 OFFSET $3;

-- name: CountListeningRecords :one
SELECT COUNT(*) FROM listening_records WHERE user_id = $1;

-- name: GetMonthlySet :one
SELECT id, user_id, year, month, songs, notes, created_at, updated_at
FROM monthly_listening_sets
WHERE user_id = $1 AND year = $2 AND month = $3;

-- name: UpsertMonthlySet :one
INSERT INTO monthly_listening_sets (user_id, year, month, songs, notes)
VALUES ($1, $2, $3, $4, $5)
ON CONFLICT (user_id, year, month)
DO UPDATE SET songs = EXCLUDED.songs, notes = EXCLUDED.notes, updated_at = now()
RETURNING id, user_id, year, month, songs, notes, created_at, updated_at;
