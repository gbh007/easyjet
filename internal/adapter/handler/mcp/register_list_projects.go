package mcp

import (
	"context"

	"github.com/gbh007/easyjet/internal/adapter/handler/mcp/internal"
	"github.com/gbh007/easyjet/internal/core/entity"
	"github.com/mark3labs/mcp-go/mcp"
)

func (s *MCPServer) registerListProjects() {
	s.server.AddTool(
		mcp.NewTool("list_projects",
			mcp.WithDescription("Get a list of all projects"),
			mcp.WithString("filter_type",
				mcp.Description("Filter type for projects: all, project, or template"),
				mcp.Enum("all", "project", "template"),
				mcp.DefaultString("all"),
			),
		),
		func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			filterTypeStr := request.GetString("filter_type", "all")
			filterType := entity.ProjectFilterType(filterTypeStr)

			projects, err := s.service.ProjectsWithRunInfo(ctx, filterType)
			if err != nil {
				return mcp.NewToolResultError(err.Error()), nil
			}

			response := make([]internal.ProjectWithRunInfoResponse, 0, len(projects))
			for _, p := range projects {
				response = append(response, internal.ToProjectResponseFromInfo(p))
			}

			return mcp.NewToolResultJSON(map[string]any{
				"projects": response,
			})
		},
	)
}
