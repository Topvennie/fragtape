-- name: SettingUserGetByUser :one
SELECT *
FROM setting_user
WHERE user_id = $1;

-- name: SettingUserCreate :one
INSERT INTO setting_user (user_id, steam_match_token, steam_import_old, steam_authentication_token, faceit_id)
VALUES ($1, $2, $3, $4, $5)
RETURNING id;

-- name: SettingUserUpdate :exec
UPDATE setting_user
SET
  steam_match_token = $2,
  steam_authentication_token = $3,
  steam_import_old = $4,
  faceit_id = $5,
  first_time_wizard = $6
WHERE id = $1;

