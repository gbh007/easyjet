package port

import (
	"context"

	"github.com/gbh007/easyjet/internal/core/entity"
)

type Exec interface {
	Exec(ctx context.Context, dir, p string, env []string) (string, error)
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
	Projects(ctx context.Context) ([]entity.Project, error)
	ProjectsWithRunInfo(ctx context.Context, filterType string) ([]entity.ProjectsWithRunInfo, error)
	SetProject(ctx context.Context, pr entity.Project) (uint, error)

	ProjectRun(ctx context.Context, id uint) (entity.ProjectRun, error)
	ProjectRuns(ctx context.Context, id uint) ([]entity.ProjectRun, error)
	ProjectRunIDs(ctx context.Context, id uint) ([]uint, error)
	SetProjectRun(ctx context.Context, run entity.ProjectRun) (uint, error)
	SetProjectRunStage(ctx context.Context, rs entity.ProjectRunStage) error
	SetProjectRunGitCommits(ctx context.Context, commits []entity.ProjectRunGitCommits) error
	DeleteProjectRuns(ctx context.Context, ids []uint) error
	PendingProjectRuns(ctx context.Context) ([]uint, error)

	GlobalEnvVars(ctx context.Context) ([]entity.EnvironmentVariable, error)
	GlobalEnvVar(ctx context.Context, id uint) (entity.EnvironmentVariable, error)
	SetGlobalEnvVar(ctx context.Context, ev entity.EnvironmentVariable) (uint, error)
	DeleteGlobalEnvVar(ctx context.Context, id uint) error
}

type Service interface {
	Project(ctx context.Context, id uint) (entity.Project, error)
	Projects(ctx context.Context) ([]entity.Project, error)
	ProjectsWithRunInfo(ctx context.Context, filterType string) ([]entity.ProjectsWithRunInfo, error)

	CreateProject(ctx context.Context, p entity.Project) (uint, error)
	UpdateProject(ctx context.Context, p entity.Project) error

	RunProject(ctx context.Context, id uint) (uint, error)
	HandleRun(ctx context.Context, runID uint) error
	PendingProjectRuns(ctx context.Context) ([]uint, error)

	ProjectRun(ctx context.Context, runID uint) (entity.ProjectRun, error)
	ProjectRuns(ctx context.Context, id uint) ([]entity.ProjectRun, error)

	GlobalEnvVars(ctx context.Context) ([]entity.EnvironmentVariable, error)
	GlobalEnvVar(ctx context.Context, id uint) (entity.EnvironmentVariable, error)
	CreateGlobalEnvVar(ctx context.Context, ev entity.EnvironmentVariable) (uint, error)
	UpdateGlobalEnvVar(ctx context.Context, ev entity.EnvironmentVariable) error
	DeleteGlobalEnvVar(ctx context.Context, id uint) error
}

type PubSub interface {
	PublishEvent(event entity.Event)
	SubscribeEvent(name string, c int) <-chan entity.Event
}
