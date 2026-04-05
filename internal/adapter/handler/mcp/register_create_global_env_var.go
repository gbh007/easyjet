package mcp

import (
	"context"
	"time"

	"github.com/gbh007/easyjet/internal/core/entity"
	"github.com/mark3labs/mcp-go/mcp"
)

func (s *MCPServer) registerCreateGlobalEnvVar() {
	s.server.AddTool(
		mcp.NewTool("create_global_env_var",
			mcp.WithDescription("Create a new global environment variable"),
			mcp.WithString("name", mcp.Required(), mcp.Description("Variable name")),
			mcp.WithString("value", mcp.Required(), mcp.Description("Variable value")),
			mcp.WithBoolean("uses_other_variables", mcp.Description("Whether this variable references other variables")),
		),
		func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			envVar := entity.EnvironmentVariable{
				Name:               request.GetString("name", ""),
				Value:              request.GetString("value", ""),
				UsesOtherVariables: request.GetBool("uses_other_variables", false),
				CreatedAt:          time.Now(),
			}

			id, err := s.service.CreateGlobalEnvVar(ctx, envVar)
			if err != nil {
				return mcp.NewToolResultError(err.Error()), nil
			}

			return mcp.NewToolResultJSON(map[string]any{
				"id": id,
			})
		},
	)
}
