package postgres

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
		).
		From("projects").
		Where(squirrel.Eq{"id": id}).
		ToSql()
	if err != nil {
		return entity.Project{}, fmt.Errorf("build project query: %w", err)
	}

	err = repo.pool.QueryRow(ctx, projectQuery, projectArgs...).Scan(
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

	rows, err := repo.pool.Query(ctx, stagesQuery, stagesArgs...)
	if err != nil {
		return entity.Project{}, fmt.Errorf("query stages: %w", err)
	}
	defer rows.Close()

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

	return p, nil
}
