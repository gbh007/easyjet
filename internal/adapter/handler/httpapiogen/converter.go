package httpapiogen

import (
	"github.com/gbh007/easyjet/internal/adapter/handler/httpapiogen/ogenapi"
	"github.com/gbh007/easyjet/internal/core/entity"
)

func convertProjectCreate(req *ogenapi.ProjectCreate) entity.Project {
	p := entity.Project{
		Name:           req.Name,
		Dir:            req.Dir.Or(""),
		CronEnabled:    req.CronEnabled.Or(false),
		CronSchedule:   req.CronSchedule.Or(""),
		GitURL:         req.GitURL.Or(""),
		GitBranch:      req.GitBranch.Or(""),
		RestartAfter:   req.RestartAfter.Or(false),
		RetentionCount: int(req.RetentionCount.Or(10)),
		WithRootEnv:    req.WithRootEnv.Or(false),
		Stages:         make([]entity.ProjectStage, len(req.Stages)),
	}

	for i, stage := range req.Stages {
		p.Stages[i] = entity.ProjectStage{
			Number: int(stage.Number),
			Script: stage.Script,
		}
	}

	return p
}

func convertProjectUpdate(req *ogenapi.ProjectUpdate, projectID uint) entity.Project {
	p := entity.Project{
		ID:             projectID,
		Name:           req.Name,
		Dir:            req.Dir.Or(""),
		CronEnabled:    req.CronEnabled.Or(false),
		CronSchedule:   req.CronSchedule.Or(""),
		GitURL:         req.GitURL.Or(""),
		GitBranch:      req.GitBranch.Or(""),
		RestartAfter:   req.RestartAfter.Or(false),
		RetentionCount: int(req.RetentionCount.Or(10)),
		WithRootEnv:    req.WithRootEnv.Or(false),
		Stages:         make([]entity.ProjectStage, len(req.Stages)),
	}

	for i, stage := range req.Stages {
		p.Stages[i] = entity.ProjectStage{
			Number: int(stage.Number),
			Script: stage.Script,
		}
	}

	return p
}

func convertProjectToOgen(p entity.Project) ogenapi.Project {
	stages := make([]ogenapi.ProjectStage, len(p.Stages))
	for i, stage := range p.Stages {
		stages[i] = ogenapi.ProjectStage{
			Number: int32(stage.Number),
			Script: stage.Script,
		}
	}

	return ogenapi.Project{
		ID:             ogenapi.NewOptUint(p.ID),
		CreatedAt:      ogenapi.NewOptDateTime(p.CreatedAt),
		UpdatedAt:      ogenapi.NewOptDateTime(p.UpdatedAt),
		CronEnabled:    ogenapi.NewOptBool(p.CronEnabled),
		CronSchedule:   ogenapi.NewOptString(p.CronSchedule),
		Name:           p.Name,
		Dir:            ogenapi.NewOptString(p.Dir),
		GitURL:         ogenapi.NewOptString(p.GitURL),
		GitBranch:      ogenapi.NewOptString(p.GitBranch),
		RestartAfter:   ogenapi.NewOptBool(p.RestartAfter),
		RetentionCount: ogenapi.NewOptInt32(int32(p.RetentionCount)),
		WithRootEnv:    ogenapi.NewOptBool(p.WithRootEnv),
		Stages:         stages,
	}
}

func convertProjectRunToOgen(run entity.ProjectRun) ogenapi.ProjectRun {
	stages := make([]ogenapi.ProjectRunStage, len(run.Stages))
	for i, stage := range run.Stages {
		stages[i] = ogenapi.ProjectRunStage{
			StageNumber: ogenapi.NewOptInt32(int32(stage.StageNumber)),
			Success:     ogenapi.NewOptBool(stage.Success),
			Log:         ogenapi.NewOptString(stage.Log),
		}
	}

	commits := make([]ogenapi.ProjectRunGitCommit, len(run.GitCommits))
	for i, commit := range run.GitCommits {
		commits[i] = ogenapi.ProjectRunGitCommit{
			Number:  ogenapi.NewOptInt32(int32(commit.Number)),
			Hash:    ogenapi.NewOptString(commit.Hash),
			Subject: ogenapi.NewOptString(commit.Subject),
		}
	}

	return ogenapi.ProjectRun{
		ID:         ogenapi.NewOptUint(run.ID),
		CreatedAt:  ogenapi.NewOptDateTime(run.CreatedAt),
		UpdatedAt:  ogenapi.NewOptDateTime(run.UpdatedAt),
		ProjectID:  ogenapi.NewOptUint(run.ProjectID),
		Success:    ogenapi.NewOptBool(run.Success),
		Pending:    ogenapi.NewOptBool(run.Pending),
		Processing: ogenapi.NewOptBool(run.Processing),
		FailLog:    ogenapi.NewOptString(run.FailLog),
		Stages:     stages,
		GitCommits: commits,
	}
}

func convertProjectRunsToOgen(runs []entity.ProjectRun) []ogenapi.ProjectRun {
	result := make([]ogenapi.ProjectRun, len(runs))
	for i, run := range runs {
		result[i] = convertProjectRunToOgen(run)
	}
	return result
}

func convertProjectsToOgen(projects []entity.Project) []ogenapi.Project {
	result := make([]ogenapi.Project, len(projects))
	for i, project := range projects {
		result[i] = convertProjectToOgen(project)
	}
	return result
}
