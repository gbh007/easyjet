package internal

import (
	"context"
	"os/exec"
)

type Config struct {
	Cmd        string
	Args       []string
	Dir        string
	Env        []string
	WithStdErr bool
}

func Run(ctx context.Context, cfg Config) (string, error) {
	cmd := exec.CommandContext(ctx, cfg.Cmd, cfg.Args...)

	if cfg.Dir != "" {
		cmd.Dir = cfg.Dir
	}

	if len(cfg.Env) > 0 {
		cmd.Env = cfg.Env
	}

	if !cfg.WithStdErr {
		out, err := cmd.Output()
		if err != nil {
			return string(out), err
		}

		return string(out), nil
	}

	out, err := cmd.CombinedOutput()
	if err != nil {
		return string(out), err
	}

	return string(out), nil
}
