package httpapi

import (
	"github.com/gbh007/easyjet/internal/adapter/handler/httpapi/ogenapi"
	"github.com/gbh007/easyjet/internal/core/entity"
)

func convertProjectListItemsToOgen(items []entity.ProjectsWithRunInfo) []ogenapi.ProjectListItem {
	result := make([]ogenapi.ProjectListItem, len(items))
	for i, item := range items {
		result[i] = convertProjectListItemToOgen(item)
	}
	return result
}

func convertProjectListItemToOgen(item entity.ProjectsWithRunInfo) ogenapi.ProjectListItem {
	result := ogenapi.ProjectListItem{
		ID:          item.ID,
		Name:        item.Name,
		CronEnabled: ogenapi.NewOptBool(item.CronEnabled),
		IsTemplate:  ogenapi.NewOptBool(item.IsTemplate),
	}

	if item.LastSuccessfulRunAt != nil {
		result.LastSuccessfulRunAt = ogenapi.NewOptNilDateTime(*item.LastSuccessfulRunAt)
	}

	if item.LastRun != nil {
		result.LastRun = ogenapi.NewOptProjectLastRun(convertProjectLastRunToOgen(*item.LastRun))
	}

	return result
}

func convertProjectLastRunToOgen(lastRun entity.ProjectLastRun) ogenapi.ProjectLastRun {
	return ogenapi.ProjectLastRun{
		CreatedAt:  ogenapi.NewOptDateTime(lastRun.CreatedAt),
		Success:    ogenapi.NewOptBool(lastRun.Success),
		Pending:    ogenapi.NewOptBool(lastRun.Pending),
		Processing: ogenapi.NewOptBool(lastRun.Processing),
		Duration:   ogenapi.NewOptInt64(lastRun.Duration.Milliseconds()),
	}
}
