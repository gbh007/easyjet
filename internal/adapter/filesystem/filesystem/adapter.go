package filesystem

import (
	"context"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strconv"
)

type Adapter struct {
	basePath string
}

func New(basePath string) Adapter {
	return Adapter{basePath: basePath}
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

	_, err = io.WriteString(f, body)
	if err != nil {
		return "", fmt.Errorf("write script: %w", err)
	}

	return p, nil
}
