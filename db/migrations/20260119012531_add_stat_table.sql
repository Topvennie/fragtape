-- +goose Up
-- +goose StatementBegin
CREATE TYPE result AS ENUM ('win', 'loss', 'tie');

CREATE TYPE team AS ENUM ('ct', 't');

CREATE TABLE stats (
  id SERIAL PRIMARY KEY,
  demo_id INTEGER NOT NULL REFERENCES demos (id) ON DELETE CASCADE,
  user_id INTEGER NOT NULL REFERENCES users (id) ON DELETE CASCADE,
  result RESULT NOT NULL,
  start_team TEAM NOT NULL,
  kills INTEGER NOT NULL,
  assists INTEGER NOT NULL,
  deaths INTEGER NOT NULL,

  UNIQUE (demo_id, user_id)
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE stats;
DROP TYPE team
DROP TYPE result;
-- +goose StatementEnd
