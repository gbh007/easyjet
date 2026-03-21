package shellexec

import (
	"context"
	"log/slog"

	"github.com/gbh007/easyjet/internal/adapter/internal"
)

const execPath = "/bin/sh"

type Adapter struct {
	logger *slog.Logger
}

func New(logger *slog.Logger) Adapter {
	return Adapter{logger: logger}
}

func (Adapter) Exec(ctx context.Context, dir, p string) (string, error) {
	return internal.Run(ctx, internal.Config{
		Cmd: execPath,
		Args: []string{
			"-xe",
			p,
		},
		Env: []string{
			"WORKSPACE=" + dir,
		},
		Dir:        dir,
		WithStdErr: true,
	})
}
