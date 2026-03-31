package httpapi

import (
	"context"

	"github.com/gbh007/easyjet/internal/adapter/handler/httpapi/ogenapi"
)

func (h *Handler) GetProjects(ctx context.Context, params ogenapi.GetProjectsParams) (*ogenapi.GetProjectsOK, error) {
	filterType := "all"
	if params.Type.Set {
		filterType = string(params.Type.Value)
	}

	items, err := h.service.ProjectsWithRunInfo(ctx, filterType)
	if err != nil {
		return nil, err
	}

	return &ogenapi.GetProjectsOK{
		Projects: ogenapi.NewOptNilProjectListItemArray(convertProjectListItemsToOgen(items)),
	}, nil
}
