package mcp

import (
	"context"
	"fmt"

	"github.com/mark3labs/mcp-go/mcp"
)

func (s *MCPServer) registerEnvVarTools() {
	// List global env vars tool
	listEnvVarsTool := mcp.NewTool("list_global_env_vars",
		mcp.WithDescription("Get a list of all global environment variables"),
	)

	s.server.AddTool(listEnvVarsTool, s.handleListGlobalEnvVars)

	// Get global env var tool
	getEnvVarTool := mcp.NewTool("get_global_env_var",
		mcp.WithDescription("Get details of a specific global environment variable"),
		mcp.WithNumber("id", mcp.Required(), mcp.Description("Environment variable ID")),
	)

	s.server.AddTool(getEnvVarTool, s.handleGetGlobalEnvVar)

	// Get project runs tool (read-only, always available)
	getProjectRunsTool := mcp.NewTool("get_project_runs",
		mcp.WithDescription("Get run history for a specific project"),
		mcp.WithNumber("id", mcp.Required(), mcp.Description("Project ID")),
	)

	s.server.AddTool(getProjectRunsTool, s.handleGetProjectRuns)

	// Get run details tool (read-only, always available)
	getRunTool := mcp.NewTool("get_run",
		mcp.WithDescription("Get detailed information about a specific run"),
		mcp.WithNumber("id", mcp.Required(), mcp.Description("Run ID")),
	)

	s.server.AddTool(getRunTool, s.handleGetRun)
}

type ListEnvVarsResponse struct {
	Success bool             `json:"success"`
	Count   int              `json:"count"`
	EnvVars []EnvVarResponse `json:"env_vars"`
}

func (s *MCPServer) handleListGlobalEnvVars(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	envVars, err := s.service.GlobalEnvVars(ctx)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to list global env vars: %v", err)), nil
	}

	responses := make([]EnvVarResponse, 0, len(envVars))
	for _, ev := range envVars {
		responses = append(responses, toEnvVarResponse(ev))
	}

	resp := ListEnvVarsResponse{
		Success: true,
		Count:   len(responses),
		EnvVars: responses,
	}

	return mcp.NewToolResultJSON(resp)
}

type GetEnvVarResponse struct {
	Success bool           `json:"success"`
	EnvVar  EnvVarResponse `json:"env_var"`
}

func (s *MCPServer) handleGetGlobalEnvVar(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	id := request.GetInt("id", 0)
	if id == 0 {
		return mcp.NewToolResultError("Invalid or missing env var ID"), nil
	}

	envVar, err := s.service.GlobalEnvVar(ctx, id)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to get env var: %v", err)), nil
	}

	resp := GetEnvVarResponse{
		Success: true,
		EnvVar:  toEnvVarResponse(envVar),
	}

	return mcp.NewToolResultJSON(resp)
}

type ListRunsResponse struct {
	Success bool                 `json:"success"`
	Count   int                  `json:"count"`
	Runs    []ProjectRunResponse `json:"runs"`
}

func (s *MCPServer) handleGetProjectRuns(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	id := request.GetInt("id", 0)
	if id == 0 {
		return mcp.NewToolResultError("Invalid or missing project ID"), nil
	}

	runs, err := s.service.ProjectRuns(ctx, id)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to get project runs: %v", err)), nil
	}

	responses := make([]ProjectRunResponse, 0, len(runs))
	for _, run := range runs {
		responses = append(responses, toProjectRunResponse(run))
	}

	resp := ListRunsResponse{
		Success: true,
		Count:   len(responses),
		Runs:    responses,
	}

	return mcp.NewToolResultJSON(resp)
}

type GetRunResponse struct {
	Success bool               `json:"success"`
	Run     ProjectRunResponse `json:"run"`
}

func (s *MCPServer) handleGetRun(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	id := request.GetInt("id", 0)
	if id == 0 {
		return mcp.NewToolResultError("Invalid or missing run ID"), nil
	}

	run, err := s.service.ProjectRun(ctx, id)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to get run: %v", err)), nil
	}

	resp := GetRunResponse{
		Success: true,
		Run:     toProjectRunResponse(run),
	}

	return mcp.NewToolResultJSON(resp)
}
