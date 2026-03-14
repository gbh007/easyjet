package port

import (
	"context"

	"github.com/gbh007/easyjet/internal/core/entity"
)

type Exec interface {
	Exec(ctx context.Context, dir, p string) (string, error)
}

type FileSystem interface {
	GetProjectDir(ctx context.Context, id uint) string
	CreateProjectDir(ctx context.Context, id uint) (string, error)
	CreateSHScript(ctx context.Context, id uint, stage int, body string) (p string, err error)
}

type Git interface {
	OriginName() string
	CurrentHash(ctx context.Context, dir string) (string, error)
	Diff(ctx context.Context, dir, from, to string) ([]entity.Commit, error)
	Init(ctx context.Context, dir, branch, originURL string) error
	Pull(ctx context.Context, dir, branch string) error
}

type Database interface {
	Project(ctx context.Context, id uint) (entity.Project, error)
	SetProject(ctx context.Context, pr entity.Project) (uint, error)
	ProjectRun(ctx context.Context, id uint) (entity.ProjectRun, error)
	ProjectRuns(ctx context.Context, id uint) ([]entity.ProjectRun, error)
	SetProjectRun(ctx context.Context, run entity.ProjectRun) (uint, error)
}

type Service interface {
	Project(ctx context.Context, id uint) (entity.Project, error)

	CreateProject(ctx context.Context, p entity.Project) (uint, error)
	UpdateProject(ctx context.Context, p entity.Project) error

	RunProject(ctx context.Context, id uint) (returnedErr error)

	ProjectRun(ctx context.Context, id uint) (entity.ProjectRun, error)
	ProjectRuns(ctx context.Context, id uint) ([]entity.ProjectRun, error)
}
