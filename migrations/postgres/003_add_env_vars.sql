-- +goose Up
CREATE TABLE
    env_vars (
        id SERIAL PRIMARY KEY,
        created_at TIMESTAMP NOT NULL DEFAULT NOW(),
        updated_at TIMESTAMP NOT NULL DEFAULT NOW(),
        project_id INTEGER,
        name TEXT NOT NULL,
        value TEXT NOT NULL,
        uses_other_variables BOOLEAN NOT NULL DEFAULT FALSE,
        UNIQUE (project_id, name),
        FOREIGN KEY (project_id) REFERENCES projects (id) ON UPDATE CASCADE ON DELETE CASCADE
    );

CREATE INDEX idx_env_vars_project_id ON env_vars (project_id);
CREATE INDEX idx_env_vars_name ON env_vars (name);

-- +goose Down
DROP TABLE env_vars;
