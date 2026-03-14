package service

import (
	"context"
	"errors"
	"fmt"

	"github.com/gbh007/easyjet/internal/core/entity"
	"github.com/samber/lo"
)

func (srv Service) RunProject(ctx context.Context, id uint) (returnedErr error) {
	p, err := srv.db.Project(ctx, id)
	if err != nil {
		return fmt.Errorf("get project: %w", err)
	}

	dir := p.Dir

	if dir == "" {
		dir = srv.fs.GetProjectDir(ctx, id)
	}

	run := entity.ProjectRun{
		ProjectID: id,
		Success:   true,
	}

	defer func() {
		if returnedErr != nil {
			run.Success = false
			run.FailLog = returnedErr.Error()
		}

		_, saveRunErr := srv.db.SetProjectRun(ctx, run)
		returnedErr = errors.Join(returnedErr, saveRunErr)
	}()

	if p.HasGIT() {
		commits, err := srv.git.Diff(ctx, dir, "HEAD", srv.git.OriginName()+"/"+p.GitBranch)
		if err != nil {
			return fmt.Errorf("git get diff: %w", err)
		}

		run.GitLogs = lo.Map(commits, func(c entity.Commit, i int) entity.ProjectRunGitLogs {
			return entity.ProjectRunGitLogs{
				Number:  i,
				Hash:    c.Hash,
				Subject: c.Subject,
			}
		})

		err = srv.git.Pull(ctx, dir, p.GitBranch)
		if err != nil {
			return fmt.Errorf("git pull origin: %w", err)
		}
	}

	for _, stage := range p.Stages {
		p, err := srv.fs.CreateSHScript(ctx, id, stage.Number, stage.Script)
		if err != nil {
			return fmt.Errorf("create stage %d script: %w", stage.Number, err)
		}

		out, err := srv.ex.Exec(ctx, dir, p)

		run.Stages = append(run.Stages, entity.ProjectRunStage{
			StageNumber: stage.Number,
			Success:     err == nil,
			Log:         out,
		})

		if err != nil {
			return fmt.Errorf("execute stage %d script: %w", stage.Number, err)
		}
	}

	return nil
}
