package postgres

import (
	"context"
	"fmt"

	"github.com/Masterminds/squirrel"
	"github.com/gbh007/easyjet/internal/core/entity"
	"github.com/jackc/pgx/v5"
)

func (repo Repo) SetProjectRun(ctx context.Context, run entity.ProjectRun) (uint, error) {
	tx, err := repo.pool.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		return 0, fmt.Errorf("begin transaction: %w", err)
	}
	defer func() {
		_ = tx.Rollback(ctx)
	}()

	runID := run.ID

	if run.ID > 0 {
		updateQuery, updateArgs, err := repo.psql.
			Update("runs").
			SetMap(map[string]any{
				"updated_at": run.UpdatedAt,
				"success":    run.Success,
				"pending":    run.Pending,
				"processing": run.Processing,
				"fail_log":   run.FailLog,
			}).
			Where(squirrel.Eq{"id": run.ID}).
			Suffix("RETURNING id").
			ToSql()
		if err != nil {
			return 0, fmt.Errorf("build update query: %w", err)
		}

		err = tx.QueryRow(ctx, updateQuery, updateArgs...).Scan(&runID)
		if err != nil {
			return 0, fmt.Errorf("update run: %w", err)
		}
	} else {
		insertQuery, insertArgs, err := repo.psql.
			Insert("runs").
			SetMap(map[string]any{
				"project_id": run.ProjectID,
				"success":    run.Success,
				"pending":    run.Pending,
				"processing": run.Processing,
				"fail_log":   run.FailLog,
			}).
			Suffix("RETURNING id").
			ToSql()
		if err != nil {
			return 0, fmt.Errorf("build insert query: %w", err)
		}

		err = tx.QueryRow(ctx, insertQuery, insertArgs...).Scan(&runID)
		if err != nil {
			return 0, fmt.Errorf("insert run: %w", err)
		}
	}

	for _, stage := range run.Stages {
		insertStageQuery, insertStageArgs, err := repo.psql.
			Insert("run_stages").
			SetMap(map[string]any{
				"run_id":    runID,
				"stage_num": stage.StageNumber,
				"success":   stage.Success,
				"log":       stage.Log,
			}).
			Suffix("ON CONFLICT (run_id, stage_num) DO UPDATE SET success = EXCLUDED.success, log = EXCLUDED.log").
			ToSql()
		if err != nil {
			return 0, fmt.Errorf("build insert stage query: %w", err)
		}

		_, err = tx.Exec(ctx, insertStageQuery, insertStageArgs...)
		if err != nil {
			return 0, fmt.Errorf("insert stage: %w", err)
		}
	}

	for _, commit := range run.GitCommits {
		insertCommitQuery, insertCommitArgs, err := repo.psql.
			Insert("run_git_commits").
			SetMap(map[string]any{
				"run_id":  runID,
				"num":     commit.Number,
				"hash":    commit.Hash,
				"subject": commit.Subject,
			}).
			Suffix(`ON CONFLICT (run_id, num) DO UPDATE SET hash = EXCLUDED.hash, subject = EXCLUDED.subject`).
			ToSql()
		if err != nil {
			return 0, fmt.Errorf("build insert commit query: %w", err)
		}

		_, err = tx.Exec(ctx, insertCommitQuery, insertCommitArgs...)
		if err != nil {
			return 0, fmt.Errorf("insert commit: %w", err)
		}
	}

	if err := tx.Commit(ctx); err != nil {
		return 0, fmt.Errorf("commit transaction: %w", err)
	}

	return runID, nil
}
