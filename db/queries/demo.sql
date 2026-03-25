-- name: DemoGet :one
SELECT *
FROM demos
WHERE id = $1;

-- name: DemoGetByUserFiltered :many
SELECT
  sqlc.embed(d),
  COUNT(*) OVER()::bigint AS total_count
FROM demos d
LEFT JOIN stats s ON s.demo_id = d.id
WHERE 
  (s.user_id = @user_id) AND
  (NOT @filter_source::bool OR d.source = @source::DEMO_SOURCE) AND
  (NOT @filter_result::bool OR s.result = @result::RESULT) AND
  (NOT @filter_played_at_start::bool OR d.played_at >= @played_at_start::timestamptz) AND
  (NOT @filter_played_at_end::bool OR d.played_at <= @played_at_end::timestamptz) AND
  (NOT @filter_has_highlight::bool OR EXISTS (SELECT 1 FROM highlights h WHERE h.demo_id = d.id AND h.user_id = @user_id) = @has_highlight::bool)
ORDER BY d.played_at DESC
LIMIT $1 OFFSET $2;

-- name: DemoGetBySourceSourceID :one
SELECT *
FROM demos
WHERE source = $1 AND source_id = $2;

-- name: DemoGetByStatus :many
SELECT *
FROM demos
WHERE status = $1
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
  error = NULL,
  attempts = attempts + 1,
  status_updated_at = NOW()
WHERE id in (SELECT id from cte)
RETURNING *;

-- name: DemoCreate :one
INSERT INTO demos (source, source_id, source_url, status, file_id, played_at)
VALUES ($1, $2, $3, $4, $5, $6)
RETURNING id;

-- name: DemoUpdateStatus :exec
UPDATE demos
SET
  status = $2,
  error = $3,
  attempts = $4,
  expired = $5,
  status_updated_at = NOW()
WHERE id = $1;

-- name: DemoUpdateFile :exec
UPDATE demos
SET file_id = $2
WHERE id = $1;

-- name: DemoUpdateData :exec
UPDATE demos
SET data_id = $2
WHERE id = $1;

-- name: DemoResetStatusAll :exec
UPDATE demos
SET 
  status = sqlc.arg('new_status'),
  status_updated_at = NOW()
WHERE
  status = sqlc.arg('old_status');

