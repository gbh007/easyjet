-- +goose Up
ALTER TABLE projects ADD COLUMN is_template BOOLEAN NOT NULL DEFAULT false;

-- +goose Down
ALTER TABLE projects DROP COLUMN is_template;
