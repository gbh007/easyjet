package exec

import (
	"context"

	"github.com/gbh007/easyjet/internal/adapters/internal/shell"
)

const execPath = "/bin/sh"

type Adapter struct{}

func New() Adapter {
	return Adapter{}
}

func (Adapter) Exec(ctx context.Context, dir, p string) (string, error) {
	return shell.Run(ctx, shell.Config{
		Cmd: execPath,
		Args: []string{
			"-xe",
			p,
		},
		Dir: dir,
	})
}
