package sqlite

import (
	"context"
	"fmt"

	"github.com/gbh007/easyjet/internal/core/entity"
)

func (repo Repo) SetProjectRunGitCommits(ctx context.Context, commits []entity.ProjectRunGitCommits) error {
	for _, commit := range commits {
		insertQuery, insertArgs, err := repo.psql.
			Insert("run_git_commits").
			SetMap(map[string]any{
				"run_id":  commit.RunID,
				"num":     commit.Number,
				"hash":    commit.Hash,
				"subject": commit.Subject,
			}).
			Suffix("ON CONFLICT (run_id, num) DO UPDATE SET hash = EXCLUDED.hash, subject = EXCLUDED.subject").
			ToSql()
		if err != nil {
			return fmt.Errorf("build insert query: %w", err)
		}

		_, err = repo.db.ExecContext(ctx, insertQuery, insertArgs...)
		if err != nil {
			return fmt.Errorf("insert commit: %w", err)
		}
	}

	return nil
}
