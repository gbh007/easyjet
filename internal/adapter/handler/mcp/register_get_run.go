package mcp

import (
	"context"

	"github.com/gbh007/easyjet/internal/adapter/handler/mcp/internal"
	"github.com/mark3labs/mcp-go/mcp"
)

func (s *MCPServer) registerGetRun() {
	s.server.AddTool(
		mcp.NewTool("get_run",
			mcp.WithDescription("Get detailed information about a specific run"),
			mcp.WithNumber("id", mcp.Required(), mcp.Description("Run ID")),
		),
		func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			id := request.GetInt("id", 0)
			if id < 1 {
				return mcp.NewToolResultError("Invalid or missing run ID"), nil
			}

			run, err := s.service.ProjectRun(ctx, id)
			if err != nil {
				return mcp.NewToolResultError(err.Error()), nil
			}

			return mcp.NewToolResultJSON(internal.ToProjectRunResponse(run))
		},
	)
}
