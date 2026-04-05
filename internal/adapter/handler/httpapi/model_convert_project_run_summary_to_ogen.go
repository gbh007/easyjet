package httpapi

import (
	"github.com/gbh007/easyjet/internal/adapter/handler/httpapi/ogenapi"
	"github.com/gbh007/easyjet/internal/core/entity"
)

func convertProjectRunsSummaryToOgen(runs []entity.ProjectRun) []ogenapi.ProjectRunSummary {
	result := make([]ogenapi.ProjectRunSummary, len(runs))
	for i, run := range runs {
		result[i] = convertProjectRunSummaryToOgen(run)
	}
	return result
}

func convertProjectRunSummaryToOgen(run entity.ProjectRun) ogenapi.ProjectRunSummary {
	return ogenapi.ProjectRunSummary{
		ID:        ogenapi.NewOptInt(run.ID),
		CreatedAt: ogenapi.NewOptDateTime(run.CreatedAt),
		UpdatedAt: ogenapi.NewOptDateTime(run.UpdatedAt),
		ProjectID: ogenapi.NewOptInt(run.ProjectID),
		Status:    ogenapi.NewOptProjectRunSummaryStatus(ogenapi.ProjectRunSummaryStatus(run.Status)),
		FailLog:   ogenapi.NewOptString(run.FailLog),
		Duration:  ogenapi.NewOptInt64(run.Duration.Milliseconds()),
	}
}
