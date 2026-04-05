package mcp

import (
	"context"

	"github.com/mark3labs/mcp-go/mcp"
)

func (s *MCPServer) registerRunProject() {
	s.server.AddTool(
		mcp.NewTool("run_project",
			mcp.WithDescription("Trigger a run of a specific project by ID"),
			mcp.WithNumber("id", mcp.Required(), mcp.Description("Project ID to run")),
		),
		func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			id := request.GetInt("id", 0)
			if id < 1 {
				return mcp.NewToolResultError("Invalid or missing project ID"), nil
			}

			runID, err := s.service.RunProject(ctx, id)
			if err != nil {
				return mcp.NewToolResultError(err.Error()), nil
			}

			return mcp.NewToolResultJSON(map[string]any{
				"run_id": runID,
			})
		},
	)
}
