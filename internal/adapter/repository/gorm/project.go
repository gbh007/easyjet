package gorm

import (
	"context"
	"fmt"

	"github.com/gbh007/easyjet/internal/core/entity"
	"gorm.io/gorm"
)

func (repo Repo) Project(ctx context.Context, id uint) (entity.Project, error) {
	var p entity.Project

	res := repo.db.WithContext(ctx).Preload("Stages").First(&p, id)
	if res.Error != nil {
		return entity.Project{}, res.Error
	}

	return p, nil
}

func (repo Repo) SetProject(ctx context.Context, p entity.Project) (uint, error) {
	err := repo.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if p.ID > 0 {
			res := tx.Where(&entity.ProjectStage{
				ProjectID: p.ID,
			}).
				Delete(&entity.ProjectStage{})
			if res.Error != nil {
				return fmt.Errorf("delete stages: %w", res.Error)
			}
		}

		res := tx.Save(&p)
		if res.Error != nil {
			return fmt.Errorf("save project: %w", res.Error)
		}

		return nil
	})
	if err != nil {
		return 0, err
	}

	return p.ID, nil
}
