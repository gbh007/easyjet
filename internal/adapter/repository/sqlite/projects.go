package sqlite

import (
	"context"
	"fmt"
	"slices"

	"github.com/Masterminds/squirrel"
	"github.com/gbh007/easyjet/internal/core/entity"
)

func (repo Repo) Projects(ctx context.Context) ([]entity.Project, error) {
	projectsQuery, projectsArgs, err := repo.psql.
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
		OrderBy("id ASC").
		ToSql()
	if err != nil {
		return nil, fmt.Errorf("build projects query: %w", err)
	}

	rows, err := repo.db.QueryContext(ctx, projectsQuery, projectsArgs...)
	if err != nil {
		return nil, fmt.Errorf("query projects: %w", err)
	}
	defer func() { _ = rows.Close() }()

	var projects []entity.Project
	projectIDs := make([]uint, 0)

	for rows.Next() {
		var p entity.Project
		if err := rows.Scan(
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
		); err != nil {
			return nil, fmt.Errorf("scan project: %w", err)
		}
		projects = append(projects, p)
		projectIDs = append(projectIDs, p.ID)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate projects: %w", err)
	}

	if len(projects) == 0 {
		return projects, nil
	}

	stagesQuery, stagesArgs, err := repo.psql.
		Select("project_id", "num", "script").
		From("stages").
		Where(squirrel.Eq{"project_id": projectIDs}).
		ToSql()
	if err != nil {
		return nil, fmt.Errorf("build stages query: %w", err)
	}

	rows, err = repo.db.QueryContext(ctx, stagesQuery, stagesArgs...)
	if err != nil {
		return nil, fmt.Errorf("query stages: %w", err)
	}
	defer func() { _ = rows.Close() }()

	stagesMap := make(map[uint][]entity.ProjectStage)
	for rows.Next() {
		var stage entity.ProjectStage
		if err := rows.Scan(&stage.ProjectID, &stage.Number, &stage.Script); err != nil {
			return nil, fmt.Errorf("scan stage: %w", err)
		}
		stagesMap[stage.ProjectID] = append(stagesMap[stage.ProjectID], stage)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate stages: %w", err)
	}

	for i := range projects {
		projects[i].Stages = stagesMap[projects[i].ID]
		slices.SortStableFunc(projects[i].Stages, func(a, b entity.ProjectStage) int {
			return a.Number - b.Number
		})
	}

	return projects, nil
}
