package sqlite

import (
	"context"
	"fmt"

	"github.com/Masterminds/squirrel"
	"github.com/gbh007/easyjet/internal/core/entity"
)

func (repo Repo) PendingProjectRuns(ctx context.Context) ([]int, error) {
	query, args, err := repo.psql.
		Select("id").
		From("runs").
		Where(squirrel.Eq{"status": entity.StatusPending}).
		OrderBy("id ASC").
		ToSql()
	if err != nil {
		return nil, fmt.Errorf("build query: %w", err)
	}

	rows, err := repo.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("query pending runs: %w", err)
	}
	defer func() { _ = rows.Close() }()

	var ids []int
	for rows.Next() {
		var runID int
		if err := rows.Scan(&runID); err != nil {
			return nil, fmt.Errorf("scan run id: %w", err)
		}
		ids = append(ids, runID)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate pending runs: %w", err)
	}

	return ids, nil
}
