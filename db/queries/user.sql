-- name: UserGet :one
SELECT *
FROM users
WHERE id = $1;

-- name: UserGetByUID :one
SELECT *
FROM users
WHERE uid = $1;

-- name: UserCreate :one
INSERT INTO users (uid, name, display_name, avatar_url)
VALUES ($1, $2, $3, $4)
RETURNING id;

-- name: UserUpdate :exec
UPDATE users
SET name = $2, display_name = $3, avatar_url = $4
WHERE id = $1;
