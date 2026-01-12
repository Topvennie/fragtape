-- +goose Up
-- +goose StatementBegin
CREATE TYPE demo_source AS ENUM ('manual', 'steam', 'faceit');

CREATE TYPE demo_status AS ENUM ('queued_parse', 'parsing', 'queued_render', 'rendering', 'rendered', 'completed', 'failed');

CREATE TABLE demos (
  id SERIAL PRIMARY KEY,
  user_id INTEGER NOT NULL REFERENCES users (id) ON DELETE CASCADE,
  source DEMO_SOURCE NOT NULL,
  source_id TEXT,
  file_id TEXT,
  status DEMO_STATUS NOT NULL,
  attempts INTEGER NOT NULL DEFAULT 0,
  error TEXT,
  created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  status_updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  deleted_at TIMESTAMPTZ DEFAULT NULL
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE demos;
DROP TYPE demo_status;
DROP TYPE demo_source;
-- +goose StatementEnd
