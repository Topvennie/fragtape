-- name: StatsDemoCreate :one
INSERT INTO stats_demos (demo_id, map, rounds_ct, rounds_t)
VALUES ($1, $2, $3, $4)
RETURNING id;
