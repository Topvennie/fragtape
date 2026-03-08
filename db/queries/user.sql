-- name: UserGet :one
SELECT *
FROM users
WHERE id = $1;

-- name: UserGetByUid :one
SELECT *
FROM users
WHERE uid = $1;

-- name: UserGetByUids :many
SELECT *
FROM users
WHERE uid = ANY($1::bigint[]);

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

-- name: UserGetAllRealWithSettingLastDemo :many
SELECT
  sqlc.embed(u),
  sqlc.embed(s_u),
  COALESCE(d.id, 0),
  COALESCE(d.source, 'manual'),
  COALESCE(d.source_id, ''),
  d.created_at
FROM users u
LEFT JOIN setting_user s_u ON s_u.user_id = u.id
LEFT JOIN LATERAL (
  SELECT d2.id, d2.source, d2.source_id, d2.created_at
  FROM stats s
  JOIN demos d2 ON d2.id = s.demo_id
  WHERE s.user_id = u.id
  ORDER BY d2.created_at DESC
  LIMIT 1
) d ON true
WHERE u.name != '';

-- name: UserCreate :one
INSERT INTO users (uid, name, display_name, avatar_url, crosshair, admin)
VALUES ($1, $2, $3, $4, $5, NOT EXISTS (SELECT 1 FROM users))
RETURNING id;

-- name: UserUpdate :exec
UPDATE users
SET name = $2, display_name = $3, avatar_url = $4, crosshair = $5, admin = $6
WHERE id = $1;
