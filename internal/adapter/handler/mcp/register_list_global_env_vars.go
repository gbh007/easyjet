package mcp

import (
	"context"

	"github.com/gbh007/easyjet/internal/adapter/handler/mcp/internal"
	"github.com/mark3labs/mcp-go/mcp"
)

func (s *MCPServer) registerListGlobalEnvVars() {
	s.server.AddTool(
		mcp.NewTool("list_global_env_vars",
			mcp.WithDescription("Get a list of all global environment variables"),
		),
		func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			envVars, err := s.service.GlobalEnvVars(ctx)
			if err != nil {
				return mcp.NewToolResultError(err.Error()), nil
			}

			return mcp.NewToolResultJSON(internal.ToEnvVarsResponse(envVars))
		},
	)
}
