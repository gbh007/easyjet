package gorm

import (
	"context"
	"slices"

	"github.com/gbh007/easyjet/internal/core/entity"
)

func (repo Repo) ProjectRuns(ctx context.Context, id uint) ([]entity.ProjectRun, error) {
	var runs []entity.ProjectRun

	res := repo.db.WithContext(ctx).
		Where(&entity.ProjectRun{
			ProjectID: id,
		}).
		Preload("Stages").
		Preload("GitCommits").
		Find(&runs)
	if res.Error != nil {
		return nil, res.Error
	}

	for i := range runs {
		slices.SortStableFunc(runs[i].Stages, func(a, b entity.ProjectRunStage) int {
			return a.StageNumber - b.StageNumber
		})
		slices.SortStableFunc(runs[i].GitCommits, func(a, b entity.ProjectRunGitCommits) int {
			return a.Number - b.Number
		})
	}

	return runs, nil
}

func (repo Repo) ProjectRun(ctx context.Context, id uint) (entity.ProjectRun, error) {
	var run entity.ProjectRun

	res := repo.db.WithContext(ctx).
		Preload("Stages").
		Preload("GitCommits").
		First(&run, id)
	if res.Error != nil {
		return entity.ProjectRun{}, res.Error
	}

	slices.SortStableFunc(run.Stages, func(a, b entity.ProjectRunStage) int {
		return a.StageNumber - b.StageNumber
	})
	slices.SortStableFunc(run.GitCommits, func(a, b entity.ProjectRunGitCommits) int {
		return a.Number - b.Number
	})

	return run, nil
}

func (repo Repo) SetProjectRun(ctx context.Context, run entity.ProjectRun) (uint, error) {
	res := repo.db.Save(&run)
	if res.Error != nil {
		return 0, res.Error
	}

	return run.ID, nil
}

func (repo Repo) SetProjectRunStage(ctx context.Context, rs entity.ProjectRunStage) error {
	res := repo.db.Create(&rs)
	if res.Error != nil {
		return res.Error
	}

	return nil
}

func (repo Repo) SetProjectRunGitCommits(ctx context.Context, commits []entity.ProjectRunGitCommits) error {
	res := repo.db.Save(&commits)
	if res.Error != nil {
		return res.Error
	}

	return nil
}

func (repo Repo) PendingProjectRuns(ctx context.Context) ([]uint, error) {
	var runIDs []uint

	res := repo.db.WithContext(ctx).
		Model(&entity.ProjectRun{}).
		Where(&entity.ProjectRun{
			Pending: true,
		}).
		Select("id").
		Pluck("id", &runIDs)
	if res.Error != nil {
		return nil, res.Error
	}

	return runIDs, nil
}

func (repo Repo) DeleteProjectRuns(ctx context.Context, ids []uint) error {
	res := repo.db.WithContext(ctx).Delete(&[]entity.ProjectRun{}, ids)
	if res.Error != nil {
		return res.Error
	}

	return nil
}

func (repo Repo) ProjectRunIDs(ctx context.Context, id uint) ([]uint, error) {
	var runIDs []uint

	res := repo.db.WithContext(ctx).
		Model(&entity.ProjectRun{}).
		Where(&entity.ProjectRun{
			ProjectID: id,
		}).
		Select("id").
		Pluck("id", &runIDs)
	if res.Error != nil {
		return nil, res.Error
	}

	return runIDs, nil
}
