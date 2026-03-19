-- +goose Up
-- +goose StatementBegin
CREATE TABLE setting_user (
  id SERIAL PRIMARY KEY,
  user_id INTEGER NOT NULL REFERENCES users (id) ON DELETE CASCADE,
  steam_match_token TEXT,
  steam_authentication_token TEXT,
  steam_import_old BOOLEAN,
  faceit_id TEXT,
  first_time_wizard BOOLEAN NOT NULL DEFAULT true
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE setting_user;
-- +goose StatementEnd
