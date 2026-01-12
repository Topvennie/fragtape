-- name: HighlightCreate :one
INSERT INTO highlights (demo_id, file_id, title)
VALUES ($1, $2, $3)
RETURNING id;

-- name: HighlightUpdateFile :exec
UPDATE highlights
SET file_id = $2
WHERE id = $1;
