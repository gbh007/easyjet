package shell

import (
	"context"
	"os/exec"
)

type Config struct {
	Cmd  string
	Args []string
	Dir  string
	Env  []string
}

func Run(ctx context.Context, cfg Config) (string, error) {
	cmd := exec.CommandContext(ctx, cfg.Cmd, cfg.Args...)

	if cfg.Dir != "" {
		cmd.Dir = cfg.Dir
	}

	if len(cfg.Env) > 0 {
		cmd.Env = cfg.Env
	}

	out, err := cmd.Output()
	if err != nil {
		return string(out), err
	}

	return string(out), nil
}
