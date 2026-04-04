package mcp

import (
	"context"
	"fmt"

	"github.com/gbh007/easyjet/internal/core/entity"
	"github.com/mark3labs/mcp-go/mcp"
)

func (s *MCPServer) registerMutationTools() {
	stageSchema := map[string]any{
		"type": "object",
		"properties": map[string]any{
			"number": map[string]any{
				"type":        "integer",
				"description": "Stage number (will be auto-assigned by server)",
			},
			"script": map[string]any{
				"type":        "string",
				"description": "Shell script content for this stage",
			},
		},
		"required": []string{"script"},
	}
	varSchema := map[string]any{
		"type": "object",
		"properties": map[string]any{
			"name": map[string]any{
				"type":        "string",
				"description": "Environment variable name",
			},
			"value": map[string]any{
				"type":        "string",
				"description": "Environment variable value",
			},
			"uses_other_variables": map[string]any{
				"type":        "boolean",
				"description": "Whether this variable references other variables",
			},
		},
		"required": []string{"name", "value"},
	}

	createProjectTool := mcp.NewTool("create_project",
		mcp.WithDescription("Create a new project"),
		mcp.WithString("name", mcp.Required(), mcp.Description("Project name")),
		mcp.WithArray("stages",
			mcp.Required(),
			mcp.Description("Project stages"),
			mcp.Items(stageSchema),
		),
		mcp.WithString("dir", mcp.Description("Project directory. If empty, directory will be auto-created")),
		mcp.WithArray("env_vars",
			mcp.Description("Project environment variables"),
			mcp.Items(varSchema),
		),
		mcp.WithString("git_url", mcp.Description("Git repository URL")),
		mcp.WithString("git_branch", mcp.Description("Git branch name")),
		mcp.WithBoolean("cron_enabled", mcp.Description("Enable cron scheduling")),
		mcp.WithString("cron_schedule", mcp.Description("Cron expression")),
		mcp.WithBoolean("restart_after", mcp.Description("Restart after completion")),
		mcp.WithNumber("retention_count", mcp.Description("Number of runs to retain")),
		mcp.WithBoolean("with_root_env", mcp.Description("Use root environment variables")),
		mcp.WithBoolean("is_template", mcp.Description("Mark as template")),
	)

	s.server.AddTool(createProjectTool, s.handleCreateProject)

	updateProjectTool := mcp.NewTool("update_project",
		mcp.WithDescription("Update an existing project. This is a full replacement - stages and env_vars are completely replaced"),
		mcp.WithNumber("id", mcp.Required(), mcp.Description("Project ID")),
		mcp.WithString("name", mcp.Required(), mcp.Description("Project name")),
		mcp.WithArray("stages",
			mcp.Required(),
			mcp.Description("Project stages (replaces all existing stages)"),
			mcp.Items(stageSchema),
		),
		mcp.WithString("dir", mcp.Description("Project directory. If empty, directory will be auto-created")),
		mcp.WithArray("env_vars",
			mcp.Description("Project environment variables (replaces all existing env vars)"),
			mcp.Items(varSchema),
		),
		mcp.WithString("git_url", mcp.Description("Git repository URL")),
		mcp.WithString("git_branch", mcp.Description("Git branch name")),
		mcp.WithBoolean("cron_enabled", mcp.Description("Enable cron scheduling")),
		mcp.WithString("cron_schedule", mcp.Description("Cron expression")),
		mcp.WithBoolean("restart_after", mcp.Description("Restart after completion")),
		mcp.WithNumber("retention_count", mcp.Description("Number of runs to retain")),
		mcp.WithBoolean("with_root_env", mcp.Description("Use root environment variables")),
		mcp.WithBoolean("is_template", mcp.Description("Mark as template")),
	)

	s.server.AddTool(updateProjectTool, s.handleUpdateProject)

	createEnvVarTool := mcp.NewTool("create_global_env_var",
		mcp.WithDescription("Create a new global environment variable"),
		mcp.WithString("name", mcp.Required(), mcp.Description("Variable name")),
		mcp.WithString("value", mcp.Required(), mcp.Description("Variable value")),
		mcp.WithBoolean("uses_other_variables", mcp.Description("Whether this variable references other variables")),
	)

	s.server.AddTool(createEnvVarTool, s.handleCreateGlobalEnvVar)

	updateEnvVarTool := mcp.NewTool("update_global_env_var",
		mcp.WithDescription("Update an existing global environment variable"),
		mcp.WithNumber("id", mcp.Required(), mcp.Description("Environment variable ID")),
		mcp.WithString("name", mcp.Required(), mcp.Description("Variable name")),
		mcp.WithString("value", mcp.Required(), mcp.Description("Variable value")),
		mcp.WithBoolean("uses_other_variables", mcp.Description("Whether this variable references other variables")),
	)

	s.server.AddTool(updateEnvVarTool, s.handleUpdateGlobalEnvVar)

	// Delete global env var tool
	deleteEnvVarTool := mcp.NewTool("delete_global_env_var",
		mcp.WithDescription("Delete a global environment variable"),
		mcp.WithNumber("id", mcp.Required(), mcp.Description("Environment variable ID")),
	)

	s.server.AddTool(deleteEnvVarTool, s.handleDeleteGlobalEnvVar)
}

func parseStages(args map[string]any) []entity.ProjectStage {
	stagesRaw, ok := args["stages"].([]any)
	if !ok {
		return nil
	}

	stages := make([]entity.ProjectStage, 0, len(stagesRaw))
	for i, stageRaw := range stagesRaw {
		stageMap, ok := stageRaw.(map[string]any)
		if !ok {
			continue
		}
		script, ok := stageMap["script"].(string)
		if !ok || script == "" {
			continue
		}
		stages = append(stages, entity.ProjectStage{
			Number: i + 1,
			Script: script,
		})
	}

	return stages
}

func parseEnvVars(args map[string]any) []entity.EnvironmentVariable {
	envVarsRaw, ok := args["env_vars"].([]any)
	if !ok {
		return nil
	}

	envVars := make([]entity.EnvironmentVariable, 0, len(envVarsRaw))
	for _, envVarRaw := range envVarsRaw {
		envVarMap, ok := envVarRaw.(map[string]any)
		if !ok {
			continue
		}
		name, nameOk := envVarMap["name"].(string)
		value, valueOk := envVarMap["value"].(string)
		if !nameOk || !valueOk {
			continue
		}
		usesOtherVars := false
		if v, ok := envVarMap["uses_other_variables"].(bool); ok {
			usesOtherVars = v
		}
		envVars = append(envVars, entity.EnvironmentVariable{
			Name:               name,
			Value:              value,
			UsesOtherVariables: usesOtherVars,
		})
	}

	return envVars
}

type CreateProjectResponse struct {
	Success bool   `json:"success"`
	ID      uint   `json:"id"`
	Message string `json:"message"`
}

func (s *MCPServer) handleCreateProject(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	args := request.GetArguments()

	stages := parseStages(args)
	if len(stages) == 0 {
		return mcp.NewToolResultError("At least one stage with a non-empty script is required"), nil
	}

	envVars := parseEnvVars(args)

	project := entity.Project{
		Name:           request.GetString("name", ""),
		Dir:            request.GetString("dir", ""),
		GitURL:         request.GetString("git_url", ""),
		GitBranch:      request.GetString("git_branch", ""),
		CronEnabled:    request.GetBool("cron_enabled", false),
		CronSchedule:   request.GetString("cron_schedule", ""),
		RestartAfter:   request.GetBool("restart_after", false),
		RetentionCount: request.GetInt("retention_count", 10),
		WithRootEnv:    request.GetBool("with_root_env", false),
		IsTemplate:     request.GetBool("is_template", false),
		Stages:         stages,
		EnvVars:        envVars,
	}

	id, err := s.service.CreateProject(ctx, project)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to create project: %v", err)), nil
	}

	resp := CreateProjectResponse{
		Success: true,
		ID:      id,
		Message: "Project created successfully",
	}

	return mcp.NewToolResultJSON(resp)
}

type UpdateProjectResponse struct {
	Success bool   `json:"success"`
	ID      uint   `json:"id"`
	Message string `json:"message"`
}

func (s *MCPServer) handleUpdateProject(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	args := request.GetArguments()

	id := uint(request.GetFloat("id", 0))
	if id == 0 {
		return mcp.NewToolResultError("Invalid or missing project ID"), nil
	}

	existingProject, err := s.service.Project(ctx, id)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to get project: %v", err)), nil
	}

	stages := parseStages(args)
	if len(stages) == 0 {
		return mcp.NewToolResultError("At least one stage with a non-empty script is required"), nil
	}

	envVars := parseEnvVars(args)

	project := entity.Project{
		ID:             id,
		CreatedAt:      existingProject.CreatedAt,
		Name:           request.GetString("name", ""),
		Dir:            request.GetString("dir", ""),
		GitURL:         request.GetString("git_url", ""),
		GitBranch:      request.GetString("git_branch", ""),
		CronEnabled:    request.GetBool("cron_enabled", false),
		CronSchedule:   request.GetString("cron_schedule", ""),
		RestartAfter:   request.GetBool("restart_after", false),
		RetentionCount: request.GetInt("retention_count", 10),
		WithRootEnv:    request.GetBool("with_root_env", false),
		IsTemplate:     request.GetBool("is_template", false),
		Stages:         stages,
		EnvVars:        envVars,
	}

	err = s.service.UpdateProject(ctx, project)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to update project: %v", err)), nil
	}

	resp := UpdateProjectResponse{
		Success: true,
		ID:      id,
		Message: "Project updated successfully",
	}

	return mcp.NewToolResultJSON(resp)
}

type CreateEnvVarResponse struct {
	Success bool   `json:"success"`
	ID      uint   `json:"id"`
	Message string `json:"message"`
}

func (s *MCPServer) handleCreateGlobalEnvVar(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	envVar := entity.EnvironmentVariable{
		Name:               request.GetString("name", ""),
		Value:              request.GetString("value", ""),
		UsesOtherVariables: request.GetBool("uses_other_variables", false),
	}

	id, err := s.service.CreateGlobalEnvVar(ctx, envVar)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to create env var: %v", err)), nil
	}

	resp := CreateEnvVarResponse{
		Success: true,
		ID:      id,
		Message: "Environment variable created successfully",
	}

	return mcp.NewToolResultJSON(resp)
}

type UpdateEnvVarResponse struct {
	Success bool   `json:"success"`
	ID      uint   `json:"id"`
	Message string `json:"message"`
}

func (s *MCPServer) handleUpdateGlobalEnvVar(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	id := uint(request.GetFloat("id", 0))
	if id == 0 {
		return mcp.NewToolResultError("Invalid or missing env var ID"), nil
	}

	// Get existing env var to preserve CreatedAt and other fields
	existingEnvVar, err := s.service.GlobalEnvVar(ctx, id)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to get env var: %v", err)), nil
	}

	envVar := entity.EnvironmentVariable{
		ID:                 id,
		CreatedAt:          existingEnvVar.CreatedAt,
		ProjectID:          existingEnvVar.ProjectID,
		Name:               request.GetString("name", ""),
		Value:              request.GetString("value", ""),
		UsesOtherVariables: request.GetBool("uses_other_variables", false),
	}

	err = s.service.UpdateGlobalEnvVar(ctx, envVar)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to update env var: %v", err)), nil
	}

	resp := UpdateEnvVarResponse{
		Success: true,
		ID:      id,
		Message: "Environment variable updated successfully",
	}

	return mcp.NewToolResultJSON(resp)
}

type DeleteEnvVarResponse struct {
	Success bool   `json:"success"`
	ID      uint   `json:"id"`
	Message string `json:"message"`
}

func (s *MCPServer) handleDeleteGlobalEnvVar(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	id := uint(request.GetFloat("id", 0))
	if id == 0 {
		return mcp.NewToolResultError("Invalid or missing env var ID"), nil
	}

	err := s.service.DeleteGlobalEnvVar(ctx, id)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to delete env var: %v", err)), nil
	}

	resp := DeleteEnvVarResponse{
		Success: true,
		ID:      id,
		Message: "Environment variable deleted successfully",
	}

	return mcp.NewToolResultJSON(resp)
}
