-- name: GetAlbumsByYearMonth :many
SELECT a.id, a.title, a.artist_id, a.release_date, a.release_year, a.release_month,
       a.release_day, a.release_date_precision, a.album_type, a.cover_art_url,
       a.cover_art_thumbnail, a.mbid, a.spotify_id, a.label, a.track_count,
       a.runtime_seconds, a.genres, a.wikipedia_url, a.description,
       a.bandcamp_url, a.primary_format, a.original_label,
       a.discogs_master_id, a.genius_pageviews, a.is_outstanding,
       a.outstanding_score, a.created_at, a.updated_at,
       ar.name AS artist_name, ar.image_url AS artist_image_url
FROM albums a
JOIN artists ar ON a.artist_id = ar.id
WHERE a.release_year = $1 AND a.release_month = $2
ORDER BY a.is_outstanding DESC, a.outstanding_score DESC NULLS LAST, a.release_day ASC NULLS LAST, a.title ASC
LIMIT $3 OFFSET $4;

-- name: GetAlbumsByYear :many
SELECT a.id, a.title, a.artist_id, a.release_date, a.release_year, a.release_month,
       a.release_day, a.release_date_precision, a.album_type, a.cover_art_url,
       a.cover_art_thumbnail, a.mbid, a.spotify_id, a.label, a.track_count,
       a.runtime_seconds, a.genres, a.wikipedia_url, a.description,
       a.bandcamp_url, a.primary_format, a.original_label,
       a.discogs_master_id, a.genius_pageviews, a.is_outstanding,
       a.outstanding_score, a.created_at, a.updated_at,
       ar.name AS artist_name, ar.image_url AS artist_image_url
FROM albums a
JOIN artists ar ON a.artist_id = ar.id
WHERE a.release_year = $1
ORDER BY a.is_outstanding DESC, a.outstanding_score DESC NULLS LAST, a.release_month ASC NULLS LAST, a.title ASC
LIMIT $2 OFFSET $3;

-- name: GetTimelineYears :many
SELECT release_year, COUNT(*) AS album_count
FROM albums
GROUP BY release_year
ORDER BY release_year DESC;

-- name: GetTimelineMonths :many
SELECT release_year, release_month, COUNT(*) AS album_count
FROM albums
GROUP BY release_year, release_month
ORDER BY release_year DESC, release_month DESC;

-- name: GetTimelineMonthDetail :many
SELECT a.id, a.title, a.artist_id, a.release_date, a.release_year, a.release_month,
       a.release_day, a.release_date_precision, a.album_type, a.cover_art_url,
       a.cover_art_thumbnail, a.label, a.track_count, a.is_outstanding,
       a.outstanding_score, a.created_at, a.updated_at,
       ar.name AS artist_name, ar.image_url AS artist_image_url
FROM albums a
JOIN artists ar ON a.artist_id = ar.id
WHERE a.release_year = $1
  AND (a.release_month = $2 OR $2 IS NULL)
ORDER BY a.is_outstanding DESC, a.outstanding_score DESC NULLS LAST, a.release_day ASC NULLS LAST, a.title ASC
LIMIT $3 OFFSET $4;

-- name: GetAlbumByID :one
SELECT a.id, a.title, a.artist_id, a.release_date, a.release_year, a.release_month,
       a.release_day, a.release_date_precision, a.album_type, a.cover_art_url,
       a.cover_art_thumbnail, a.mbid, a.spotify_id, a.label, a.track_count,
       a.runtime_seconds, a.genres, a.wikipedia_url, a.description,
       a.bandcamp_url, a.primary_format, a.original_label,
       a.discogs_master_id, a.genius_pageviews, a.is_outstanding,
       a.outstanding_score, a.created_at, a.updated_at,
       ar.name AS artist_name, ar.image_url AS artist_image_url, ar.bio AS artist_bio,
       ar.genres AS artist_genres, ar.country AS artist_country
FROM albums a
JOIN artists ar ON a.artist_id = ar.id
WHERE a.id = $1;

-- name: GetAlbumArtists :many
SELECT aa.artist_id, aa.role, aa.position,
       ar.name AS artist_name, ar.image_url AS artist_image_url
FROM album_artists aa
JOIN artists ar ON aa.artist_id = ar.id
WHERE aa.album_id = $1
ORDER BY aa.position ASC;

-- name: GetAlbumSongs :many
SELECT s.id, s.title, s.album_id, s.artist_id, s.track_number, s.disc_number,
       s.duration_seconds, s.spotify_id, s.isrc,
       ar.name AS artist_name
FROM songs s
LEFT JOIN artists ar ON s.artist_id = ar.id
WHERE s.album_id = $1
ORDER BY s.disc_number ASC, s.track_number ASC;

-- name: ListAlbums :many
SELECT a.id, a.title, a.artist_id, a.release_date, a.release_year, a.release_month,
       a.release_day, a.release_date_precision, a.album_type, a.cover_art_url,
       a.cover_art_thumbnail, a.mbid, a.spotify_id, a.label, a.track_count,
       a.runtime_seconds, a.genres, a.wikipedia_url, a.description,
       a.bandcamp_url, a.primary_format, a.original_label,
       a.discogs_master_id, a.genius_pageviews, a.is_outstanding,
       a.outstanding_score, a.created_at, a.updated_at,
       ar.name AS artist_name, ar.image_url AS artist_image_url
FROM albums a
JOIN artists ar ON a.artist_id = ar.id
WHERE
  ($1::SMALLINT IS NULL OR a.release_year = $1) AND
  ($2::SMALLINT IS NULL OR a.release_month = $2) AND
  ($3::BOOLEAN IS NULL OR a.is_outstanding = $3) AND
  ($4::TEXT IS NULL OR a.genres @> to_jsonb($4::TEXT))
ORDER BY
  CASE WHEN $5::TEXT = 'release_date' AND $6::TEXT = 'asc' THEN a.release_date END ASC,
  CASE WHEN $5::TEXT = 'release_date' AND $6::TEXT = 'desc' THEN a.release_date END DESC,
  CASE WHEN $5::TEXT = 'title' AND $6::TEXT = 'asc' THEN a.title END ASC,
  CASE WHEN $5::TEXT = 'title' AND $6::TEXT = 'desc' THEN a.title END DESC,
  a.outstanding_score DESC NULLS LAST, a.title ASC
LIMIT $7 OFFSET $8;

-- name: CountAlbums :one
SELECT COUNT(*) FROM albums a
WHERE
  ($1::SMALLINT IS NULL OR a.release_year = $1) AND
  ($2::SMALLINT IS NULL OR a.release_month = $2) AND
  ($3::BOOLEAN IS NULL OR a.is_outstanding = $3) AND
  ($4::TEXT IS NULL OR a.genres @> to_jsonb($4::TEXT));

-- name: SearchAlbums :many
SELECT a.id, a.title, a.artist_id, a.release_date, a.release_year, a.release_month,
       a.cover_art_url, a.cover_art_thumbnail, a.album_type, a.is_outstanding,
       a.outstanding_score,
       ar.name AS artist_name, ar.image_url AS artist_image_url,
       ts_rank(a.search_vector, plainto_tsquery('english', $1)) AS rank
FROM albums a
JOIN artists ar ON a.artist_id = ar.id
WHERE
  a.search_vector @@ plainto_tsquery('english', $1)
  OR a.title % $1
ORDER BY rank DESC, a.outstanding_score DESC NULLS LAST
LIMIT $2 OFFSET $3;

-- name: CountSearchAlbums :one
SELECT COUNT(*) FROM albums a
WHERE
  a.search_vector @@ plainto_tsquery('english', $1)
  OR a.title % $1;
