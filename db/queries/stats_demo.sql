-- name: StatsDemoGetByDemo :one
SELECT *
FROM stats_demos
WHERE demo_id = $1;

-- name: StatsDemoGetByDemos :many
SELECT *
FROM stats_demos
WHERE demo_id = ANY($1::int[]);

-- name: StatsDemoCreate :one
INSERT INTO stats_demos (demo_id, map, rounds_ct, rounds_t)
VALUES ($1, $2, $3, $4)
RETURNING id;

-- name: StatsDemoUpdate :exec
UPDATE stats_demos
SET map = $2, rounds_ct = $3, rounds_ct = $4, rounds_t = $5
WHERE id = $1;
