-- +goose Up
CREATE TABLE
    projects (
        id INTEGER PRIMARY KEY AUTOINCREMENT,
        created_at DATETIME NOT NULL DEFAULT (datetime ('now')),
        updated_at DATETIME NOT NULL DEFAULT (datetime ('now')),
        cron_enabled BOOLEAN NOT NULL DEFAULT 0,
        cron_schedule TEXT NOT NULL DEFAULT '',
        dir TEXT NOT NULL DEFAULT '',
        git_url TEXT NOT NULL DEFAULT '',
        git_branch TEXT NOT NULL DEFAULT '',
        name TEXT NOT NULL,
        restart_after BOOLEAN NOT NULL DEFAULT 0,
        retention_count INTEGER NOT NULL DEFAULT 0,
        with_root_env BOOLEAN NOT NULL DEFAULT 0
    );

CREATE TABLE
    stages (
        project_id INTEGER NOT NULL,
        num INTEGER NOT NULL,
        script TEXT NOT NULL,
        PRIMARY KEY (project_id, num),
        FOREIGN KEY (project_id) REFERENCES projects (id) ON UPDATE CASCADE ON DELETE CASCADE
    );

CREATE INDEX idx_stages_project_id ON stages (project_id);

CREATE TABLE
    runs (
        id INTEGER PRIMARY KEY AUTOINCREMENT,
        created_at DATETIME DEFAULT (datetime ('now')),
        updated_at DATETIME DEFAULT (datetime ('now')),
        project_id INTEGER NOT NULL,
        success BOOLEAN NOT NULL DEFAULT 0,
        pending BOOLEAN NOT NULL DEFAULT 0,
        processing BOOLEAN NOT NULL DEFAULT 0,
        fail_log TEXT NOT NULL DEFAULT '',
        FOREIGN KEY (project_id) REFERENCES projects (id) ON UPDATE CASCADE ON DELETE CASCADE
    );

CREATE INDEX idx_runs_project_id ON runs (project_id);

CREATE TABLE
    run_stages (
        run_id INTEGER NOT NULL,
        stage_num INTEGER NOT NULL,
        success BOOLEAN NOT NULL DEFAULT 0,
        log TEXT NOT NULL DEFAULT '',
        PRIMARY KEY (run_id, stage_num),
        FOREIGN KEY (run_id) REFERENCES runs (id) ON UPDATE CASCADE ON DELETE CASCADE
    );

CREATE INDEX idx_run_stages_run_id ON run_stages (run_id);

CREATE TABLE
    run_git_commits (
        run_id INTEGER NOT NULL,
        num INTEGER NOT NULL,
        hash TEXT NOT NULL,
        subject TEXT NOT NULL,
        PRIMARY KEY (run_id, num),
        FOREIGN KEY (run_id) REFERENCES runs (id) ON UPDATE CASCADE ON DELETE CASCADE
    );

CREATE INDEX idx_run_git_commits_run_id ON run_git_commits (run_id);

-- +goose Down
DROP TABLE run_git_commits;

DROP TABLE run_stages;

DROP TABLE runs;

DROP TABLE stages;

DROP TABLE projects;