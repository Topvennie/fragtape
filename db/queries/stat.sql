-- name: StatGetByDemo :many
SELECT *
FROM stats
WHERE demo_id = $1;

-- name: StatGetByDemos :many
SELECT *
FROM stats
WHERE demo_id = ANY($1::int[]);

-- name: StatCreateUpdateAtomic :one
INSERT INTO stats (demo_id, user_id, result, start_team, kills, assists, deaths)
VALUES ($1, $2, $3, $4, $5, $6, $7)
ON CONFLICT (demo_id, user_id)
DO UPDATE SET
  result = $3, start_team = $4, kills = $5, assists = $6, deaths = $7
RETURNING id;

-- name: StatCreateNoConflict :one
INSERT INTO stats (demo_id, user_id, result, start_team, kills, assists, deaths)
VALUES ($1, $2, $3, $4, $5, $6, $7)
ON CONFLICT (demo_id, user_id)
DO UPDATE SET
  user_id = stats.user_id
RETURNING id;
