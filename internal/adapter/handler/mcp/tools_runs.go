package mcp

import (
	"context"
	"fmt"

	"github.com/mark3labs/mcp-go/mcp"
)

func (s *MCPServer) registerRunTools() {
	runProjectTool := mcp.NewTool("run_project",
		mcp.WithDescription("Trigger a run of a specific project by ID"),
		mcp.WithNumber("id", mcp.Required(), mcp.Description("Project ID to run")),
	)

	s.server.AddTool(runProjectTool, s.handleRunProject)
}

type RunProjectResponse struct {
	Success bool   `json:"success"`
	RunID   int    `json:"run_id"`
	Message string `json:"message"`
}

func (s *MCPServer) handleRunProject(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	id := request.GetInt("id", 0)
	if id == 0 {
		return mcp.NewToolResultError("Invalid or missing project ID"), nil
	}

	runID, err := s.service.RunProject(ctx, id)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to run project: %v", err)), nil
	}

	resp := RunProjectResponse{
		Success: true,
		RunID:   runID,
		Message: "Project run started",
	}

	return mcp.NewToolResultJSON(resp)
}
