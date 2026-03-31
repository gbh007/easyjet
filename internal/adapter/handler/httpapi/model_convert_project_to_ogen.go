package httpapi

import (
	"github.com/gbh007/easyjet/internal/adapter/handler/httpapi/ogenapi"
	"github.com/gbh007/easyjet/internal/core/entity"
)

func convertProjectToOgen(p entity.Project) ogenapi.Project {
	stages := make([]ogenapi.ProjectStage, len(p.Stages))
	for i, stage := range p.Stages {
		stages[i] = ogenapi.ProjectStage{
			Number: int32(stage.Number),
			Script: stage.Script,
		}
	}

	envVars := make([]ogenapi.ProjectEnvironmentVariable, 0, len(p.EnvVars))
	for _, ev := range p.EnvVars {
		envVars = append(envVars, ogenapi.ProjectEnvironmentVariable{
			ID:                 ogenapi.NewOptUint(ev.ID),
			Name:               ev.Name,
			Value:              ev.Value,
			UsesOtherVariables: ogenapi.NewOptBool(ev.UsesOtherVariables),
		})
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
		IsTemplate:     ogenapi.NewOptBool(p.IsTemplate),
		Stages:         stages,
		EnvVars:        ogenapi.NewOptNilProjectEnvironmentVariableArray(envVars),
	}
}
