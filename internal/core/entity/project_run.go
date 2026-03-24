package entity

import "time"

type ProjectRun struct {
	ID        uint
	CreatedAt time.Time
	UpdatedAt time.Time

	Project    Project
	Stages     []ProjectRunStage
	GitCommits []ProjectRunGitCommits

	ProjectID  uint
	Success    bool
	Pending    bool
	Processing bool
	FailLog    string
	Duration   time.Duration
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
