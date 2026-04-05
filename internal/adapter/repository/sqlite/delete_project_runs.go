package sqlite

import (
	"context"
	"fmt"

	"github.com/Masterminds/squirrel"
)

func (repo Repo) DeleteProjectRuns(ctx context.Context, ids []int) error {
	if len(ids) == 0 {
		return nil
	}

	query, args, err := repo.psql.
		Delete("runs").
		Where(squirrel.Eq{"id": ids}).
		ToSql()
	if err != nil {
		return fmt.Errorf("build delete query: %w", err)
	}

	_, err = repo.db.ExecContext(ctx, query, args...)
	if err != nil {
		return fmt.Errorf("delete runs: %w", err)
	}

	return nil
}
