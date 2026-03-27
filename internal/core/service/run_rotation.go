package service

import (
	"context"
	"fmt"
	"slices"
	"time"

	"github.com/gbh007/easyjet/internal/core/entity"
)

func (srv Service) rotateProjectRuns(ctx context.Context, project entity.Project, runID uint) (err error) {
	start := time.Now()

	defer func() {
		srv.pubsub.PublishEvent(entity.Event{
			Type:      entity.EventRunRotateFinished,
			ProjectID: project.ID,
			RunID:     runID,
			Err:       err,
			Duration:  time.Since(start),
		})
	}()

	retentionCount := project.RetentionCount
	if retentionCount <= 0 {
		return nil
	}

	runs, err := srv.db.ProjectRunIDs(ctx, project.ID)
	if err != nil {
		return fmt.Errorf("get project runs: %w", err)
	}

	if len(runs) <= retentionCount {
		return nil
	}

	slices.Sort(runs)
	toRemove := runs[:len(runs)-retentionCount]

	err = srv.db.DeleteProjectRuns(ctx, toRemove)
	if err != nil {
		return fmt.Errorf("delete runs: %w", err)
	}

	return nil
}
