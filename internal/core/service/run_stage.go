package service

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/gbh007/easyjet/internal/core/entity"
)

func (srv Service) runStage(
	ctx context.Context,
	project entity.Project,
	runID int,
	dir string,
	stage entity.ProjectStage,
	env func(entity.ProjectStage) []string,
) (err error) {
	start := time.Now()

	defer func() {
		srv.pubsub.PublishEvent(entity.Event{
			Type:      entity.EventRunStageFinished,
			ProjectID: project.ID,
			RunID:     runID,
			Stage:     stage.Number,
			Err:       err,
			Duration:  time.Since(start),
		})
	}()

	p, err := srv.fs.CreateSHScript(ctx, project.ID, stage.Number, stage.Script)
	if err != nil {
		return fmt.Errorf("create stage %d script: %w", stage.Number, err)
	}

	out, err := srv.ex.Exec(ctx, dir, p, env(stage))
	if err != nil {
		err = fmt.Errorf("execute stage %d script: %w", stage.Number, err)
	}

	saveStageErr := srv.db.SetProjectRunStage(ctx, entity.ProjectRunStage{
		RunID:       runID,
		StageNumber: stage.Number,
		Success:     err == nil,
		Log:         out,
		Duration:    time.Since(start),
	})
	if saveStageErr != nil {
		err = errors.Join(err, fmt.Errorf("save stage %d result: %w", stage.Number, saveStageErr))
	}

	if err != nil {
		return err
	}

	return nil
}
