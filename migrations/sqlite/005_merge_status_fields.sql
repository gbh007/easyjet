-- +goose Up
ALTER TABLE runs ADD COLUMN status TEXT NOT NULL DEFAULT 'pending';

UPDATE runs SET status = 'pending' WHERE pending = 1;
UPDATE runs SET status = 'processing' WHERE processing = 1;
UPDATE runs SET status = 'success' WHERE success = 1 AND pending = 0 AND processing = 0;
UPDATE runs SET status = 'failed' WHERE success = 0 AND pending = 0 AND processing = 0;

ALTER TABLE runs DROP COLUMN success;
ALTER TABLE runs DROP COLUMN pending;
ALTER TABLE runs DROP COLUMN processing;

-- +goose Down

ALTER TABLE runs ADD COLUMN success BOOLEAN NOT NULL DEFAULT 0;
ALTER TABLE runs ADD COLUMN pending BOOLEAN NOT NULL DEFAULT 0;
ALTER TABLE runs ADD COLUMN processing BOOLEAN NOT NULL DEFAULT 0;

UPDATE runs SET success = 0 WHERE status = 'success';
UPDATE runs SET pending = 0 WHERE status = 'pending';
UPDATE runs SET processing = 0 WHERE status = 'processing';

ALTER TABLE runs DROP COLUMN status;