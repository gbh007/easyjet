package service

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/gbh007/easyjet/internal/core/entity"
	"github.com/samber/lo"
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

		var saveRunErr error

		_, saveRunErr = srv.db.SetProjectRun(ctx, run)
		returnedErr = errors.Join(returnedErr, saveRunErr)
	}()

	project, err := srv.db.Project(ctx, run.ProjectID)
	if err != nil {
		return fmt.Errorf("get project: %w", err)
	}

	defer func() {
		rotErr := srv.rotateProjectRuns(ctx, project)
		if rotErr != nil {
			srv.logger.Error("failed to rotate runs", "error", rotErr)
		}
	}()

	dir := project.Dir

	if dir == "" {
		dir = srv.fs.GetProjectDir(ctx, project.ID)
	}

	if project.HasGIT() {
		commits, err := srv.git.Diff(ctx, dir, "HEAD", srv.git.OriginName()+"/"+project.GitBranch)
		if err != nil {
			return fmt.Errorf("git get diff: %w", err)
		}

		if len(commits) > 0 {
			err = srv.db.SetProjectRunGitCommits(
				ctx,
				lo.Map(commits, func(c entity.Commit, i int) entity.ProjectRunGitCommits {
					return entity.ProjectRunGitCommits{
						RunID:   runID,
						Number:  i,
						Hash:    c.Hash,
						Subject: c.Subject,
					}
				}),
			)
			if err != nil {
				return fmt.Errorf("save git commits: %w", err)
			}
		}

		err = srv.git.Pull(ctx, dir, project.GitBranch)
		if err != nil {
			return fmt.Errorf("git pull origin: %w", err)
		}
	}

	for _, stage := range project.Stages {
		stageStart := time.Now()

		p, err := srv.fs.CreateSHScript(ctx, project.ID, stage.Number, stage.Script)
		if err != nil {
			return fmt.Errorf("create stage %d script: %w", stage.Number, err)
		}

		out, err := srv.ex.Exec(ctx, dir, p, project.WithRootEnv)
		if err != nil {
			err = fmt.Errorf("execute stage %d script: %w", stage.Number, err)
		}

		saveStageErr := srv.db.SetProjectRunStage(ctx, entity.ProjectRunStage{
			RunID:       runID,
			StageNumber: stage.Number,
			Success:     err == nil,
			Log:         out,
			Duration:    time.Since(stageStart),
		})
		if saveStageErr != nil {
			err = errors.Join(err, fmt.Errorf("save stage %d result: %w", stage.Number, saveStageErr))
		}

		if err != nil {
			return err
		}
	}

	if project.RestartAfter {
		srv.pubsub.PublishAppEvent(entity.AppEvent{
			Type:      entity.EventRequireAppRestart,
			ProjectID: project.ID,
		})
	}

	return nil
}
