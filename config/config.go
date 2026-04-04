package config

import (
	"errors"
	"fmt"
	"log/slog"
	"os"
	"strings"

	"github.com/BurntSushi/toml"
)

type Config struct {
	App      App      `toml:"app"`
	Log      Log      `toml:"log"`
	Database Database `toml:"database"`
	Server   Server   `toml:"server"`
	MCP      MCP      `toml:"mcp"`
}

type App struct {
	ProjectDir string `toml:"project_dir"`
}

type Server struct {
	Addr            string         `toml:"addr"`
	User            string         `toml:"user"`
	Pass            string         `toml:"pass"`
	StaticFilesPath string         `toml:"static_files_path"`
	External        ServerExternal `toml:"external"`
}

type ServerExternal struct {
	API string `toml:"api"`
	Web string `toml:"web"`
}

type Database struct {
	Type string `toml:"type"`
	DNS  string `toml:"dns"`
}

type Log struct {
	Level  string `toml:"level"`
	Format string `toml:"format"`
}

type MCP struct {
	Enabled        bool `toml:"enabled"`
	AllowRuns      bool `toml:"allow_runs"`
	AllowMutations bool `toml:"allow_mutations"`
}

func (l Log) SlogLevel() slog.Level {
	switch strings.ToLower(l.Level) {
	case "debug", "dbg":
		return slog.LevelDebug
	case "info", "inf":
		return slog.LevelInfo
	case "warn", "warning", "wrn":
		return slog.LevelWarn
	case "error", "err":
		return slog.LevelError
	default:
		return slog.LevelInfo
	}
}

func Write(filename string, cfg Config) (returnedErr error) {
	f, err := os.Create(filename)
	if err != nil {
		return fmt.Errorf("create file: %w", err)
	}

	defer func() {
		err := f.Close()
		if err != nil {
			returnedErr = errors.Join(returnedErr, fmt.Errorf("close config file: %w", err))
		}
	}()

	enc := toml.NewEncoder(f)
	enc.Indent = ""

	err = enc.Encode(cfg)
	if err != nil {
		return fmt.Errorf("encode toml: %w", err)
	}

	return nil
}

func Read[T any](filename string) (cfg T, returnedErr error) {
	f, err := os.Open(filename)
	if err != nil {
		return cfg, fmt.Errorf("open config file: %w", err)
	}

	defer func() {
		err := f.Close()
		if err != nil {
			returnedErr = errors.Join(returnedErr, fmt.Errorf("close config file: %w", err))
		}
	}()

	_, err = toml.NewDecoder(f).Decode(&cfg)
	if err != nil {
		return cfg, fmt.Errorf("decode toml: %w", err)
	}

	return cfg, nil
}
