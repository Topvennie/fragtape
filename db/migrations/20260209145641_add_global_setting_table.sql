-- +goose Up
-- +goose StatementBegin
CREATE TABLE setting_global (
  id SERIAL PRIMARY KEY,
  demo_upload BOOLEAN NOT NULL DEFAULT false,
  custom_criteria BOOLEAN NOT NULL DEFAULT false,
  chat_command BOOLEAN NOT NULL DEFAULT false,
  chat_trigger TEXT NOT NULL DEFAULT 'fragtape'
);

INSERT INTO setting_global (id)
VALUES (1);

SELECT setval(pg_get_serial_sequence('setting_global', 'id'), GREATEST((SELECT MAX(id) FROM setting_global), 1));
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE setting_global;
-- +goose StatementEnd
