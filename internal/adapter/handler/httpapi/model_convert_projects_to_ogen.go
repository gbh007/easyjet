package httpapi

import (
	"github.com/gbh007/easyjet/internal/adapter/handler/httpapi/ogenapi"
	"github.com/gbh007/easyjet/internal/core/entity"
)

func convertProjectsToOgen(projects []entity.Project) []ogenapi.Project {
	result := make([]ogenapi.Project, len(projects))
	for i, project := range projects {
		result[i] = convertProjectToOgen(project)
	}
	return result
}
