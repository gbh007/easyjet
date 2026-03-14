package entity

import "time"

type ProjectRun struct {
	ID        uint      `gorm:"column:id;not null;primarykey"`
	CreatedAt time.Time `gorm:"column:created_at;not null;<-:create;autoCreateTime"`
	UpdatedAt time.Time `gorm:"column:updated_at;not null;autoUpdateTime"`

	Project Project             `gorm:"foreignKey:ProjectID"`
	Stages  []ProjectRunStage   `gorm:"foreignKey:RunID"`
	GitLogs []ProjectRunGitLogs `gorm:"foreignKey:RunID"`

	ProjectID uint   `gorm:"column:project_id;not null"`
	Success   bool   `gorm:"column:success;not null"`
	FailLog   string `gorm:"column:fail_log;not null"`
}

func (ProjectRun) TableName() string {
	return "runs"
}

type ProjectRunStage struct {
	RunID       uint   `gorm:"column:run_id;not null"`
	StageNumber int    `gorm:"column:stage_num;not null"`
	Success     bool   `gorm:"column:success;not null"`
	Log         string `gorm:"column:log;not null"`
}

func (ProjectRunStage) TableName() string {
	return "run_stages"
}

type ProjectRunGitLogs struct {
	RunID   uint   `gorm:"column:run_id;not null"`
	Number  int    `gorm:"column:num;not null"`
	Hash    string `gorm:"column:hash;not null"`
	Subject string `gorm:"column:subject;not null"`
}

func (ProjectRunGitLogs) TableName() string {
	return "run_git_logs"
}
