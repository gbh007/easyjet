package mcp

import (
	"context"

	"github.com/gbh007/easyjet/internal/adapter/handler/mcp/internal"
	"github.com/mark3labs/mcp-go/mcp"
)

func (s *MCPServer) registerGetProject() {
	s.server.AddTool(
		mcp.NewTool("get_project",
			mcp.WithDescription("Get details of a specific project by ID"),
			mcp.WithNumber("id", mcp.Required(), mcp.Description("Project ID")),
		),
		func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			id := request.GetInt("id", 0)
			if id < 1 {
				return mcp.NewToolResultError("Invalid or missing project ID"), nil
			}

			project, err := s.service.Project(ctx, id)
			if err != nil {
				return mcp.NewToolResultError(err.Error()), nil
			}

			return mcp.NewToolResultJSON(internal.ToProjectResponse(project))
		},
	)
}
