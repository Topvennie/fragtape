-- +goose Up
-- +goose StatementBegin
CREATE TABLE stats_demos (
  id SERIAL PRIMARY KEY,
  demo_id INTEGER NOT NULL REFERENCES demos (id) ON DELETE CASCADE,
  map TEXT,
  rounds_ct INTEGER,
  rounds_t INTEGER,

  UNIQUE (demo_id)
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE stats_demos;
-- +goose StatementEnd
