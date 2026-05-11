-- name: UnifiedSearch :many
-- Returns combined results from artists and albums, ranked by relevance
SELECT 'artist' AS result_type,
       a.id::TEXT AS id,
       a.name AS title,
       a.name AS subtitle,
       a.image_url,
       NULL::SMALLINT AS release_year,
       NULL::VARCHAR(20) AS album_type,
       NULL::BOOLEAN AS is_outstanding,
       (ts_rank(a.search_vector, plainto_tsquery('english', $1)) *
        (1.0 + COALESCE(similarity(a.name, $1), 0)))::REAL AS rank
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
       al.is_outstanding,
       (ts_rank(al.search_vector, plainto_tsquery('english', $1)) *
        (1.0 + COALESCE(similarity(al.title, $1), 0)))::REAL AS rank
FROM albums al
JOIN artists ar ON al.artist_id = ar.id
WHERE al.search_vector @@ plainto_tsquery('english', $1)
   OR al.title % $1

ORDER BY rank DESC
LIMIT $2 OFFSET $3;
