-- name: HighlightGet :one
SELECT *
FROM highlights
WHERE id = $1;

-- name: HighlightGetByDemo :many
SELECT *
FROM highlights
WHERE demo_id = $1
ORDER BY created_at;

-- name: HighlightGetByDemos :many
SELECT *
FROM highlights
WHERE demo_id = ANY($1::int[]);

-- name: HighlightCreate :one
INSERT INTO highlights (user_id, demo_id, title, round, duration_s)
VALUES ($1, $2, $3, $4, $5)
RETURNING id;

-- name: HighlightUpdate :exec
UPDATE highlights
SET 
  demo_id = coalesce(sqlc.narg('demo_id'), demo_id),
  file_id = coalesce(sqlc.narg('file_id'), file_id),
  file_web_id = coalesce(sqlc.narg('file_web_id'), file_web_id),
  title = coalesce(sqlc.narg('title'), title)
WHERE id = $1;

-- name: HighlightDeleteFile :exec
UPDATE highlights
SET file_id = NULL, file_web_id = NULL
WHERE id = $1;
