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
			Duration:    ogenapi.NewOptInt64(stage.Duration.Milliseconds()),
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
		ID:         ogenapi.NewOptInt(run.ID),
		CreatedAt:  ogenapi.NewOptDateTime(run.CreatedAt),
		UpdatedAt:  ogenapi.NewOptDateTime(run.UpdatedAt),
		ProjectID:  ogenapi.NewOptInt(run.ProjectID),
		Status:     ogenapi.NewOptProjectRunStatus(ogenapi.ProjectRunStatus(run.Status)),
		FailLog:    ogenapi.NewOptString(run.FailLog),
		Duration:   ogenapi.NewOptInt64(run.Duration.Milliseconds()),
		Stages:     ogenapi.NewOptNilProjectRunStageArray(stages),
		GitCommits: ogenapi.NewOptNilProjectRunGitCommitArray(commits),
	}
}
