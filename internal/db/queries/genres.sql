-- name: ListGenres :many
SELECT genre, COUNT(*) AS album_count
FROM (
  SELECT jsonb_array_elements_text(genres) AS genre
  FROM albums
  WHERE genres IS NOT NULL AND jsonb_array_length(genres) > 0
) genre_rows
GROUP BY genre
ORDER BY album_count DESC, genre ASC;

-- name: ListArtistGenres :many
SELECT genre, COUNT(*) AS artist_count
FROM (
  SELECT jsonb_array_elements_text(genres) AS genre
  FROM artists
  WHERE genres IS NOT NULL AND jsonb_array_length(genres) > 0
) genre_rows
GROUP BY genre
ORDER BY artist_count DESC, genre ASC;
