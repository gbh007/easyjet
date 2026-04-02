package shellgit

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"os"
	"path/filepath"
	"slices"
	"strings"

	"github.com/gbh007/easyjet/internal/adapter/internal"
	"github.com/gbh007/easyjet/internal/core/entity"
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

func (Adapter) OriginName() string {
	return originName
}

func (Adapter) CurrentHash(ctx context.Context, dir string) (string, error) {
	out, err := internal.Run(ctx, internal.Config{
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

func (adp Adapter) Diff(ctx context.Context, dir, from, to string) ([]entity.Commit, error) {
	_, err := internal.Run(ctx, internal.Config{
		Cmd: execPath,
		Args: []string{
			"fetch",
			"--force",
			adp.OriginName(),
		},
		Dir: dir,
	})
	if err != nil {
		return nil, fmt.Errorf("fetch: %w", err)
	}

	out, err := internal.Run(ctx, internal.Config{
		Cmd: execPath,
		Args: []string{
			"log",
			"--no-color",
			// Перейти в будущем на
			// --format="{\"hash\":\"%H\",\"author\":\"%an\",\"email\":\"%ae\",\"commit_message\": \"%s\",\"date\": \"%ad\"}"
			"--format=%H %s",
			from + ".." + to,
		},
		Dir: dir,
	})
	if err != nil {
		return nil, fmt.Errorf("log: %w", err)
	}

	result := lo.Filter(
		lo.Map(strings.Split(strings.TrimSpace(out), "\n"), func(s string, _ int) entity.Commit {
			a, b, _ := strings.Cut(s, " ")
			return entity.Commit{
				Hash:    a,
				Subject: b,
			}
		}),
		func(c entity.Commit, _ int) bool {
			return c.Hash != "" || c.Subject != ""
		},
	)

	slices.Reverse(result)

	return result, nil
}

func (adp Adapter) Init(ctx context.Context, dir, branch, originURL string) error {
	_, err := internal.Run(ctx, internal.Config{
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

	_, err = internal.Run(ctx, internal.Config{
		Cmd: execPath,
		Args: []string{
			"remote",
			"add",
			adp.OriginName(),
			originURL,
		},
		Dir: dir,
	})
	if err != nil {
		return fmt.Errorf("add remote: %w", err)
	}

	return nil
}

func (adp Adapter) Pull(ctx context.Context, dir, branch string) error {
	_, err := internal.Run(ctx, internal.Config{
		Cmd: execPath,
		Args: []string{
			"pull",
			"--force",
			adp.OriginName(),
			branch,
		},
		Dir: dir,
	})
	if err != nil {
		return err
	}

	return nil
}

func (adp Adapter) Fetch(ctx context.Context, dir string) error {
	_, err := internal.Run(ctx, internal.Config{
		Cmd: execPath,
		Args: []string{
			"fetch",
			"--force",
			adp.OriginName(),
		},
		Dir: dir,
	})
	if err != nil {
		return err
	}

	return nil
}

func (adp Adapter) HardReset(ctx context.Context, dir, branch string) error {
	_, err := internal.Run(ctx, internal.Config{
		Cmd: execPath,
		Args: []string{
			"reset",
			"--hard",
			branch,
		},
		Dir: dir,
	})
	if err != nil {
		return err
	}

	return nil
}

func (adp Adapter) Exists(ctx context.Context, dir string) (bool, error) {
	info, err := os.Stat(filepath.Join(dir, ".git"))
	if errors.Is(err, os.ErrNotExist) {
		return false, nil
	}

	if err != nil {
		return false, err
	}

	if info.IsDir() {
		return true, nil
	}

	return false, errors.New(".git is not dir")
}

func (adp Adapter) CurrentBranch(ctx context.Context, dir string) (string, error) {
	name, err := internal.Run(ctx, internal.Config{
		Cmd: execPath,
		Args: []string{
			"branch",
			"--show-current",
		},
		Dir: dir,
	})
	if err != nil {
		return "", err
	}

	return strings.TrimSpace(name), nil
}

func (adp Adapter) CurrentOriginURL(ctx context.Context, dir string) (string, error) {
	u, err := internal.Run(ctx, internal.Config{
		Cmd: execPath,
		Args: []string{
			"remote",
			"get-url",
			adp.OriginName(),
		},
		Dir: dir,
	})
	if err != nil {
		return "", err
	}

	return strings.TrimSpace(u), nil
}

func (adp Adapter) SetOriginURL(ctx context.Context, dir, originURL string) error {
	_, err := internal.Run(ctx, internal.Config{
		Cmd: execPath,
		Args: []string{
			"remote",
			"set-url",
			adp.OriginName(),
			originURL,
		},
		Dir: dir,
	})
	if err != nil {
		return err
	}

	return nil
}

func (adp Adapter) Branches(ctx context.Context, dir string) ([]string, error) {
	name, err := internal.Run(ctx, internal.Config{
		Cmd: execPath,
		Args: []string{
			"branch",
			"--no-color",
			`--format=%(refname:short)`,
		},
		Dir: dir,
	})
	if err != nil {
		return nil, err
	}

	return lo.Map(strings.Split(name, "\n"), func(s string, _ int) string {
		return strings.TrimSpace(s)
	}), nil
}

func (adp Adapter) DeleteBranch(ctx context.Context, dir, branch string) error {
	_, err := internal.Run(ctx, internal.Config{
		Cmd: execPath,
		Args: []string{
			"branch",
			"-D",
			branch,
		},
		Dir: dir,
	})
	if err != nil {
		return err
	}

	return nil
}

func (adp Adapter) GC(ctx context.Context, dir string) error {
	_, err := internal.Run(ctx, internal.Config{
		Cmd: execPath,
		Args: []string{
			"gc",
		},
		Dir: dir,
	})
	if err != nil {
		return err
	}

	return nil
}

func (adp Adapter) SwitchBranch(ctx context.Context, dir, branch string, create bool) error {
	args := []string{
		"checkout",
		branch,
	}

	if create {
		args = []string{
			"checkout",
			"-b",
			branch,
		}
	}

	_, err := internal.Run(ctx, internal.Config{
		Cmd:  execPath,
		Args: args,
		Dir:  dir,
	})
	if err != nil {
		return err
	}

	return nil
}
