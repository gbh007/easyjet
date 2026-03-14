package gorm

import (
	"context"

	"github.com/gbh007/easyjet/internal/core/entity"
)

func (repo Repo) ProjectRuns(ctx context.Context, id uint) ([]entity.ProjectRun, error) {
	var runs []entity.ProjectRun

	res := repo.db.WithContext(ctx).
		Where(&entity.ProjectRun{
			ProjectID: id,
		}).
		Preload("Stages").
		Find(&runs)
	if res.Error != nil {
		return nil, res.Error
	}

	return runs, nil
}

func (repo Repo) SetProjectRun(ctx context.Context, run entity.ProjectRun) (uint, error) {
	res := repo.db.Save(&run)
	if res.Error != nil {
		return 0, res.Error
	}

	return run.ID, nil
}
