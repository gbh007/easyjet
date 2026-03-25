package service

import (
	"context"
	"fmt"

	"github.com/gbh007/easyjet/internal/core/entity"
	"github.com/samber/lo"
)

func (srv Service) UpdateProject(ctx context.Context, p entity.Project) error {
	p.Stages = lo.FilterMap(p.Stages, func(s entity.ProjectStage, i int) (entity.ProjectStage, bool) {
		s.Number = i + 1

		return s, s.Script != ""
	})

	_, err := srv.db.SetProject(ctx, p)
	if err != nil {
		return fmt.Errorf("create project: %w", err)
	}

	srv.pubsub.PublishEvent(entity.Event{
		Type:      entity.EventProjectUpdated,
		ProjectID: p.ID,
	})

	// TODO: обновлять и настройки гита + создавать папку если нужно

	return nil
}
