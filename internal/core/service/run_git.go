package service

import (
	"context"
	"fmt"
	"time"

	"github.com/gbh007/easyjet/internal/core/entity"
	"github.com/samber/lo"
)

func (srv Service) runGitPull(ctx context.Context, project entity.Project, runID uint, dir string) (err error) {
	start := time.Now()

	defer func() {
		srv.pubsub.PublishEvent(entity.Event{
			Type:      entity.EventRunGitFinished,
			ProjectID: project.ID,
			RunID:     runID,
			Err:       err,
			Duration:  time.Since(start),
		})
	}()

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

	return nil
}
