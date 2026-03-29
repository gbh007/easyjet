package postgres

import (
	"context"
	"fmt"

	"github.com/Masterminds/squirrel"
	"github.com/gbh007/easyjet/internal/core/entity"
	"github.com/jackc/pgx/v5"
)

func (repo Repo) SetProject(ctx context.Context, p entity.Project) (uint, error) {
	tx, err := repo.pool.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		return 0, fmt.Errorf("begin transaction: %w", err)
	}
	defer func() {
		_ = tx.Rollback(ctx)
	}()

	if p.ID > 0 {
		deleteStagesQuery, deleteStagesArgs, err := repo.psql.
			Delete("stages").
			Where(squirrel.Eq{"project_id": p.ID}).
			ToSql()
		if err != nil {
			return 0, fmt.Errorf("build delete stages query: %w", err)
		}

		_, err = tx.Exec(ctx, deleteStagesQuery, deleteStagesArgs...)
		if err != nil {
			return 0, fmt.Errorf("delete stages: %w", err)
		}
	}

	projectID := p.ID

	if p.ID > 0 {
		updateQuery, updateArgs, err := repo.psql.
			Update("projects").
			SetMap(map[string]any{
				"updated_at":      p.UpdatedAt,
				"cron_enabled":    p.CronEnabled,
				"cron_schedule":   p.CronSchedule,
				"dir":             p.Dir,
				"git_url":         p.GitURL,
				"git_branch":      p.GitBranch,
				"name":            p.Name,
				"restart_after":   p.RestartAfter,
				"retention_count": p.RetentionCount,
				"with_root_env":   p.WithRootEnv,
			}).
			Where(squirrel.Eq{"id": p.ID}).
			ToSql()
		if err != nil {
			return 0, fmt.Errorf("build update query: %w", err)
		}

		_, err = tx.Exec(ctx, updateQuery, updateArgs...)
		if err != nil {
			return 0, fmt.Errorf("update project: %w", err)
		}
	} else {
		insertQuery, insertArgs, err := repo.psql.
			Insert("projects").
			SetMap(map[string]any{
				"cron_enabled":    p.CronEnabled,
				"cron_schedule":   p.CronSchedule,
				"dir":             p.Dir,
				"git_url":         p.GitURL,
				"git_branch":      p.GitBranch,
				"name":            p.Name,
				"restart_after":   p.RestartAfter,
				"retention_count": p.RetentionCount,
				"with_root_env":   p.WithRootEnv,
			}).
			Suffix("RETURNING id").
			ToSql()
		if err != nil {
			return 0, fmt.Errorf("build insert query: %w", err)
		}

		err = tx.QueryRow(ctx, insertQuery, insertArgs...).Scan(&projectID)
		if err != nil {
			return 0, fmt.Errorf("insert project: %w", err)
		}
	}

	for _, stage := range p.Stages {
		insertStageQuery, insertStageArgs, err := repo.psql.
			Insert("stages").
			SetMap(map[string]any{
				"project_id": projectID,
				"num":        stage.Number,
				"script":     stage.Script,
			}).
			ToSql()
		if err != nil {
			return 0, fmt.Errorf("build insert stage query: %w", err)
		}

		_, err = tx.Exec(ctx, insertStageQuery, insertStageArgs...)
		if err != nil {
			return 0, fmt.Errorf("insert stage: %w", err)
		}
	}

	if p.ID > 0 {
		deleteEnvVarsQuery, deleteEnvVarsArgs, err := repo.psql.
			Delete("env_vars").
			Where(squirrel.Eq{"project_id": projectID}).
			ToSql()
		if err != nil {
			return 0, fmt.Errorf("build delete env vars query: %w", err)
		}

		_, err = tx.Exec(ctx, deleteEnvVarsQuery, deleteEnvVarsArgs...)
		if err != nil {
			return 0, fmt.Errorf("delete env vars: %w", err)
		}
	}

	for _, ev := range p.EnvVars {
		insertEnvVarQuery, insertEnvVarArgs, err := repo.psql.
			Insert("env_vars").
			SetMap(map[string]any{
				"project_id":           projectID,
				"name":                 ev.Name,
				"value":                ev.Value,
				"uses_other_variables": ev.UsesOtherVariables,
			}).
			ToSql()
		if err != nil {
			return 0, fmt.Errorf("build insert env var query: %w", err)
		}

		_, err = tx.Exec(ctx, insertEnvVarQuery, insertEnvVarArgs...)
		if err != nil {
			return 0, fmt.Errorf("insert env var: %w", err)
		}
	}

	if err := tx.Commit(ctx); err != nil {
		return 0, fmt.Errorf("commit transaction: %w", err)
	}

	return projectID, nil
}
