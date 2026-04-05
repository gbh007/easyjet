package mcp

import (
	"errors"

	"github.com/gbh007/easyjet/internal/core/entity"
	"github.com/mark3labs/mcp-go/mcp"
)

var projectCommon = []mcp.ToolOption{
	mcp.WithString("name", mcp.Required(), mcp.Description("Project name")),
	mcp.WithArray("stages",
		mcp.Required(),
		mcp.Description("Project stages"),
		mcp.Items(map[string]any{
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
		}),
	),
	mcp.WithString("dir", mcp.Description("Project directory. If empty, directory will be auto-created")),
	mcp.WithArray("env_vars",
		mcp.Description("Project environment variables"),
		mcp.Items(map[string]any{
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
		}),
	),
	mcp.WithString("git_url", mcp.Description("Git repository URL")),
	mcp.WithString("git_branch", mcp.Description("Git branch name")),
	mcp.WithBoolean("cron_enabled", mcp.Description("Enable cron scheduling")),
	mcp.WithString("cron_schedule", mcp.Description("Cron expression")),
	mcp.WithBoolean("restart_after", mcp.Description("Restart after completion")),
	mcp.WithNumber("retention_count", mcp.Description("Number of runs to retain")),
	mcp.WithBoolean("with_root_env", mcp.Description("Use root environment variables")),
	mcp.WithBoolean("is_template", mcp.Description("Mark as template")),
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

func parseProject(request mcp.CallToolRequest) (entity.Project, error) {
	args := request.GetArguments()

	stages := parseStages(args)
	if len(stages) == 0 {
		return entity.Project{}, errors.New("at least one stage with a non-empty script is required")
	}

	envVars := parseEnvVars(args)

	return entity.Project{
		Name:           request.GetString("name", ""),
		Dir:            request.GetString("dir", ""),
		GitURL:         request.GetString("git_url", ""),
		GitBranch:      request.GetString("git_branch", ""),
		CronEnabled:    request.GetBool("cron_enabled", false),
		CronSchedule:   request.GetString("cron_schedule", ""),
		RestartAfter:   request.GetBool("restart_after", false),
		RetentionCount: request.GetInt("retention_count", 0),
		WithRootEnv:    request.GetBool("with_root_env", false),
		IsTemplate:     request.GetBool("is_template", false),
		Stages:         stages,
		EnvVars:        envVars,
	}, nil
}
