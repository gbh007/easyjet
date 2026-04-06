package sqlite

import (
	"context"
	"fmt"

	"github.com/Masterminds/squirrel"
)

func (repo Repo) DeleteProject(ctx context.Context, id int) error {
	query, args, err := repo.psql.
		Delete("projects").
		Where(squirrel.Eq{"id": id}).
		ToSql()
	if err != nil {
		return fmt.Errorf("build delete project query: %w", err)
	}

	_, err = repo.db.ExecContext(ctx, query, args...)
	if err != nil {
		return fmt.Errorf("delete project: %w", err)
	}

	return nil
}
