package service

import (
	"context"
	"fmt"

	"github.com/gbh007/easyjet/internal/adapter/scheduler"
	"github.com/gbh007/easyjet/internal/core/entity"
	"github.com/samber/lo"
)

func (srv Service) CreateProject(ctx context.Context, p entity.Project) (uint, error) {
	// Validate cron expression if provided
	if err := validateCronExpression(p.CronSchedule); err != nil {
		return 0, err
	}

	p.Stages = lo.FilterMap(p.Stages, func(s entity.ProjectStage, i int) (entity.ProjectStage, bool) {
		s.Number = i + 1

		return s, s.Script != ""
	})

	id, err := srv.db.SetProject(ctx, p)
	if err != nil {
		return 0, fmt.Errorf("create project: %w", err)
	}

	// Publish scheduler event for cron scheduling
	if srv.publisher != nil {
		err = srv.publisher.Publish(scheduler.SchedulerEvent{
			Type:      scheduler.EventCreated,
			ProjectID: id,
			Schedule:  p.CronSchedule,
			Enabled:   p.CronEnabled,
		})
		if err != nil {
			srv.logger.Error("failed to publish scheduler event", "error", err, "project_id", id)
		}
	}

	dir := p.Dir

	if p.Dir == "" {
		dir, err = srv.fs.CreateProjectDir(ctx, id)
		if err != nil {
			return 0, fmt.Errorf("create project dir: %w", err)
		}
	}

	if !p.HasGIT() {
		return id, nil
	}

	err = srv.git.Init(ctx, dir, p.GitBranch, p.GitURL)
	if err != nil {
		return 0, fmt.Errorf("init git: %w", err)
	}

	err = srv.git.Pull(ctx, dir, p.GitBranch)
	if err != nil {
		return 0, fmt.Errorf("pull git: %w", err)
	}

	return id, nil
}
