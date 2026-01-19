-- +goose Up
-- +goose StatementBegin
CREATE TABLE stats (
  id SERIAL PRIMARY KEY,
  demo_id INTEGER NOT NULL REFERENCES demos (id) ON DELETE CASCADE,
  user_id INTEGER NOT NULL REFERENCES users (id) ON DELETE CASCADE,
  kills INTEGER NOT NULL,
  assists INTEGER NOT NULL,
  deaths INTEGER NOT NULL,

  UNIQUE (demo_id, user_id)
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE stats;
-- +goose StatementEnd
