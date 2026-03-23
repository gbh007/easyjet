package postgres

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/gbh007/easyjet/internal/core/entity"
)

func (repo Repo) ProjectsWithRunInfo(ctx context.Context) ([]entity.ProjectsWithRunInfo, error) {
	query := `
		SELECT 
			p.id,
			p.name,
			p.cron_enabled,
			last_successful_run.last_successful_run_at,
			last_run.created_at as last_run_created_at,
			last_run.success as last_run_success,
			last_run.pending as last_run_pending,
			last_run.processing as last_run_processing
		FROM projects p
		LEFT JOIN (
			SELECT 
				project_id,
				MAX(created_at) as last_successful_run_at
			FROM runs
			WHERE success = true
			GROUP BY project_id
		) last_successful_run ON p.id = last_successful_run.project_id
		LEFT JOIN (
			SELECT 
				r1.project_id,
				r1.created_at,
				r1.success,
				r1.pending,
				r1.processing
			FROM runs r1
			INNER JOIN (
				SELECT project_id, MAX(created_at) as max_created_at
				FROM runs
				GROUP BY project_id
			) r2 ON r1.project_id = r2.project_id AND r1.created_at = r2.max_created_at
		) last_run ON p.id = last_run.project_id
		ORDER BY p.id ASC
	`

	rows, err := repo.pool.Query(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("query project list items: %w", err)
	}
	defer rows.Close()

	var items []entity.ProjectsWithRunInfo
	for rows.Next() {
		var item entity.ProjectsWithRunInfo
		var lastSuccessfulRunAt sql.NullTime
		var lastRunCreatedAt sql.NullTime
		var lastRunSuccess sql.NullBool
		var lastRunPending sql.NullBool
		var lastRunProcessing sql.NullBool

		if err := rows.Scan(
			&item.ID,
			&item.Name,
			&item.CronEnabled,
			&lastSuccessfulRunAt,
			&lastRunCreatedAt,
			&lastRunSuccess,
			&lastRunPending,
			&lastRunProcessing,
		); err != nil {
			return nil, fmt.Errorf("scan project list item: %w", err)
		}

		if lastSuccessfulRunAt.Valid {
			item.LastSuccessfulRunAt = &lastSuccessfulRunAt.Time
		}

		if lastRunCreatedAt.Valid {
			item.LastRun = &entity.ProjectLastRun{
				CreatedAt:  lastRunCreatedAt.Time,
				Success:    lastRunSuccess.Bool,
				Pending:    lastRunPending.Bool,
				Processing: lastRunProcessing.Bool,
			}
		}

		items = append(items, item)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate project list items: %w", err)
	}

	return items, nil
}
