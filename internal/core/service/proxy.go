package service

import (
	"context"

	"github.com/gbh007/easyjet/internal/core/entity"
)

func (srv Service) Project(ctx context.Context, id uint) (entity.Project, error) {
	return srv.db.Project(ctx, id)
}

func (srv Service) ProjectRun(ctx context.Context, id uint) (entity.ProjectRun, error) {
	return srv.db.ProjectRun(ctx, id)
}

func (srv Service) ProjectRuns(ctx context.Context, id uint) ([]entity.ProjectRun, error) {
	return srv.db.ProjectRuns(ctx, id)
}
