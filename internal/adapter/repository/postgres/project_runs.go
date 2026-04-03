package postgres

import (
	"context"
	"fmt"
	"time"

	"github.com/Masterminds/squirrel"
	"github.com/gbh007/easyjet/internal/core/entity"
)

func (repo Repo) ProjectRuns(ctx context.Context, id uint) ([]entity.ProjectRun, error) {
	runsQuery, runsArgs, err := repo.psql.
		Select(
			"id",
			"created_at",
			"updated_at",
			"project_id",
			"status",
			"fail_log",
			"duration",
		).
		From("runs").
		Where(squirrel.Eq{"project_id": id}).
		OrderBy("id ASC").
		ToSql()
	if err != nil {
		return nil, fmt.Errorf("build runs query: %w", err)
	}

	rows, err := repo.pool.Query(ctx, runsQuery, runsArgs...)
	if err != nil {
		return nil, fmt.Errorf("query runs: %w", err)
	}
	defer rows.Close()

	var runs []entity.ProjectRun

	for rows.Next() {
		var run entity.ProjectRun
		var durationMs int64
		if err := rows.Scan(
			&run.ID,
			&run.CreatedAt,
			&run.UpdatedAt,
			&run.ProjectID,
			&run.Status,
			&run.FailLog,
			&durationMs,
		); err != nil {
			return nil, fmt.Errorf("scan run: %w", err)
		}
		run.Duration = time.Duration(durationMs) * time.Millisecond
		runs = append(runs, run)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate runs: %w", err)
	}

	return runs, nil
}
