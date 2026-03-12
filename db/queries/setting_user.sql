-- name: SettingUserGetByUser :one
SELECT *
FROM setting_user
WHERE user_id = $1;

-- name: SettingUserCreate :one
INSERT INTO setting_user (user_id, steam_match_token, steam_authentication_token, faceit_id)
VALUES ($1, $2, $3, $4)
RETURNING id;

-- name: SettingUserUpdate :exec
UPDATE setting_user
SET
  steam_match_token = $2,
  steam_authentication_token = $3,
  faceit_id = $4
WHERE id = $1;

