package filesystem

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"os"
	"path/filepath"
	"strconv"
)

type Adapter struct {
	basePath string
	logger   *slog.Logger
}

func New(logger *slog.Logger, basePath string) Adapter {
	return Adapter{
		basePath: basePath,
		logger:   logger,
	}
}

func (a Adapter) GetProjectDir(ctx context.Context, id uint) string {
	return filepath.Join(a.basePath, strconv.Itoa(int(id)))
}

func (a Adapter) CreateProjectDir(ctx context.Context, id uint) (string, error) {
	p := a.GetProjectDir(ctx, id)

	err := os.MkdirAll(p, os.ModeDir|os.ModePerm)
	if err != nil {
		return "", err
	}

	return p, nil
}

func (a Adapter) ProjectDirExists(ctx context.Context, id uint) (bool, error) {
	p := a.GetProjectDir(ctx, id)

	info, err := os.Stat(p)
	if errors.Is(err, os.ErrNotExist) {
		return false, nil
	}

	if err != nil {
		return false, err
	}

	if info.IsDir() {
		return true, nil
	}

	return false, errors.New("is not dir")
}

func (Adapter) CreateSHScript(ctx context.Context, id uint, stage int, body string) (p string, err error) {
	d := filepath.Join(os.TempDir(), "easyjet", "scripts")

	err = os.MkdirAll(d, os.ModeDir|os.ModePerm)
	if err != nil {
		return "", fmt.Errorf("mk temp dir: %w", err)
	}

	p = filepath.Join(d, fmt.Sprintf("%d_%d.sh", id, stage))

	f, err := os.Create(p)
	if err != nil {
		return "", fmt.Errorf("create script: %w", err)
	}

	defer func() {
		err = errors.Join(err, f.Close())
	}()

	_, err = f.WriteString(body)
	if err != nil {
		return "", fmt.Errorf("write script: %w", err)
	}

	return p, nil
}
