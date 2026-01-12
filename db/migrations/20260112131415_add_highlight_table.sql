-- +goose Up
-- +goose StatementBegin
CREATE TABLE highlights (
  id SERIAL PRIMARY KEY,
  demo_id INTEGER NOT NULL REFERENCES demos (id) ON DELETE CASCADE,
  file_id TEXT,
  title TEXT NOT NULL,
  created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE highlights;
-- +goose StatementEnd
