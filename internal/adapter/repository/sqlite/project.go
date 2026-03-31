package sqlite

import (
	"context"
	"fmt"
	"slices"

	"github.com/Masterminds/squirrel"
	"github.com/gbh007/easyjet/internal/core/entity"
)

func (repo Repo) Project(ctx context.Context, id uint) (entity.Project, error) {
	var p entity.Project

	projectQuery, projectArgs, err := repo.psql.
		Select(
			"id",
			"created_at",
			"updated_at",
			"cron_enabled",
			"cron_schedule",
			"dir",
			"git_url",
			"git_branch",
			"name",
			"restart_after",
			"retention_count",
			"with_root_env",
			"is_template",
		).
		From("projects").
		Where(squirrel.Eq{"id": id}).
		ToSql()
	if err != nil {
		return entity.Project{}, fmt.Errorf("build project query: %w", err)
	}

	err = repo.db.QueryRowContext(ctx, projectQuery, projectArgs...).Scan(
		&p.ID,
		&p.CreatedAt,
		&p.UpdatedAt,
		&p.CronEnabled,
		&p.CronSchedule,
		&p.Dir,
		&p.GitURL,
		&p.GitBranch,
		&p.Name,
		&p.RestartAfter,
		&p.RetentionCount,
		&p.WithRootEnv,
		&p.IsTemplate,
	)
	if err != nil {
		return entity.Project{}, fmt.Errorf("query project: %w", err)
	}

	stagesQuery, stagesArgs, err := repo.psql.
		Select("project_id", "num", "script").
		From("stages").
		Where(squirrel.Eq{"project_id": id}).
		OrderBy("num ASC").
		ToSql()
	if err != nil {
		return entity.Project{}, fmt.Errorf("build stages query: %w", err)
	}

	rows, err := repo.db.QueryContext(ctx, stagesQuery, stagesArgs...)
	if err != nil {
		return entity.Project{}, fmt.Errorf("query stages: %w", err)
	}
	defer func() { _ = rows.Close() }()

	for rows.Next() {
		var stage entity.ProjectStage
		if err := rows.Scan(&stage.ProjectID, &stage.Number, &stage.Script); err != nil {
			return entity.Project{}, fmt.Errorf("scan stage: %w", err)
		}
		p.Stages = append(p.Stages, stage)
	}

	if err := rows.Err(); err != nil {
		return entity.Project{}, fmt.Errorf("iterate stages: %w", err)
	}

	slices.SortStableFunc(p.Stages, func(a, b entity.ProjectStage) int {
		return a.Number - b.Number
	})

	envVarsQuery, envVarsArgs, err := repo.psql.
		Select("id", "created_at", "updated_at", "project_id", "name", "value", "uses_other_variables").
		From("env_vars").
		Where(squirrel.Eq{"project_id": id}).
		OrderBy("id ASC").
		ToSql()
	if err != nil {
		return entity.Project{}, fmt.Errorf("build env vars query: %w", err)
	}

	rows, err = repo.db.QueryContext(ctx, envVarsQuery, envVarsArgs...)
	if err != nil {
		return entity.Project{}, fmt.Errorf("query env vars: %w", err)
	}
	defer func() { _ = rows.Close() }()

	for rows.Next() {
		var ev entity.EnvironmentVariable
		if err := rows.Scan(&ev.ID, &ev.CreatedAt, &ev.UpdatedAt, &ev.ProjectID, &ev.Name, &ev.Value, &ev.UsesOtherVariables); err != nil {
			return entity.Project{}, fmt.Errorf("scan env var: %w", err)
		}
		p.EnvVars = append(p.EnvVars, ev)
	}

	if err := rows.Err(); err != nil {
		return entity.Project{}, fmt.Errorf("iterate env vars: %w", err)
	}

	return p, nil
}
