package gorm

import (
	"context"
	"fmt"
	"slices"

	"github.com/gbh007/easyjet/internal/core/entity"
	"gorm.io/gorm"
)

func (repo Repo) Project(ctx context.Context, id uint) (entity.Project, error) {
	var p entity.Project

	res := repo.db.WithContext(ctx).Preload("Stages").First(&p, id)
	if res.Error != nil {
		return entity.Project{}, res.Error
	}

	slices.SortStableFunc(p.Stages, func(a, b entity.ProjectStage) int {
		return a.Number - b.Number
	})

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

func (repo Repo) Projects(ctx context.Context) ([]entity.Project, error) {
	var projects []entity.Project

	res := repo.db.WithContext(ctx).Find(&projects)
	if res.Error != nil {
		return nil, res.Error
	}

	slices.SortStableFunc(projects, func(a, b entity.Project) int {
		return int(a.ID) - int(b.ID)
	})

	return projects, nil
}
