package httpapi

import (
	"context"

	"github.com/gbh007/easyjet/internal/adapter/handler/httpapi/ogenapi"
	"github.com/gbh007/easyjet/internal/core/entity"
)

func (h *Handler) GetGlobalEnvVars(ctx context.Context) (*ogenapi.GetGlobalEnvVarsOK, error) {
	envVars, err := h.service.GlobalEnvVars(ctx)
	if err != nil {
		return nil, err
	}

	result := make([]ogenapi.EnvironmentVariable, 0, len(envVars))
	for _, ev := range envVars {
		result = append(result, convertEnvVarToOgen(ev))
	}

	return &ogenapi.GetGlobalEnvVarsOK{
		EnvVars: ogenapi.NewOptNilEnvironmentVariableArray(result),
	}, nil
}

func (h *Handler) GetGlobalEnvVar(ctx context.Context, params ogenapi.GetGlobalEnvVarParams) (*ogenapi.EnvironmentVariable, error) {
	ev, err := h.service.GlobalEnvVar(ctx, params.EnvVarID)
	if err != nil {
		return nil, err
	}

	result := convertEnvVarToOgen(ev)
	return &result, nil
}

func (h *Handler) CreateGlobalEnvVar(ctx context.Context, req *ogenapi.EnvironmentVariableCreate) (*ogenapi.CreateGlobalEnvVarCreated, error) {
	ev := entity.EnvironmentVariable{
		Name:               req.Name,
		Value:              req.Value,
		UsesOtherVariables: req.UsesOtherVariables.Or(false),
	}

	id, err := h.service.CreateGlobalEnvVar(ctx, ev)
	if err != nil {
		return nil, err
	}

	return &ogenapi.CreateGlobalEnvVarCreated{
		ID: ogenapi.NewOptInt(id),
	}, nil
}

func (h *Handler) UpdateGlobalEnvVar(ctx context.Context, req *ogenapi.EnvironmentVariableUpdate, params ogenapi.UpdateGlobalEnvVarParams) error {
	ev := entity.EnvironmentVariable{
		ID:                 params.EnvVarID,
		Name:               req.Name,
		Value:              req.Value,
		UsesOtherVariables: req.UsesOtherVariables.Or(false),
	}

	return h.service.UpdateGlobalEnvVar(ctx, ev)
}

func (h *Handler) DeleteGlobalEnvVar(ctx context.Context, params ogenapi.DeleteGlobalEnvVarParams) error {
	return h.service.DeleteGlobalEnvVar(ctx, params.EnvVarID)
}

func convertEnvVarToOgen(ev entity.EnvironmentVariable) ogenapi.EnvironmentVariable {
	result := ogenapi.EnvironmentVariable{
		ID:                 ogenapi.NewOptInt(ev.ID),
		CreatedAt:          ogenapi.NewOptDateTime(ev.CreatedAt),
		UpdatedAt:          ogenapi.NewOptDateTime(ev.UpdatedAt),
		Name:               ev.Name,
		Value:              ev.Value,
		UsesOtherVariables: ogenapi.NewOptBool(ev.UsesOtherVariables),
	}

	if ev.ProjectID != nil {
		result.ProjectID = ogenapi.NewOptNilInt(*ev.ProjectID)
	}

	return result
}
