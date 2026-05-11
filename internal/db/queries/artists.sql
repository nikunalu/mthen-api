-- name: GetArtistByID :one
SELECT id, name, sort_name, mbid, spotify_id, image_url, bio, genres,
       country, bandcamp_url, created_at, updated_at
FROM artists
WHERE id = $1;

-- name: ListArtists :many
SELECT id, name, sort_name, mbid, spotify_id, image_url, bio, genres,
       country, bandcamp_url, created_at, updated_at
FROM artists
WHERE
  ($1::TEXT IS NULL OR name ILIKE '%' || $1 || '%'
   OR name % $1
   OR search_vector @@ plainto_tsquery('english', $1))
ORDER BY
  CASE WHEN $1::TEXT IS NOT NULL THEN
    similarity(name, $1)
  ELSE 0.0
  END DESC,
  name ASC
LIMIT $2 OFFSET $3;

-- name: CountArtists :one
SELECT COUNT(*) FROM artists
WHERE
  ($1::TEXT IS NULL OR name ILIKE '%' || $1 || '%'
   OR name % $1
   OR search_vector @@ plainto_tsquery('english', $1));

-- name: GetArtistDiscography :many
SELECT a.id, a.title, a.release_year, a.release_month, a.release_day,
       a.album_type, a.cover_art_url, a.cover_art_thumbnail,
       a.is_outstanding, a.outstanding_score
FROM albums a
WHERE a.artist_id = $1
ORDER BY a.release_year DESC, a.release_month DESC NULLS LAST, a.release_day DESC NULLS LAST
LIMIT $2 OFFSET $3;

-- name: GetArtistReleaseYears :many
SELECT DISTINCT a.release_year
FROM albums a
WHERE a.artist_id = $1
ORDER BY a.release_year DESC;

-- name: GetArtistRelatedArtists :many
SELECT ar.id, ar.name, ar.image_url, arr.relationship_type
FROM artist_relationships arr
JOIN artists ar ON
  CASE WHEN arr.from_artist_id = $1 THEN arr.to_artist_id = ar.id
       WHEN arr.to_artist_id = $1 THEN arr.from_artist_id = ar.id
  END
WHERE arr.from_artist_id = $1 OR arr.to_artist_id = $1
ORDER BY ar.name;

-- name: SearchArtists :many
SELECT id, name, image_url, country, genres,
       ts_rank(search_vector, plainto_tsquery('english', $1)) *
       (1.0 + COALESCE(similarity(name, $1), 0)) AS rank
FROM artists
WHERE
  search_vector @@ plainto_tsquery('english', $1)
  OR name % $1
ORDER BY rank DESC, name ASC
LIMIT $2 OFFSET $3;

-- name: CountSearchArtists :one
SELECT COUNT(*) FROM artists
WHERE
  search_vector @@ plainto_tsquery('english', $1)
  OR name % $1;
