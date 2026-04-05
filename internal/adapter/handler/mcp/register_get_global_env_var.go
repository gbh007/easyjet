package mcp

import (
	"context"

	"github.com/gbh007/easyjet/internal/adapter/handler/mcp/internal"
	"github.com/mark3labs/mcp-go/mcp"
)

func (s *MCPServer) registerGetGlobalEnvVar() {
	s.server.AddTool(
		mcp.NewTool("get_global_env_var",
			mcp.WithDescription("Get details of a specific global environment variable"),
			mcp.WithNumber("id", mcp.Required(), mcp.Description("Environment variable ID")),
		),
		func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			id := request.GetInt("id", 0)
			if id < 1 {
				return mcp.NewToolResultError("Invalid or missing env var ID"), nil
			}

			envVar, err := s.service.GlobalEnvVar(ctx, id)
			if err != nil {
				return mcp.NewToolResultError(err.Error()), nil
			}

			return mcp.NewToolResultJSON(internal.ToEnvVarResponse(envVar))
		},
	)
}
