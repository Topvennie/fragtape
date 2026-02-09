-- name: SettingGlobalGet :one
SELECT *
FROM setting_global
WHERE id = 1;

-- name: SettingGlobalUpdate :exec
UPDATE setting_global
SET
  demo_upload = coalesce(sqlc.narg('demo_upload'), demo_upload),
  custom_criteria = coalesce(sqlc.narg('custom_criteria'), custom_criteria),
  chat_command = coalesce(sqlc.narg('chat_command'), chat_command),
  chat_trigger = coalesce(sqlc.narg('chat_trigger'), chat_trigger)
WHERE id = 1;
