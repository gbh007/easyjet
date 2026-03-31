package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"strings"
	"time"

	"github.com/gbh007/easyjet/internal/core/entity"
)

func joinStrings(strs []string, sep string) string {
	if len(strs) == 0 {
		return ""
	}
	return strings.Join(strs, sep)
}

func (repo Repo) ProjectsWithRunInfo(ctx context.Context, filterType string) ([]entity.ProjectsWithRunInfo, error) {
	query := `
		SELECT
			p.id,
			p.name,
			p.cron_enabled,
			p.is_template,
			last_successful_run.last_successful_run_at,
			last_run.created_at as last_run_created_at,
			last_run.success as last_run_success,
			last_run.pending as last_run_pending,
			last_run.processing as last_run_processing,
			last_run.duration as last_run_duration
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
				r1.processing,
				r1.duration
			FROM runs r1
			INNER JOIN (
				SELECT project_id, MAX(created_at) as max_created_at
				FROM runs
				GROUP BY project_id
			) r2 ON r1.project_id = r2.project_id AND r1.created_at = r2.max_created_at
		) last_run ON p.id = last_run.project_id
	`

	args := make([]any, 0)
	whereClauses := make([]string, 0)

	switch filterType {
	case "project":
		whereClauses = append(whereClauses, "p.is_template = false")
	case "template":
		whereClauses = append(whereClauses, "p.is_template = true")
	}

	if len(whereClauses) > 0 {
		query += " WHERE " + joinStrings(whereClauses, " AND ")
	}

	query += " ORDER BY p.id ASC"

	rows, err := repo.pool.Query(ctx, query, args...)
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
		var lastRunDuration sql.NullInt64

		if err := rows.Scan(
			&item.ID,
			&item.Name,
			&item.CronEnabled,
			&item.IsTemplate,
			&lastSuccessfulRunAt,
			&lastRunCreatedAt,
			&lastRunSuccess,
			&lastRunPending,
			&lastRunProcessing,
			&lastRunDuration,
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
				Duration:   time.Duration(lastRunDuration.Int64) * time.Millisecond,
			}
		}

		items = append(items, item)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate project list items: %w", err)
	}

	return items, nil
}
