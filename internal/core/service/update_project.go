package service

import (
	"context"
	"fmt"
	"slices"

	"github.com/gbh007/easyjet/internal/core/entity"
	"github.com/samber/lo"
)

func (srv Service) UpdateProject(ctx context.Context, p entity.Project) error {
	if p.Dir == "" {
		ok, err := srv.fs.ProjectDirExists(ctx, p.ID)
		if err != nil {
			return fmt.Errorf("check project dir: %w", err)
		}

		if !ok {
			_, err = srv.fs.CreateProjectDir(ctx, p.ID)
			if err != nil {
				return fmt.Errorf("create project dir: %w", err)
			}
		}
	}

	if p.HasGIT() {
		err := srv.resolveGit(ctx, p)
		if err != nil {
			return fmt.Errorf("resolve git: %w", err)
		}
	}

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

	return nil
}

func (srv Service) resolveGit(ctx context.Context, p entity.Project) error {
	dir := p.Dir

	if dir == "" {
		dir = srv.fs.GetProjectDir(ctx, p.ID)
	}

	ok, err := srv.git.Exists(ctx, dir)
	if err != nil {
		return fmt.Errorf("check git exists: %w", err)
	}

	if !ok {
		err = srv.git.Init(ctx, dir, p.GitBranch, p.GitURL)
		if err != nil {
			return fmt.Errorf("init git: %w", err)
		}

		err = srv.git.Pull(ctx, dir, p.GitBranch)
		if err != nil {
			return fmt.Errorf("pull git: %w", err)
		}

		return nil
	}

	oldBranch, err := srv.git.CurrentBranch(ctx, dir)
	if err != nil {
		return fmt.Errorf("get git branch: %w", err)
	}

	oldURL, err := srv.git.CurrentOriginURL(ctx, dir)
	if err != nil {
		return fmt.Errorf("get git remote: %w", err)
	}

	if p.GitBranch == oldBranch && p.GitURL == oldURL {
		return nil
	}

	branches, err := srv.git.Branches(ctx, dir)
	if err != nil {
		return fmt.Errorf("get git branches: %w", err)
	}

	if p.GitURL != oldURL {
		err = srv.git.SetOriginURL(ctx, dir, p.GitURL)
		if err != nil {
			return fmt.Errorf("set git remote: %w", err)
		}
	}

	if p.GitBranch != oldBranch {
		newExists := slices.Contains(branches, p.GitBranch)

		err = srv.git.SwitchBranch(ctx, dir, p.GitBranch, !newExists)
		if err != nil {
			return fmt.Errorf("switch git branch: %w", err)
		}
	}

	err = srv.git.Fetch(ctx, dir)
	if err != nil {
		return fmt.Errorf("git fetch: %w", err)
	}

	err = srv.git.HardReset(ctx, dir, srv.git.OriginName()+"/"+p.GitBranch)
	if err != nil {
		return fmt.Errorf("git reset: %w", err)
	}

	for _, b := range branches {
		if b == p.GitBranch {
			continue
		}

		err = srv.git.DeleteBranch(ctx, dir, b)
		if err != nil { // По неизвестной причине удаляет ветку но выходит с кодом 1, поэтому просто логируем
			srv.logger.WarnContext(ctx, "fail delete git branch", "error", err)
		}
	}

	err = srv.git.GC(ctx, dir)
	if err != nil {
		return fmt.Errorf("git gc: %w", err)
	}

	return nil
}
