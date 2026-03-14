package database

import (
	"context"
	"fmt"

	"github.com/gbh007/easyjet/internal/entities"
	"github.com/samber/lo"
	"gorm.io/gorm"
)

func (repo Repo) ProjectRuns(ctx context.Context, id uint) ([]entities.ProjectRun, error) {
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

	result := make([]entities.ProjectRun, 0, len(runs))

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

		result = append(result, entities.ProjectRun{
			ID:        run.ID,
			ProjectID: run.ProjectID,
			Success:   run.Success,
			Stages: lo.Map(stages, func(r modelProjectStageRun, _ int) entities.ProjectRunStage {
				return entities.ProjectRunStage{
					StageNumber: r.StageNumber,
					Success:     r.Success,
					Log:         r.Log,
				}
			}),
		})
	}

	return result, nil
}

func (repo Repo) SetProjectRun(ctx context.Context, run entities.ProjectRun) (uint, error) {
	runRaw := modelProjectRun{
		Model: gorm.Model{
			ID: run.ID,
		},
		ProjectID: run.ProjectID,
		Success:   run.Success,
		Stages: lo.Map(run.Stages, func(r entities.ProjectRunStage, _ int) modelProjectStageRun {
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
