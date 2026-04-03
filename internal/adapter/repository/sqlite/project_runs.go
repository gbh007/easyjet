package sqlite

import (
	"context"
	"fmt"
	"slices"
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

	rows, err := repo.db.QueryContext(ctx, runsQuery, runsArgs...)
	if err != nil {
		return nil, fmt.Errorf("query runs: %w", err)
	}
	defer func() { _ = rows.Close() }()

	var runs []entity.ProjectRun
	runIDs := make([]uint, 0)

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
		runIDs = append(runIDs, run.ID)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate runs: %w", err)
	}

	if len(runs) == 0 {
		return runs, nil
	}

	stagesQuery, stagesArgs, err := repo.psql.
		Select("run_id", "stage_num", "success", "log", "duration").
		From("run_stages").
		Where(squirrel.Eq{"run_id": runIDs}).
		ToSql()
	if err != nil {
		return nil, fmt.Errorf("build stages query: %w", err)
	}

	rows, err = repo.db.QueryContext(ctx, stagesQuery, stagesArgs...)
	if err != nil {
		return nil, fmt.Errorf("query stages: %w", err)
	}
	defer func() { _ = rows.Close() }()

	stagesMap := make(map[uint][]entity.ProjectRunStage)
	for rows.Next() {
		var stage entity.ProjectRunStage
		var durationMs int64
		if err := rows.Scan(&stage.RunID, &stage.StageNumber, &stage.Success, &stage.Log, &durationMs); err != nil {
			return nil, fmt.Errorf("scan stage: %w", err)
		}
		stage.Duration = time.Duration(durationMs) * time.Millisecond
		stagesMap[stage.RunID] = append(stagesMap[stage.RunID], stage)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate stages: %w", err)
	}

	commitsQuery, commitsArgs, err := repo.psql.
		Select("run_id", "num", "hash", "subject").
		From("run_git_commits").
		Where(squirrel.Eq{"run_id": runIDs}).
		ToSql()
	if err != nil {
		return nil, fmt.Errorf("build commits query: %w", err)
	}

	rows, err = repo.db.QueryContext(ctx, commitsQuery, commitsArgs...)
	if err != nil {
		return nil, fmt.Errorf("query commits: %w", err)
	}
	defer func() { _ = rows.Close() }()

	commitsMap := make(map[uint][]entity.ProjectRunGitCommits)
	for rows.Next() {
		var commit entity.ProjectRunGitCommits
		if err := rows.Scan(&commit.RunID, &commit.Number, &commit.Hash, &commit.Subject); err != nil {
			return nil, fmt.Errorf("scan commit: %w", err)
		}
		commitsMap[commit.RunID] = append(commitsMap[commit.RunID], commit)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate commits: %w", err)
	}

	for i := range runs {
		if stages, ok := stagesMap[runs[i].ID]; ok {
			runs[i].Stages = stages
			slices.SortStableFunc(runs[i].Stages, func(a, b entity.ProjectRunStage) int {
				return a.StageNumber - b.StageNumber
			})
		}

		if commits, ok := commitsMap[runs[i].ID]; ok {
			runs[i].GitCommits = commits
			slices.SortStableFunc(runs[i].GitCommits, func(a, b entity.ProjectRunGitCommits) int {
				return a.Number - b.Number
			})
		}
	}

	return runs, nil
}
