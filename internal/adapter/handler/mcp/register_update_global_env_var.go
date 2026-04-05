package mcp

import (
	"context"
	"time"

	"github.com/gbh007/easyjet/internal/core/entity"
	"github.com/mark3labs/mcp-go/mcp"
)

func (s *MCPServer) registerUpdateGlobalEnvVar() {
	s.server.AddTool(
		mcp.NewTool("update_global_env_var",
			mcp.WithDescription("Update an existing global environment variable"),
			mcp.WithNumber("id", mcp.Required(), mcp.Description("Environment variable ID")),
			mcp.WithString("name", mcp.Required(), mcp.Description("Variable name")),
			mcp.WithString("value", mcp.Required(), mcp.Description("Variable value")),
			mcp.WithBoolean("uses_other_variables", mcp.Description("Whether this variable references other variables")),
		),
		func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			id := request.GetInt("id", 0)
			if id < 1 {
				return mcp.NewToolResultError("Invalid or missing env var ID"), nil
			}

			envVar := entity.EnvironmentVariable{
				ID:                 id,
				Name:               request.GetString("name", ""),
				Value:              request.GetString("value", ""),
				UsesOtherVariables: request.GetBool("uses_other_variables", false),
				UpdatedAt:          time.Now(),
			}

			err := s.service.UpdateGlobalEnvVar(ctx, envVar)
			if err != nil {
				return mcp.NewToolResultError(err.Error()), nil
			}

			return mcp.NewToolResultJSON(map[string]any{
				"success": true,
			})
		},
	)
}
