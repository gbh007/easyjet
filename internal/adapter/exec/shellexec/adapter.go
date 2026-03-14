package shellexec

import (
	"context"

	"github.com/gbh007/easyjet/internal/adapter/internal"
)

const execPath = "/bin/sh"

type Adapter struct{}

func New() Adapter {
	return Adapter{}
}

func (Adapter) Exec(ctx context.Context, dir, p string) (string, error) {
	return internal.Run(ctx, internal.Config{
		Cmd: execPath,
		Args: []string{
			"-xe",
			p,
		},
		Dir: dir,
	})
}
