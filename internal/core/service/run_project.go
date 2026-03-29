package service

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/gbh007/easyjet/internal/core/entity"
)

func (srv Service) RunProject(ctx context.Context, id uint) (uint, error) {
	// TODO: проверять что проект существует

	run := entity.ProjectRun{
		ProjectID: id,
		Pending:   true,
	}

	runID, err := srv.db.SetProjectRun(ctx, run)
	if err != nil {
		return 0, fmt.Errorf("create run: %w", err)
	}

	return runID, nil
}

func (srv Service) HandleRun(ctx context.Context, runID uint) (returnedErr error) {
	run, err := srv.db.ProjectRun(ctx, runID)
	if err != nil {
		return fmt.Errorf("get run: %w", err)
	}

	if run.Processing {
		return errors.New("run already processing")
	} else if !run.Pending {
		return errors.New("run already finished")
	}

	run.Pending = false
	run.Processing = true

	_, err = srv.db.SetProjectRun(ctx, run)
	if err != nil {
		return fmt.Errorf("update run status: %w", err)
	}

	run.Processing = false
	run.Success = true

	runStart := time.Now()

	defer func() {
		if returnedErr != nil {
			run.Success = false
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

	env, err := srv.CalculateEffectiveEnvVars(ctx, project, dir)
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
