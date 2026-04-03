-- +goose Up
ALTER TABLE runs ADD COLUMN status VARCHAR(20) NOT NULL DEFAULT 'pending';

UPDATE runs SET status = 'pending' WHERE pending = true;
UPDATE runs SET status = 'processing' WHERE processing = true;
UPDATE runs SET status = 'success' WHERE success = true AND pending = false AND processing = false;
UPDATE runs SET status = 'failed' WHERE success = false AND pending = false AND processing = false;

ALTER TABLE runs DROP COLUMN success;
ALTER TABLE runs DROP COLUMN pending;
ALTER TABLE runs DROP COLUMN processing;

-- +goose Down

ALTER TABLE runs ADD COLUMN success BOOLEAN NOT NULL DEFAULT false;
ALTER TABLE runs ADD COLUMN pending BOOLEAN NOT NULL DEFAULT false;
ALTER TABLE runs ADD COLUMN processing BOOLEAN NOT NULL DEFAULT false;

UPDATE runs SET success = true WHERE status = 'success';
UPDATE runs SET pending = true WHERE status = 'pending';
UPDATE runs SET processing = true WHERE status = 'processing';

ALTER TABLE runs DROP COLUMN status;
