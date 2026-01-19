-- +goose Up
-- +goose StatementBegin
CREATE TABLE highlights (
  id SERIAL PRIMARY KEY,
  user_id INTEGER NOT NULL REFERENCES users (id) ON DELETE CASCADE,
  demo_id INTEGER NOT NULL REFERENCES demos (id) ON DELETE CASCADE,
  file_id TEXT,
  file_web_id TEXT,
  title TEXT NOT NULL,
  created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE highlight_segments (
  id SERIAL PRIMARY KEY,
  highlight_id INTEGER NOT NULL REFERENCES highlights (id) ON DELETE CASCADE,
  start_tick INTEGER NOT NULL,
  end_tick INTEGER NOT NULL
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE highlight_segments;

DROP TABLE highlights;
-- +goose StatementEnd
