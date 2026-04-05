package mcp

import (
	"context"

	"github.com/gbh007/easyjet/internal/adapter/handler/mcp/internal"
	"github.com/mark3labs/mcp-go/mcp"
)

func (s *MCPServer) registerListProjects() {
	s.server.AddTool(
		mcp.NewTool("list_projects",
			mcp.WithDescription("Get a list of all projects"),
		),
		func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			projects, err := s.service.Projects(ctx)
			if err != nil {
				return mcp.NewToolResultError(err.Error()), nil
			}

			response := make([]internal.ProjectResponse, 0, len(projects))
			for _, p := range projects {
				response = append(response, internal.ToProjectResponse(p))
			}

			return mcp.NewToolResultJSON(response)
		},
	)
}
