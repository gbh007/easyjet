package service

import (
	"context"
	"fmt"

	"github.com/gbh007/easyjet/internal/core/entity"
	"github.com/samber/lo"
)

func (srv Service) CreateProject(ctx context.Context, p entity.Project) (uint, error) {
	p.Stages = lo.FilterMap(p.Stages, func(s entity.ProjectStage, i int) (entity.ProjectStage, bool) {
		s.Number = i + 1

		return s, s.Script != ""
	})

	id, err := srv.db.SetProject(ctx, p)
	if err != nil {
		return 0, fmt.Errorf("create project: %w", err)
	}

	if !p.HasGIT() {
		return id, nil
	}

	dir := p.Dir

	if p.Dir == "" {
		dir, err = srv.fs.CreateProjectDir(ctx, id)
		if err != nil {
			return 0, fmt.Errorf("create project dir: %w", err)
		}
	}

	err = srv.git.Init(ctx, dir, p.GitBranch, p.GitURL)
	if err != nil {
		return 0, fmt.Errorf("init git: %w", err)
	}

	return id, nil
}
