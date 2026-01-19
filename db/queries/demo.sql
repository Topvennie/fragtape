-- name: DemoGet :one
SELECT *
FROM demos
WHERE id = $1;

-- name: DemoGetByUserPopulated :many
SELECT sqlc.embed(d), sqlc.embed(sd)
FROM demos d
LEFT JOIN stats s ON s.demo_id = d.id
LEFT JOIN stats_demos sd ON sd.demo_id = d.id
WHERE s.user_id = $1
ORDER BY d.created_at DESC;

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
  attempts = attempts + 1,
  status_updated_at = NOW()
WHERE id in (SELECT id from cte)
RETURNING *;

-- name: DemoCreate :one
INSERT INTO demos (source, source_id, file_id)
VALUES ($1, $2, $3)
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

