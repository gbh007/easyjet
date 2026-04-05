package mcp

import (
	"context"
	"time"

	"github.com/mark3labs/mcp-go/mcp"
)

func (s *MCPServer) registerUpdateProject() {
	s.server.AddTool(
		mcp.NewTool("update_project",
			append([]mcp.ToolOption{
				mcp.WithDescription("Update an existing project. This is a full replacement - stages and env_vars are completely replaced"),
				mcp.WithNumber("id", mcp.Required(), mcp.Description("Project ID")),
			},
				projectCommon...,
			)...,
		),
		func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			id := request.GetInt("id", 0)
			if id < 1 {
				return mcp.NewToolResultError("Invalid or missing project ID"), nil
			}

			project, err := parseProject(request)
			if err != nil {
				return mcp.NewToolResultError(err.Error()), nil
			}

			project.ID = id
			project.UpdatedAt = time.Now()

			err = s.service.UpdateProject(ctx, project)
			if err != nil {
				return mcp.NewToolResultError(err.Error()), nil
			}

			return mcp.NewToolResultJSON(map[string]any{
				"success": true,
			})
		},
	)
}
