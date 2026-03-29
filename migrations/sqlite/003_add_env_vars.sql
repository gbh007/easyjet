-- +goose Up
CREATE TABLE
    env_vars (
        id INTEGER PRIMARY KEY AUTOINCREMENT,
        created_at DATETIME NOT NULL DEFAULT (datetime ('now')),
        updated_at DATETIME NOT NULL DEFAULT (datetime ('now')),
        project_id INTEGER,
        name TEXT NOT NULL,
        value TEXT NOT NULL,
        uses_other_variables BOOLEAN NOT NULL DEFAULT 0,
        UNIQUE (project_id, name)
    );

CREATE INDEX idx_env_vars_project_id ON env_vars (project_id);
CREATE INDEX idx_env_vars_name ON env_vars (name);

-- +goose Down
DROP TABLE env_vars;
