-- +goose Up
CREATE TABLE
    projects (
        id SERIAL PRIMARY KEY,
        created_at TIMESTAMPTZ NOT NULL DEFAULT NOW (),
        updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW (),
        cron_enabled BOOLEAN NOT NULL DEFAULT FALSE,
        cron_schedule TEXT NOT NULL DEFAULT '',
        dir TEXT NOT NULL DEFAULT '',
        git_url TEXT NOT NULL DEFAULT '',
        git_branch TEXT NOT NULL DEFAULT '',
        name TEXT NOT NULL,
        restart_after BOOLEAN NOT NULL DEFAULT FALSE,
        retention_count INTEGER NOT NULL DEFAULT 0,
        with_root_env BOOLEAN NOT NULL DEFAULT FALSE
    );

CREATE TABLE
    stages (
        project_id INTEGER NOT NULL REFERENCES projects (id) ON UPDATE CASCADE ON DELETE CASCADE,
        num INTEGER NOT NULL,
        script TEXT NOT NULL,
        PRIMARY KEY (project_id, num)
    );

CREATE INDEX idx_stages_project_id ON stages (project_id);

CREATE TABLE
    runs (
        id SERIAL PRIMARY KEY,
        created_at TIMESTAMPTZ DEFAULT NOW (),
        updated_at TIMESTAMPTZ DEFAULT NOW (),
        project_id INTEGER NOT NULL REFERENCES projects (id) ON UPDATE CASCADE ON DELETE CASCADE,
        success BOOLEAN NOT NULL DEFAULT FALSE,
        pending BOOLEAN NOT NULL DEFAULT FALSE,
        processing BOOLEAN NOT NULL DEFAULT FALSE,
        fail_log TEXT NOT NULL DEFAULT ''
    );

CREATE INDEX idx_runs_project_id ON runs (project_id);

CREATE TABLE
    run_stages (
        run_id INTEGER NOT NULL REFERENCES runs (id) ON UPDATE CASCADE ON DELETE CASCADE,
        stage_num INTEGER NOT NULL,
        success BOOLEAN NOT NULL DEFAULT FALSE,
        log TEXT NOT NULL DEFAULT '',
        PRIMARY KEY (run_id, stage_num)
    );

CREATE INDEX idx_run_stages_run_id ON run_stages (run_id);

CREATE TABLE
    run_git_commits (
        run_id INTEGER NOT NULL REFERENCES runs (id) ON UPDATE CASCADE ON DELETE CASCADE,
        num INTEGER NOT NULL,
        hash TEXT NOT NULL,
        subject TEXT NOT NULL,
        PRIMARY KEY (run_id, num)
    );

CREATE INDEX idx_run_git_commits_run_id ON run_git_commits (run_id);

-- +goose Down
DROP TABLE IF EXISTS run_git_commits;

DROP TABLE IF EXISTS runs;

DROP TABLE IF EXISTS stages;

DROP TABLE IF EXISTS projects;

DROP TABLE IF EXISTS run_stages;