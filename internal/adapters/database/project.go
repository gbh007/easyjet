package database

import (
	"context"
	"fmt"

	"github.com/gbh007/easyjet/internal/entities"
	"github.com/samber/lo"
	"gorm.io/gorm"
)

func (repo Repo) Project(ctx context.Context, id uint) (entities.Project, error) {
	var (
		p      modelProject
		stages []*modelProjectStage
	)

	res := repo.db.WithContext(ctx).First(&p, id)
	if res.Error != nil {
		return entities.Project{}, fmt.Errorf("get project: %w", res.Error)
	}

	res = repo.db.WithContext(ctx).
		Where(&modelProjectStage{
			ProjectID: id,
		}).
		Order("num").
		Find(&stages)
	if res.Error != nil {
		return entities.Project{}, fmt.Errorf("get stages: %w", res.Error)
	}

	return entities.Project{
		ID:     p.ID,
		Dir:    p.Dir,
		GitURL: p.GitURL,
		Name:   p.Name,
		Stages: lo.Map(stages, func(s *modelProjectStage, _ int) string {
			return s.Script
		}),
	}, nil
}

func (repo Repo) SetProject(ctx context.Context, pr entities.Project) (uint, error) {
	p := modelProject{
		Model: gorm.Model{
			ID: pr.ID,
		},
		Dir:    pr.Dir,
		GitURL: pr.GitURL,
		Name:   pr.Name,
		Stages: lo.Map(pr.Stages, func(s string, num int) modelProjectStage {
			return modelProjectStage{
				Script: s,
				Number: num,
			}
		}),
	}

	err := repo.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if pr.ID > 0 {
			res := tx.Where(&modelProjectStage{
				ProjectID: pr.ID,
			}).
				Delete(&modelProjectStage{})
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
