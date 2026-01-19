-- name: StatGetByDemos :many
SELECT *
FROM stats
WHERE demo_id = ANY($1::int[]);

-- name: StatCreate :one
INSERT INTO stats (demo_id, user_id, kills, assists, deaths)
VALUES ($1, $2, $3, $4, $5)
RETURNING id;
