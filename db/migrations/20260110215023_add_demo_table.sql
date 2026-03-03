-- +goose Up
-- +goose StatementBegin
CREATE TYPE demo_source AS ENUM ('manual', 'steam', 'faceit');

CREATE TYPE demo_status AS ENUM ('queued_download', 'downloading', 'queued_parse', 'parsing', 'queued_render', 'rendering', 'queued_finalize', 'finalizing', 'finished', 'failed');

CREATE TABLE demos (
  id SERIAL PRIMARY KEY,
  source DEMO_SOURCE NOT NULL,
  source_id TEXT NOT NULL,
  source_url TEXT,
  file_id TEXT,
  data_id TEXT,
  status DEMO_STATUS NOT NULL,
  attempts INTEGER NOT NULL DEFAULT 0,
  error TEXT,
  status_updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),

  UNIQUE(souce, source_id)
);

CREATE INDEX IF NOT EXISTS demos_created_at_idx ON demos(created_at);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP INDEX demos_created_at_idx;
DROP TABLE demos;
DROP TYPE demo_status;
DROP TYPE demo_source;
-- +goose StatementEnd
