-- +goose Up
ALTER TABLE projects ADD COLUMN is_template BOOLEAN NOT NULL DEFAULT 0;

-- +goose Down
ALTER TABLE projects DROP COLUMN is_template;
