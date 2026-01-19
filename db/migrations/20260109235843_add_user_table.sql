-- +goose Up
-- +goose StatementBegin
CREATE TABLE users (
  id SERIAL PRIMARY KEY,
  uid INTEGER NOT NULL,
  name TEXT,
  display_name TEXT NOT NULL,
  avatar_url TEXT,
  crosshair TEXT,

  UNIQUE (uid)
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE users;
-- +goose StatementEnd

