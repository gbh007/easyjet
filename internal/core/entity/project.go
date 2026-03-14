package entity

import (
	"time"
)

type Project struct {
	ID        uint      `gorm:"column:id;not null;primarykey"`
	CreatedAt time.Time `gorm:"column:created_at;not null;<-:create;autoCreateTime"`
	UpdatedAt time.Time `gorm:"column:updated_at;not null;autoUpdateTime"`

	Dir       string `gorm:"column:dir;not null"`
	GitURL    string `gorm:"column:git_url;not null"`
	GitBranch string `gorm:"column:git_branch;not null"`
	Name      string `gorm:"column:name;not null"`

	Stages []ProjectStage `gorm:"foreignKey:ProjectID"`
}

func (Project) TableName() string {
	return "projects"
}

func (p Project) HasGIT() bool {
	return p.GitURL != "" && p.GitBranch != ""
}

type ProjectStage struct {
	ProjectID uint   `gorm:"column:project_id;not null"`
	Number    int    `gorm:"column:num;not null"`
	Script    string `gorm:"column:script;not null"`
}

func (ProjectStage) TableName() string {
	return "stages"
}
