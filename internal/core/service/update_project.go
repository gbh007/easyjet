package service

import (
	"context"
	"fmt"

	"github.com/gbh007/easyjet/internal/adapter/scheduler"
	"github.com/gbh007/easyjet/internal/core/entity"
	"github.com/samber/lo"
)

func (srv Service) UpdateProject(ctx context.Context, p entity.Project) error {
	// Validate cron expression if provided
	if err := validateCronExpression(p.CronSchedule); err != nil {
		return err
	}

	p.Stages = lo.FilterMap(p.Stages, func(s entity.ProjectStage, i int) (entity.ProjectStage, bool) {
		s.Number = i + 1

		return s, s.Script != ""
	})

	_, err := srv.db.SetProject(ctx, p)
	if err != nil {
		return fmt.Errorf("create project: %w", err)
	}

	// Publish scheduler event for cron scheduling
	// FIXME(ai-shit): убрать бесполезные проверки.
	if srv.publisher != nil {
		err = srv.publisher.Publish(scheduler.SchedulerEvent{
			Type:      scheduler.EventUpdated,
			ProjectID: p.ID,
			Schedule:  p.CronSchedule,
			Enabled:   p.CronEnabled,
		})
		if err != nil {
			srv.logger.Error("failed to publish scheduler event", "error", err, "project_id", p.ID)
		}
	}

	// TODO: обновлять и настройки гита + создавать папку если нужно

	return nil
}
