package entity

import "time"

const (
	StatusPending    = "pending"
	StatusProcessing = "processing"
	StatusSuccess    = "success"
	StatusFailed     = "failed"
)

type ProjectRun struct {
	ID        uint
	CreatedAt time.Time
	UpdatedAt time.Time

	Project    Project
	Stages     []ProjectRunStage
	GitCommits []ProjectRunGitCommits

	ProjectID uint
	Status    string
	FailLog   string
	Duration  time.Duration
}

type ProjectRunStage struct {
	RunID       uint
	StageNumber int
	Success     bool
	Log         string
	Duration    time.Duration
}

type ProjectRunGitCommits struct {
	RunID   uint
	Number  int
	Hash    string
	Subject string
}
