package mcp

import (
	"context"
	"fmt"
	"slices"

	"github.com/gbh007/easyjet/internal/core/entity"
	"github.com/mark3labs/mcp-go/mcp"
)

func (s *MCPServer) registerProjectTools() {
	listProjectsTool := mcp.NewTool("list_projects",
		mcp.WithDescription("Get a list of all projects"),
	)

	s.server.AddTool(listProjectsTool, s.handleListProjects)

	getProjectTool := mcp.NewTool("get_project",
		mcp.WithDescription("Get details of a specific project by ID"),
		mcp.WithNumber("id", mcp.Required(), mcp.Description("Project ID")),
	)

	s.server.AddTool(getProjectTool, s.handleGetProject)
}

type ListProjectsResponse struct {
	Success  bool              `json:"success"`
	Count    int               `json:"count"`
	Projects []ProjectResponse `json:"projects"`
}

func (s *MCPServer) handleListProjects(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	projects, err := s.service.Projects(ctx)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to list projects: %v", err)), nil
	}

	responses := make([]ProjectResponse, 0, len(projects))
	for _, p := range projects {
		responses = append(responses, toProjectResponse(p))
	}

	resp := ListProjectsResponse{
		Success:  true,
		Count:    len(responses),
		Projects: responses,
	}

	return mcp.NewToolResultJSON(resp)
}

type GetProjectResponse struct {
	Success bool            `json:"success"`
	Project ProjectResponse `json:"project"`
}

func (s *MCPServer) handleGetProject(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	id := uint(request.GetFloat("id", 0))
	if id == 0 {
		return mcp.NewToolResultError("Invalid or missing project ID"), nil
	}

	project, err := s.service.Project(ctx, id)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to get project: %v", err)), nil
	}

	slices.SortFunc(project.EnvVars, func(a, b entity.EnvironmentVariable) int {
		return int(a.ID) - int(b.ID)
	})

	resp := GetProjectResponse{
		Success: true,
		Project: toProjectResponse(project),
	}

	return mcp.NewToolResultJSON(resp)
}
