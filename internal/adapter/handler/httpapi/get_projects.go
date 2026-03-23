package httpapi

import (
	"context"

	"github.com/gbh007/easyjet/internal/adapter/handler/httpapi/ogenapi"
)

func (h *Handler) GetProjects(ctx context.Context) (*ogenapi.GetProjectsOK, error) {
	items, err := h.service.ProjectsWithRunInfo(ctx)
	if err != nil {
		return nil, err
	}

	return &ogenapi.GetProjectsOK{
		Projects: ogenapi.NewOptNilProjectListItemArray(convertProjectListItemsToOgen(items)),
	}, nil
}
