package mcp

import (
	"context"
	"time"

	"github.com/mark3labs/mcp-go/mcp"
)

func (s *MCPServer) registerCreateProject() {
	s.server.AddTool(
		mcp.NewTool("create_project",
			append([]mcp.ToolOption{
				mcp.WithDescription("Create a new project"),
			},
				projectCommon...,
			)...,
		),
		func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			project, err := parseProject(request)
			if err != nil {
				return mcp.NewToolResultError(err.Error()), nil
			}

			project.CreatedAt = time.Now()

			id, err := s.service.CreateProject(ctx, project)
			if err != nil {
				return mcp.NewToolResultError(err.Error()), nil
			}

			return mcp.NewToolResultJSON(map[string]any{
				"id": id,
			})
		},
	)
}
