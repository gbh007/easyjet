package httpapi

import (
	"github.com/gbh007/easyjet/internal/adapter/handler/httpapi/ogenapi"
	"github.com/gbh007/easyjet/internal/core/entity"
)

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
