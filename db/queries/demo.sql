-- name: DemoGet :one
SELECT *
FROM demos
WHERE id = $1;

-- name: DemoGetByUser :many
SELECT *
FROM demos
WHERE user_id = $1 AND deleted_at = NULL
ORDER BY created_at DESC;

-- name: DemoCreate :one
INSERT INTO demos (user_id, source, source_id, status, demo_file_id)
VALUES ($1, $2, $3, 'queued_parse', $4)
RETURNING id;

-- name: DemoUpdateStatus :exec
UPDATE demos
SET status = $2, status_updated_at = NOW()
WHERE id = $1;

-- name: DemoUpdateFile :exec
UPDATE demos
SET demo_file_id = $2
WHERE id = $1;

-- name: DemoDelete :exec
UPDATE demos
SET deleted_at = NOW()
WHERE id = $1;
