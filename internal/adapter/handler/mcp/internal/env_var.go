package internal

import (
	"time"

	"github.com/gbh007/easyjet/internal/core/entity"
)

type EnvVarResponse struct {
	ID                 int       `json:"id"`
	Name               string    `json:"name"`
	Value              string    `json:"value"`
	UsesOtherVariables bool      `json:"uses_other_variables"`
	ProjectID          *int      `json:"project_id,omitempty"`
	CreatedAt          time.Time `json:"created_at,omitzero"`
	UpdatedAt          time.Time `json:"updated_at,omitzero"`
}

func ToEnvVarResponse(ev entity.EnvironmentVariable) EnvVarResponse {
	return EnvVarResponse{
		ID:                 ev.ID,
		Name:               ev.Name,
		Value:              ev.Value,
		UsesOtherVariables: ev.UsesOtherVariables,
		ProjectID:          ev.ProjectID,
		CreatedAt:          ev.CreatedAt,
		UpdatedAt:          ev.UpdatedAt,
	}
}

func ToEnvVarsResponse(envVars []entity.EnvironmentVariable) []EnvVarResponse {
	result := make([]EnvVarResponse, 0, len(envVars))
	for _, ev := range envVars {
		result = append(result, ToEnvVarResponse(ev))
	}
	return result
}
