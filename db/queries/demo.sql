-- name: DemoGet :one
SELECT *
FROM demos
WHERE id = $1;

-- name: DemoGetByUser :many
SELECT *
FROM demos
WHERE user_id = $1 AND deleted_at IS NULL
ORDER BY created_at DESC;

-- name: DemoGetByStatus :many
SELECT *
FROM demos
WHERE status = $1 AND deleted_at IS NULL
ORDER BY created_at ASC;

-- name: DemoGetByStatusUpdateAtomic :many
WITH cte AS (
  SELECT d.id
  FROM demos d
  WHERE d.status = sqlc.arg('old_status')
  ORDER BY d.attempts, d.created_at
  FOR UPDATE SKIP LOCKED
  LIMIT sqlc.arg('amount')
)
UPDATE demos
SET
  status = sqlc.arg('new_status'),
  attempts = attempts + 1,
  status_updated_at = NOW()
WHERE id in (SELECT id from cte)
RETURNING *;

-- name: DemoCreate :one
INSERT INTO demos (user_id, source, source_id, status, file_id)
VALUES ($1, $2, $3, 'queued_parse', $4)
RETURNING id;

-- name: DemoUpdateStatus :exec
UPDATE demos
SET
  status = $2,
  error = $3,
  status_updated_at = NOW()
WHERE id = $1;

-- name: DemoUpdateFile :exec
UPDATE demos
SET file_id = $2
WHERE id = $1;

-- name: DemoResetStatusAll :exec
UPDATE demos
SET 
  status = sqlc.arg('new_status'),
  status_updated_at = NOW()
WHERE
  status = sqlc.arg('old_status');

-- name: DemoDelete :exec
UPDATE demos
SET deleted_at = NOW()
WHERE id = $1;
