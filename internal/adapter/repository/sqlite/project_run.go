package sqlite

import (
	"context"
	"fmt"
	"slices"
	"time"

	"github.com/Masterminds/squirrel"
	"github.com/gbh007/easyjet/internal/core/entity"
)

func (repo Repo) ProjectRun(ctx context.Context, id int) (entity.ProjectRun, error) {
	var run entity.ProjectRun

	runQuery, runArgs, err := repo.psql.
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
		Where(squirrel.Eq{"id": id}).
		ToSql()
	if err != nil {
		return entity.ProjectRun{}, fmt.Errorf("build run query: %w", err)
	}

	var durationMs int64

	err = repo.db.QueryRowContext(ctx, runQuery, runArgs...).Scan(
		&run.ID,
		&run.CreatedAt,
		&run.UpdatedAt,
		&run.ProjectID,
		&run.Status,
		&run.FailLog,
		&durationMs,
	)
	if err != nil {
		return entity.ProjectRun{}, fmt.Errorf("query run: %w", err)
	}

	run.Duration = time.Duration(durationMs) * time.Millisecond

	stagesQuery, stagesArgs, err := repo.psql.
		Select("run_id", "stage_num", "success", "log", "duration").
		From("run_stages").
		Where(squirrel.Eq{"run_id": id}).
		ToSql()
	if err != nil {
		return entity.ProjectRun{}, fmt.Errorf("build stages query: %w", err)
	}

	rows, err := repo.db.QueryContext(ctx, stagesQuery, stagesArgs...)
	if err != nil {
		return entity.ProjectRun{}, fmt.Errorf("query stages: %w", err)
	}
	defer func() { _ = rows.Close() }()

	for rows.Next() {
		var stage entity.ProjectRunStage
		var durationMs int64
		if err := rows.Scan(&stage.RunID, &stage.StageNumber, &stage.Success, &stage.Log, &durationMs); err != nil {
			return entity.ProjectRun{}, fmt.Errorf("scan stage: %w", err)
		}
		stage.Duration = time.Duration(durationMs) * time.Millisecond
		run.Stages = append(run.Stages, stage)
	}

	if err := rows.Err(); err != nil {
		return entity.ProjectRun{}, fmt.Errorf("iterate stages: %w", err)
	}

	slices.SortStableFunc(run.Stages, func(a, b entity.ProjectRunStage) int {
		return a.StageNumber - b.StageNumber
	})

	commitsQuery, commitsArgs, err := repo.psql.
		Select("run_id", "num", "hash", "subject").
		From("run_git_commits").
		Where(squirrel.Eq{"run_id": id}).
		ToSql()
	if err != nil {
		return entity.ProjectRun{}, fmt.Errorf("build commits query: %w", err)
	}

	rows, err = repo.db.QueryContext(ctx, commitsQuery, commitsArgs...)
	if err != nil {
		return entity.ProjectRun{}, fmt.Errorf("query commits: %w", err)
	}
	defer func() { _ = rows.Close() }()

	for rows.Next() {
		var commit entity.ProjectRunGitCommits
		if err := rows.Scan(&commit.RunID, &commit.Number, &commit.Hash, &commit.Subject); err != nil {
			return entity.ProjectRun{}, fmt.Errorf("scan commit: %w", err)
		}
		run.GitCommits = append(run.GitCommits, commit)
	}

	if err := rows.Err(); err != nil {
		return entity.ProjectRun{}, fmt.Errorf("iterate commits: %w", err)
	}

	slices.SortStableFunc(run.GitCommits, func(a, b entity.ProjectRunGitCommits) int {
		return a.Number - b.Number
	})

	return run, nil
}
