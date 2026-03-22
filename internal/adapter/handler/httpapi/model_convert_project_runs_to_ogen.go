package httpapi

import (
	"github.com/gbh007/easyjet/internal/adapter/handler/httpapi/ogenapi"
	"github.com/gbh007/easyjet/internal/core/entity"
)

func convertProjectRunsToOgen(runs []entity.ProjectRun) []ogenapi.ProjectRun {
	result := make([]ogenapi.ProjectRun, len(runs))
	for i, run := range runs {
		result[i] = convertProjectRunToOgen(run)
	}
	return result
}
