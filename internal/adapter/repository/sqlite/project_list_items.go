package sqlite

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/Masterminds/squirrel"
	"github.com/gbh007/easyjet/internal/core/entity"
)

func (repo Repo) ProjectsWithRunInfo(ctx context.Context, filterType entity.ProjectFilterType) ([]entity.ProjectsWithRunInfo, error) {
	builder := repo.psql.
		Select(
			"p.id",
			"p.name",
			"p.cron_enabled",
			"p.is_template",
			"last_successful_run.last_successful_run_at",
			"last_run.created_at",
			"last_run.status",
			"last_run.duration",
		).
		From("projects p").
		JoinClause(
			squirrel.
				Select(
					"project_id",
					"MAX(created_at) as last_successful_run_at",
				).
				From("runs").
				Where(squirrel.Eq{
					"status": entity.StatusSuccess,
				}).
				GroupBy("project_id").
				Prefix("LEFT JOIN (").
				Suffix(") last_successful_run ON p.id = last_successful_run.project_id"),
		).
		JoinClause(
			squirrel.
				Select(
					"r1.project_id",
					"r1.created_at",
					"r1.status",
					"r1.duration",
				).
				From("runs r1").
				JoinClause(
					squirrel.
						Select(
							"project_id",
							"MAX(id) as max_id",
						).
						From("runs").
						GroupBy("project_id").
						Prefix("INNER JOIN (").
						Suffix(") r2 ON r1.project_id = r2.project_id AND r1.id = r2.max_id"),
				).
				Prefix("LEFT JOIN (").
				Suffix(") last_run ON p.id = last_run.project_id"),
		).
		OrderBy("p.id ASC")

	switch filterType {
	case entity.ProjectFilterTypeAll:
		// No filter, return all projects
	case entity.ProjectFilterTypeProject:
		builder = builder.Where(squirrel.Eq{
			"p.is_template": false,
		})
	case entity.ProjectFilterTypeTemplate:
		builder = builder.Where(squirrel.Eq{
			"p.is_template": true,
		})
	}

	query, args, err := builder.ToSql()
	if err != nil {
		return nil, fmt.Errorf("build query: %w", err)
	}

	rows, err := repo.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("query project list items: %w", err)
	}
	defer func() { _ = rows.Close() }()

	var items []entity.ProjectsWithRunInfo
	for rows.Next() {
		var item entity.ProjectsWithRunInfo
		var lastSuccessfulRunAt sql.NullString
		var lastRunCreatedAt sql.NullString
		var lastRunStatus sql.NullString
		var lastRunDuration sql.NullInt64

		if err := rows.Scan(
			&item.ID,
			&item.Name,
			&item.CronEnabled,
			&item.IsTemplate,
			&lastSuccessfulRunAt,
			&lastRunCreatedAt,
			&lastRunStatus,
			&lastRunDuration,
		); err != nil {
			return nil, fmt.Errorf("scan project list item: %w", err)
		}

		if lastSuccessfulRunAt.Valid {
			t, err := time.Parse("2006-01-02 15:04:05", lastSuccessfulRunAt.String)
			if err != nil {
				return nil, fmt.Errorf("parse last success run time: %w", err)
			}

			item.LastSuccessfulRunAt = &t
		}

		if lastRunCreatedAt.Valid {
			t, err := time.Parse(time.RFC3339, lastRunCreatedAt.String)
			if err != nil {
				return nil, fmt.Errorf("parse last run time: %w", err)
			}

			item.LastRun = &entity.ProjectLastRun{
				CreatedAt: t,
				Status:    lastRunStatus.String,
				Duration:  time.Duration(lastRunDuration.Int64) * time.Millisecond,
			}
		}

		items = append(items, item)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate project list items: %w", err)
	}

	return items, nil
}
