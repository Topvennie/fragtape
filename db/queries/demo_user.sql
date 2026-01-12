-- name: DemoUserGetByDemoUser :one
SELECT *
FROM demo_users
WHERE demo_id = $1 AND user_id = $2;

-- name: DemoUserCreate :one
INSERT INTO demo_users (demo_id, user_id)
VALUES ($1, $2)
RETURNING id;

-- name: DemoUserDeleteByDemoUser :exec
UPDATE demo_users
SET deleted_at = NOW()
WHERE demo_id = $1 AND user_id = $2;
