-- name: HighlightGet :one
SELECT *
FROM highlights
WHERE id = $1;

-- name: HighlightGetByDemo :many
SELECT *
FROM highlights
WHERE demo_id = $1
ORDER BY created_at;

-- name: HighlightGetByDemoPopulated :many
SELECT sqlc.embed(h), sqlc.embed(s)
FROM highlights h
LEFT JOIN highlight_segments s ON s.highlight_id = h.id
WHERE demo_id = $1
ORDER BY h.created_at, s.start_tick;

-- name: HighlightGetByDemos :many
SELECT *
FROM highlights
WHERE demo_id = ANY($1::int[]);

-- name: HighlightCreate :one
INSERT INTO highlights (user_id, demo_id, title, round, kills, duration_s)
VALUES ($1, $2, $3, $4, $5, $6)
RETURNING id;

-- name: HighlightUpdate :exec
UPDATE highlights
SET 
  demo_id = coalesce(sqlc.narg('demo_id'), demo_id),
  file_id = coalesce(sqlc.narg('file_id'), file_id)
WHERE id = $1;

-- name: HighlightDeleteByDemo :exec
DELETE FROM highlights
WHERE demo_id = $1;

-- name: HighlightDeleteFile :exec
UPDATE highlights
SET file_id = NULL
WHERE id = $1;
