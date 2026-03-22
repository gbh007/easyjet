package postgres

import (
	"context"
	"fmt"

	"github.com/Masterminds/squirrel"
)

func (repo Repo) ProjectRunIDs(ctx context.Context, id uint) ([]uint, error) {
	query, args, err := repo.psql.
		Select("id").
		From("runs").
		Where(squirrel.Eq{"project_id": id}).
		OrderBy("id ASC").
		ToSql()
	if err != nil {
		return nil, fmt.Errorf("build query: %w", err)
	}

	rows, err := repo.pool.Query(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("query run ids: %w", err)
	}
	defer rows.Close()

	var ids []uint
	for rows.Next() {
		var runID uint
		if err := rows.Scan(&runID); err != nil {
			return nil, fmt.Errorf("scan run id: %w", err)
		}
		ids = append(ids, runID)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate run ids: %w", err)
	}

	return ids, nil
}
