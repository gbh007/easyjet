package entity

import (
	"time"
)

type Project struct {
	ID        uint      `param:"project_id" json:"id" gorm:"column:id;not null;primarykey"`
	CreatedAt time.Time `json:"created_at" gorm:"column:created_at;not null;<-:create;autoCreateTime"`
	UpdatedAt time.Time `json:"updated_at" gorm:"column:updated_at;not null;autoUpdateTime"`

	CronEnabled  bool   `json:"cron_enabled" gorm:"column:cron_enabled;not null;default:false"`
	CronSchedule string `json:"cron_schedule" gorm:"column:cron_schedule;type:text;not null;default:''" validate:"omitzero,cron"`

	Dir            string `json:"dir" gorm:"column:dir;not null"`
	GitURL         string `json:"git_url" gorm:"column:git_url;not null"`
	GitBranch      string `json:"git_branch" gorm:"column:git_branch;not null"`
	Name           string `json:"name" gorm:"column:name;not null"`
	RestartAfter   bool   `json:"restart_after" gorm:"column:restart_after;not null;default:false"`
	RetentionCount int    `json:"retention_count" gorm:"column:retention_count;not null;default:0"`
	WithRootEnv    bool   `json:"with_root_env" gorm:"column:with_root_env;not null;default:false"`

	Stages []ProjectStage `json:"stages" gorm:"foreignKey:ProjectID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE" validate:"min=1,dive"`
}

func (Project) TableName() string {
	return "projects"
}

func (p Project) HasGIT() bool {
	return p.GitURL != "" && p.GitBranch != ""
}

type ProjectStage struct {
	ProjectID uint   `json:"project_id" gorm:"column:project_id;not null;primaryKey;autoIncrement:false"`
	Number    int    `json:"number" gorm:"column:num;not null;primaryKey;autoIncrement:false" validate:"min=1"`
	Script    string `json:"script" gorm:"column:script;not null" validate:"required,gt=1"`
}

func (ProjectStage) TableName() string {
	return "stages"
}
