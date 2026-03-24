-- +goose Up
ALTER TABLE runs ADD COLUMN duration INTEGER NOT NULL DEFAULT 0;

ALTER TABLE run_stages ADD COLUMN duration INTEGER NOT NULL DEFAULT 0;

-- +goose Down
ALTER TABLE runs DROP COLUMN duration;

ALTER TABLE run_stages DROP COLUMN duration;