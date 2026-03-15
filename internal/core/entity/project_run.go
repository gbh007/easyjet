package entity

import "time"

type ProjectRun struct {
	ID        uint      `param:"run_id" json:"id" gorm:"column:id;not null;primarykey"`
	CreatedAt time.Time `json:"created_at" gorm:"column:created_at;not null;<-:create;autoCreateTime"`
	UpdatedAt time.Time `json:"updated_at" gorm:"column:updated_at;not null;autoUpdateTime"`

	Project Project             `json:"project,omitzero" gorm:"foreignKey:ProjectID"`
	Stages  []ProjectRunStage   `json:"stages,omitempty,omitzero" gorm:"foreignKey:RunID" validate:"min=1"`
	GitLogs []ProjectRunGitLogs `json:"git_logs,omitempty,omitzero" gorm:"foreignKey:RunID"`

	ProjectID uint   `json:"project_id" gorm:"column:project_id;not null;index:idx_project_id"`
	Success   bool   `json:"success" gorm:"column:success;not null"`
	FailLog   string `json:"fail_log" gorm:"column:fail_log;not null"`
}

func (ProjectRun) TableName() string {
	return "runs"
}

type ProjectRunStage struct {
	RunID       uint   `json:"run_id" gorm:"column:run_id;not null;index:idx_run_id"`
	StageNumber int    `json:"stage_number" gorm:"column:stage_num;not null" validate:"min=1"`
	Success     bool   `json:"success" gorm:"column:success;not null"`
	Log         string `json:"log" gorm:"column:log;not null"`
}

func (ProjectRunStage) TableName() string {
	return "run_stages"
}

type ProjectRunGitLogs struct {
	RunID   uint   `json:"run_id" gorm:"column:run_id;not null;index:idx_run_id"`
	Number  int    `json:"number" gorm:"column:num;not null"`
	Hash    string `json:"hash" gorm:"column:hash;not null"`
	Subject string `json:"subject" gorm:"column:subject;not null"`
}

func (ProjectRunGitLogs) TableName() string {
	return "run_git_logs"
}
