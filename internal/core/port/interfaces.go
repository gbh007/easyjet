package port

import (
	"context"

	"github.com/gbh007/easyjet/internal/core/entity"
)

type Exec interface {
	Exec(ctx context.Context, dir, p string, env []string) (string, error)
}

type FileSystem interface {
	GetProjectDir(ctx context.Context, id int) string
	CreateProjectDir(ctx context.Context, id int) (string, error)
	ProjectDirExists(ctx context.Context, id int) (bool, error)
	CreateSHScript(ctx context.Context, id, stage int, body string) (p string, err error)
}

type Git interface {
	OriginName() string
	CurrentHash(ctx context.Context, dir string) (string, error)
	Diff(ctx context.Context, dir, from, to string) ([]entity.Commit, error)
	Init(ctx context.Context, dir, branch, originURL string) error
	Pull(ctx context.Context, dir, branch string) error
	Exists(ctx context.Context, dir string) (bool, error)
	CurrentBranch(ctx context.Context, dir string) (string, error)
	SwitchBranch(ctx context.Context, dir, branch string, create bool) error
	CurrentOriginURL(ctx context.Context, dir string) (string, error)
	Branches(ctx context.Context, dir string) ([]string, error)
	DeleteBranch(ctx context.Context, dir, branch string) error
	GC(ctx context.Context, dir string) error
	HardReset(ctx context.Context, dir, branch string) error
	SetOriginURL(ctx context.Context, dir, originURL string) error
	Fetch(ctx context.Context, dir string) error
}

type Database interface {
	Project(ctx context.Context, id int) (entity.Project, error)
	ProjectsWithRunInfo(ctx context.Context, filterType entity.ProjectFilterType) ([]entity.ProjectsWithRunInfo, error)
	SetProject(ctx context.Context, pr entity.Project) (int, error)

	ProjectRun(ctx context.Context, id int) (entity.ProjectRun, error)
	ProjectRuns(ctx context.Context, id int) ([]entity.ProjectRun, error)
	ProjectRunIDs(ctx context.Context, id int) ([]int, error)
	SetProjectRun(ctx context.Context, run entity.ProjectRun) (int, error)
	SetProjectRunStage(ctx context.Context, rs entity.ProjectRunStage) error
	SetProjectRunGitCommits(ctx context.Context, commits []entity.ProjectRunGitCommits) error
	DeleteProjectRuns(ctx context.Context, ids []int) error
	PendingProjectRuns(ctx context.Context) ([]int, error)

	GlobalEnvVars(ctx context.Context) ([]entity.EnvironmentVariable, error)
	GlobalEnvVar(ctx context.Context, id int) (entity.EnvironmentVariable, error)
	SetGlobalEnvVar(ctx context.Context, ev entity.EnvironmentVariable) (int, error)
	DeleteGlobalEnvVar(ctx context.Context, id int) error
}

type Service interface {
	Project(ctx context.Context, id int) (entity.Project, error)
	ProjectsWithRunInfo(ctx context.Context, filterType entity.ProjectFilterType) ([]entity.ProjectsWithRunInfo, error)

	CreateProject(ctx context.Context, p entity.Project) (int, error)
	UpdateProject(ctx context.Context, p entity.Project) error

	RunProject(ctx context.Context, id int) (int, error)
	HandleRun(ctx context.Context, runID int) error
	PendingProjectRuns(ctx context.Context) ([]int, error)

	ProjectRun(ctx context.Context, runID int) (entity.ProjectRun, error)
	ProjectRuns(ctx context.Context, id int) ([]entity.ProjectRun, error)

	GlobalEnvVars(ctx context.Context) ([]entity.EnvironmentVariable, error)
	GlobalEnvVar(ctx context.Context, id int) (entity.EnvironmentVariable, error)
	CreateGlobalEnvVar(ctx context.Context, ev entity.EnvironmentVariable) (int, error)
	UpdateGlobalEnvVar(ctx context.Context, ev entity.EnvironmentVariable) error
	DeleteGlobalEnvVar(ctx context.Context, id int) error
}

type PubSub interface {
	PublishEvent(event entity.Event)
	SubscribeEvent(name string, c int) <-chan entity.Event
}
