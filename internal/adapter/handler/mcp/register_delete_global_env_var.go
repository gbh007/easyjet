package mcp

import (
	"context"

	"github.com/mark3labs/mcp-go/mcp"
)

func (s *MCPServer) registerDeleteGlobalEnvVar() {
	s.server.AddTool(
		mcp.NewTool("delete_global_env_var",
			mcp.WithDescription("Delete a global environment variable"),
			mcp.WithNumber("id", mcp.Required(), mcp.Description("Environment variable ID")),
		),
		func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			id := request.GetInt("id", 0)
			if id < 1 {
				return mcp.NewToolResultError("Invalid or missing env var ID"), nil
			}

			err := s.service.DeleteGlobalEnvVar(ctx, id)
			if err != nil {
				return mcp.NewToolResultError(err.Error()), nil
			}

			return mcp.NewToolResultJSON(map[string]any{
				"success": true,
			})
		},
	)
}
