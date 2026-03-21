package shellexec

import (
	"context"
	"log/slog"
	"os"

	"github.com/gbh007/easyjet/internal/adapter/internal"
)

const execPath = "/bin/sh"

type Adapter struct {
	logger *slog.Logger
}

func New(logger *slog.Logger) Adapter {
	return Adapter{logger: logger}
}

func (Adapter) Exec(ctx context.Context, dir, p string, withRootEnv bool) (string, error) {
	var env []string

	if withRootEnv {
		env = append(env, os.Environ()...)
	}

	env = append(env, "WORKSPACE="+dir)

	return internal.Run(ctx, internal.Config{
		Cmd: execPath,
		Args: []string{
			"-xe",
			p,
		},
		Env:        env,
		Dir:        dir,
		WithStdErr: true,
	})
}
