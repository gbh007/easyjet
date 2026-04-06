-- +goose Up

CREATE TABLE env_vars_new (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    created_at DATETIME NOT NULL DEFAULT (datetime ('now')),
    updated_at DATETIME NOT NULL DEFAULT (datetime ('now')),
    project_id INTEGER,
    name TEXT NOT NULL,
    value TEXT NOT NULL,
    uses_other_variables BOOLEAN NOT NULL DEFAULT 0,
    UNIQUE (project_id, name),
    FOREIGN KEY (project_id) REFERENCES projects (id) ON UPDATE CASCADE ON DELETE CASCADE
);

INSERT INTO env_vars_new (id, created_at, updated_at, project_id, name, value, uses_other_variables)
SELECT id, created_at, updated_at, project_id, name, value, uses_other_variables
FROM env_vars;

DROP TABLE env_vars;

ALTER TABLE env_vars_new RENAME TO env_vars;

CREATE INDEX idx_env_vars_project_id ON env_vars (project_id);
CREATE INDEX idx_env_vars_name ON env_vars (name);

-- +goose Down
CREATE TABLE env_vars_old (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    created_at DATETIME NOT NULL DEFAULT (datetime ('now')),
    updated_at DATETIME NOT NULL DEFAULT (datetime ('now')),
    project_id INTEGER,
    name TEXT NOT NULL,
    value TEXT NOT NULL,
    uses_other_variables BOOLEAN NOT NULL DEFAULT 0,
    UNIQUE (project_id, name)
);

INSERT INTO env_vars_old (id, created_at, updated_at, project_id, name, value, uses_other_variables)
SELECT id, created_at, updated_at, project_id, name, value, uses_other_variables
FROM env_vars;

DROP TABLE env_vars;

ALTER TABLE env_vars_old RENAME TO env_vars;

CREATE INDEX idx_env_vars_project_id ON env_vars (project_id);
CREATE INDEX idx_env_vars_name ON env_vars (name);
