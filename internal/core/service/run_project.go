package service

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/gbh007/easyjet/internal/core/entity"
)

func (srv Service) RunProject(ctx context.Context, id int) (int, error) {
	project, err := srv.db.Project(ctx, id)
	if err != nil {
		return 0, fmt.Errorf("get project: %w", err)
	}

	if project.IsTemplate {
		return 0, errors.New("cannot run a template project")
	}

	run := entity.ProjectRun{
		ProjectID: id,
		Status:    entity.StatusPending,
	}

	runID, err := srv.db.SetProjectRun(ctx, run)
	if err != nil {
		return 0, fmt.Errorf("create run: %w", err)
	}

	return runID, nil
}

func (srv Service) HandleRun(ctx context.Context, runID int) (returnedErr error) {
	run, err := srv.db.ProjectRun(ctx, runID)
	if err != nil {
		return fmt.Errorf("get run: %w", err)
	}

	if run.Status == entity.StatusProcessing {
		return errors.New("run already processing")
	} else if run.Status != entity.StatusPending {
		return errors.New("run already finished")
	}

	run.Status = entity.StatusProcessing

	_, err = srv.db.SetProjectRun(ctx, run)
	if err != nil {
		return fmt.Errorf("update run status: %w", err)
	}

	run.Status = entity.StatusSuccess

	runStart := time.Now()

	defer func() {
		if returnedErr != nil {
			run.Status = entity.StatusFailed
			run.FailLog = returnedErr.Error()
		}

		run.Duration = time.Since(runStart)

		srv.pubsub.PublishEvent(entity.Event{
			Type:      entity.EventRunFinished,
			ProjectID: run.ProjectID,
			RunID:     runID,
			Err:       returnedErr,
			Duration:  run.Duration,
		})

		var saveRunErr error

		_, saveRunErr = srv.db.SetProjectRun(ctx, run)
		returnedErr = errors.Join(returnedErr, saveRunErr)
	}()

	project, err := srv.db.Project(ctx, run.ProjectID)
	if err != nil {
		return fmt.Errorf("get project: %w", err)
	}

	defer func() {
		rotErr := srv.rotateProjectRuns(ctx, project, runID)
		if rotErr != nil {
			srv.logger.Error("failed to rotate runs", "error", rotErr)
		}
	}()

	dir := project.Dir

	if dir == "" {
		dir = srv.fs.GetProjectDir(ctx, project.ID)
	}

	if project.HasGIT() {
		err := srv.runGitPull(ctx, project, runID, dir)
		if err != nil {
			return err
		}
	}

	env, err := srv.CalculateEffectiveEnvVars(ctx, project, run, dir)
	if err != nil {
		return fmt.Errorf("calculate env vars: %w", err)
	}

	for _, stage := range project.Stages {
		err = srv.runStage(ctx, project, runID, dir, stage, env)
		if err != nil {
			return err
		}
	}

	if project.RestartAfter {
		srv.pubsub.PublishEvent(entity.Event{
			Type:      entity.EventRequireAppRestart,
			ProjectID: project.ID,
			RunID:     runID,
		})
	}

	return nil
}
