package entity

import "time"

const (
	StatusPending    = "pending"
	StatusProcessing = "processing"
	StatusSuccess    = "success"
	StatusFailed     = "failed"
)

type ProjectRun struct {
	ID        int
	CreatedAt time.Time
	UpdatedAt time.Time

	Project    Project
	Stages     []ProjectRunStage
	GitCommits []ProjectRunGitCommits

	ProjectID int
	Status    string
	FailLog   string
	Duration  time.Duration
}

type ProjectRunStage struct {
	RunID       int
	StageNumber int
	Success     bool
	Log         string
	Duration    time.Duration
}

type ProjectRunGitCommits struct {
	RunID   int
	Number  int
	Hash    string
	Subject string
}
