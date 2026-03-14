package gorm

import (
	"context"
	"fmt"

	"github.com/gbh007/easyjet/internal/core/entity"
	"github.com/samber/lo"
	"gorm.io/gorm"
)

func (repo Repo) ProjectRuns(ctx context.Context, id uint) ([]entity.ProjectRun, error) {
	var (
		runs   []modelProjectRun
		stages []modelProjectStageRun
	)

	res := repo.db.WithContext(ctx).
		Where(&modelProjectRun{
			ProjectID: id,
		}).
		Find(&runs)
	if res.Error != nil {
		return nil, fmt.Errorf("get run: %w", res.Error)
	}

	result := make([]entity.ProjectRun, 0, len(runs))

	for _, run := range runs {
		res = repo.db.WithContext(ctx).
			Where(&modelProjectStageRun{
				RunID: run.ID,
			}).
			Order("stage_num").
			Find(&stages)
		if res.Error != nil {
			return nil, fmt.Errorf("get stages: %w", res.Error)
		}

		result = append(result, entity.ProjectRun{
			ID:        run.ID,
			ProjectID: run.ProjectID,
			Success:   run.Success,
			Stages: lo.Map(stages, func(r modelProjectStageRun, _ int) entity.ProjectRunStage {
				return entity.ProjectRunStage{
					StageNumber: r.StageNumber,
					Success:     r.Success,
					Log:         r.Log,
				}
			}),
		})
	}

	return result, nil
}

func (repo Repo) SetProjectRun(ctx context.Context, run entity.ProjectRun) (uint, error) {
	runRaw := modelProjectRun{
		Model: gorm.Model{
			ID: run.ID,
		},
		ProjectID: run.ProjectID,
		Success:   run.Success,
		Stages: lo.Map(run.Stages, func(r entity.ProjectRunStage, _ int) modelProjectStageRun {
			return modelProjectStageRun{
				StageNumber: r.StageNumber,
				Success:     r.Success,
				Log:         r.Log,
			}
		}),
	}

	res := repo.db.Save(&runRaw)
	if res.Error != nil {
		return 0, res.Error
	}

	return runRaw.ID, nil
}
