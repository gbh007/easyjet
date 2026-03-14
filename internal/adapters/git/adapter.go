package git

import (
	"context"
	"fmt"
	"log/slog"
	"strings"

	"github.com/gbh007/easyjet/internal/adapters/internal/shell"
	"github.com/gbh007/easyjet/internal/entities"
	"github.com/samber/lo"
)

const (
	execPath   = "/usr/bin/git"
	originName = "origin"
)

type Adapter struct {
	logger *slog.Logger
}

func New(logger *slog.Logger) Adapter {
	return Adapter{logger: logger}
}

func (Adapter) OriginName() string { return originName }

func (Adapter) CurrentHash(ctx context.Context, dir string) (string, error) {
	out, err := shell.Run(ctx, shell.Config{
		Cmd: execPath,
		Args: []string{
			"rev-parse",
			"HEAD",
		},
		Dir: dir,
	})
	if err != nil {
		return "", err
	}

	return strings.TrimSpace(out), nil
}

func (Adapter) Diff(ctx context.Context, dir, from, to string) ([]entities.Commit, error) {
	out, err := shell.Run(ctx, shell.Config{
		Cmd: execPath,
		Args: []string{
			"log",
			"--no-color",
			// Перейти в будущем на
			// --format="{\"hash\":\"%H\",\"author\":\"%an\",\"email\":\"%ae\",\"commit_message\": \"%s\",\"date\": \"%ad\"}"
			"--format=\"%H %s\"",
			from + ".." + to,
		},
		Dir: dir,
	})
	if err != nil {
		return nil, err
	}

	return lo.Map(strings.Split(strings.TrimSpace(out), "\n"), func(s string, _ int) entities.Commit {
		a, b, _ := strings.Cut(s, " ")
		return entities.Commit{
			Hash:    a,
			Subject: b,
		}
	}), nil
}

func (Adapter) Init(ctx context.Context, dir, branch, originURL string) error {
	_, err := shell.Run(ctx, shell.Config{
		Cmd: execPath,
		Args: []string{
			"init",
			"-b",
			branch,
		},
		Dir: dir,
	})
	if err != nil {
		return fmt.Errorf("init: %w", err)
	}

	_, err = shell.Run(ctx, shell.Config{
		Cmd: execPath,
		Args: []string{
			"remote",
			"add",
			originName,
			originURL,
		},
		Dir: dir,
	})
	if err != nil {
		return fmt.Errorf("add remote: %w", err)
	}

	return nil
}

func (Adapter) Pull(ctx context.Context, dir, branch string) error {
	_, err := shell.Run(ctx, shell.Config{
		Cmd: execPath,
		Args: []string{
			"pull",
			"--force",
			originName,
			branch,
		},
		Dir: dir,
	})
	if err != nil {
		return err
	}

	return nil
}
