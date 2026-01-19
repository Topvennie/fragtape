-- +goose Up
-- +goose StatementBegin
CREATE TABLE stats_demos (
  id SERIAL PRIMARY KEY,
  demo_id INTEGER NOT NULL REFERENCES demos (id) ON DELETE CASCADE,
  map TEXT NOT NULL,
  rounds_ct INTEGER NOT NULL,
  rounds_t INTEGER NOT NULL
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE stats_demos;
-- +goose StatementEnd
