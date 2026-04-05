package mcp

import (
	"context"

	"github.com/gbh007/easyjet/internal/adapter/handler/mcp/internal"
	"github.com/mark3labs/mcp-go/mcp"
)

func (s *MCPServer) registerGetProjectRuns() {
	s.server.AddTool(
		mcp.NewTool("get_project_runs",
			mcp.WithDescription("Get run history for a specific project"),
			mcp.WithNumber("id", mcp.Required(), mcp.Description("Project ID")),
		),
		func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			id := request.GetInt("id", 0)
			if id < 1 {
				return mcp.NewToolResultError("Invalid or missing project ID"), nil
			}

			runs, err := s.service.ProjectRuns(ctx, id)
			if err != nil {
				return mcp.NewToolResultError(err.Error()), nil
			}

			response := make([]internal.ProjectRunResponse, 0, len(runs))
			for _, run := range runs {
				response = append(response, internal.ToProjectRunResponse(run))
			}

			return mcp.NewToolResultJSON(response)
		},
	)
}
