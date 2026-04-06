package mcp

import (
	"context"

	"github.com/mark3labs/mcp-go/mcp"
)

func (s *MCPServer) registerDeleteProject() {
	s.server.AddTool(
		mcp.NewTool("delete_project",
			mcp.WithDescription("Delete a project and all associated runs. This action cannot be undone."),
			mcp.WithNumber("id", mcp.Required(), mcp.Description("Project ID")),
		),
		func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			id := request.GetInt("id", 0)
			if id < 1 {
				return mcp.NewToolResultError("Invalid or missing project ID"), nil
			}

			err := s.service.DeleteProject(ctx, id)
			if err != nil {
				return mcp.NewToolResultError(err.Error()), nil
			}

			return mcp.NewToolResultJSON(map[string]any{
				"success": true,
			})
		},
	)
}
