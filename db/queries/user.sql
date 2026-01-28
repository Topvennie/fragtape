-- name: UserGet :one
SELECT *
FROM users
WHERE id = $1;

-- name: UserGetByUid :one
SELECT *
FROM users
WHERE uid = $1;

-- name: UserGetByIds :many
SELECT *
FROM users
WHERE id = ANY($1::int[]);

-- name: UserGetAdmin :many
SELECT *
FROM users
WHERE admin;

-- name: UserGetFiltered :many
SELECT
  sqlc.embed(u),
  COUNT(*) OVER()::bigint AS total_count
FROM users u
WHERE
  (u.name ILIKE '%' || @name::text || '%' OR u.display_name ILIKE '%' || @name::text || '%') AND
  (u.admin = @admin::bool OR NOT @filter_admin::bool) AND
  (u.name != '' OR NOT @filter_real::bool)
ORDER BY u.name, u.display_name
LIMIT $1 OFFSET $2;

-- name: UserCreate :one
INSERT INTO users (uid, name, display_name, avatar_url, crosshair, admin)
VALUES ($1, $2, $3, $4, $5, NOT EXISTS (SELECT 1 FROM users))
RETURNING id;

-- name: UserUpdate :exec
UPDATE users
SET name = $2, display_name = $3, avatar_url = $4, crosshair = $5, admin = $6
WHERE id = $1;
