package httpapi

import (
	"github.com/gbh007/easyjet/internal/adapter/handler/httpapi/ogenapi"
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

	if envVars, ok := req.EnvVars.Get(); ok {
		for _, ev := range envVars {
			p.EnvVars = append(p.EnvVars, entity.EnvironmentVariable{
				Name:               ev.Name,
				Value:              ev.Value,
				UsesOtherVariables: ev.UsesOtherVariables.Value,
			})
		}
	}

	return p
}
