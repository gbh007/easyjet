package postgres

import (
	"context"
	"fmt"
	"time"

	"github.com/Masterminds/squirrel"
	"github.com/gbh007/easyjet/internal/core/entity"
)

func (repo Repo) GlobalEnvVars(ctx context.Context) ([]entity.EnvironmentVariable, error) {
	var envVars []entity.EnvironmentVariable

	query, args, err := repo.psql.
		Select("id", "created_at", "updated_at", "name", "value", "uses_other_variables").
		From("env_vars").
		Where(squirrel.Eq{"project_id": nil}).
		OrderBy("id ASC").
		ToSql()
	if err != nil {
		return nil, fmt.Errorf("build query: %w", err)
	}

	rows, err := repo.pool.Query(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("query env vars: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var ev entity.EnvironmentVariable
		if err := rows.Scan(&ev.ID, &ev.CreatedAt, &ev.UpdatedAt, &ev.Name, &ev.Value, &ev.UsesOtherVariables); err != nil {
			return nil, fmt.Errorf("scan env var: %w", err)
		}
		envVars = append(envVars, ev)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate env vars: %w", err)
	}

	return envVars, nil
}

func (repo Repo) GlobalEnvVar(ctx context.Context, id uint) (entity.EnvironmentVariable, error) {
	var ev entity.EnvironmentVariable

	query, args, err := repo.psql.
		Select("id", "created_at", "updated_at", "name", "value", "uses_other_variables").
		From("env_vars").
		Where(squirrel.Eq{"id": id, "project_id": nil}).
		ToSql()
	if err != nil {
		return entity.EnvironmentVariable{}, fmt.Errorf("build query: %w", err)
	}

	err = repo.pool.QueryRow(ctx, query, args...).Scan(
		&ev.ID,
		&ev.CreatedAt,
		&ev.UpdatedAt,
		&ev.Name,
		&ev.Value,
		&ev.UsesOtherVariables,
	)
	if err != nil {
		return entity.EnvironmentVariable{}, fmt.Errorf("query env var: %w", err)
	}

	return ev, nil
}

func (repo Repo) SetGlobalEnvVar(ctx context.Context, ev entity.EnvironmentVariable) (uint, error) {
	ev.UpdatedAt = time.Now()

	if ev.ID == 0 {
		ev.CreatedAt = time.Now()

		query, args, err := repo.psql.
			Insert("env_vars").
			SetMap(map[string]any{
				"created_at":           ev.CreatedAt,
				"updated_at":           ev.UpdatedAt,
				"name":                 ev.Name,
				"value":                ev.Value,
				"uses_other_variables": ev.UsesOtherVariables,
			}).
			Suffix("RETURNING id").
			ToSql()
		if err != nil {
			return 0, fmt.Errorf("build insert query: %w", err)
		}

		err = repo.pool.QueryRow(ctx, query, args...).Scan(&ev.ID)
		if err != nil {
			return 0, fmt.Errorf("insert env var: %w", err)
		}

		return ev.ID, nil
	}

	query, args, err := repo.psql.
		Update("env_vars").
		SetMap(map[string]any{
			"updated_at":           ev.UpdatedAt,
			"name":                 ev.Name,
			"value":                ev.Value,
			"uses_other_variables": ev.UsesOtherVariables,
		}).
		Where(squirrel.Eq{"id": ev.ID, "project_id": nil}).
		ToSql()
	if err != nil {
		return ev.ID, fmt.Errorf("build update query: %w", err)
	}

	_, err = repo.pool.Exec(ctx, query, args...)
	if err != nil {
		return ev.ID, fmt.Errorf("update env var: %w", err)
	}

	return ev.ID, nil
}

func (repo Repo) DeleteGlobalEnvVar(ctx context.Context, id uint) error {
	query, args, err := repo.psql.
		Delete("env_vars").
		Where(squirrel.Eq{"id": id, "project_id": nil}).
		ToSql()
	if err != nil {
		return fmt.Errorf("build delete query: %w", err)
	}

	_, err = repo.pool.Exec(ctx, query, args...)
	if err != nil {
		return fmt.Errorf("delete env var: %w", err)
	}

	return nil
}
