package gorm

import (
	"context"
	"fmt"
	"slices"

	"github.com/gbh007/easyjet/internal/core/entity"
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
	res := repo.db.WithContext(ctx).Save(&p)
	if res.Error != nil {
		return 0, fmt.Errorf("save stage: %w", res.Error)
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
