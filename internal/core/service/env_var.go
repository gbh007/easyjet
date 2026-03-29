package service

import (
	"context"
	"fmt"
	"os"
	"regexp"
	"strings"

	"github.com/gbh007/easyjet/internal/core/entity"
)

func (s Service) GlobalEnvVars(ctx context.Context) ([]entity.EnvironmentVariable, error) {
	return s.db.GlobalEnvVars(ctx)
}

func (s Service) GlobalEnvVar(ctx context.Context, id uint) (entity.EnvironmentVariable, error) {
	return s.db.GlobalEnvVar(ctx, id)
}

func (s Service) CreateGlobalEnvVar(ctx context.Context, ev entity.EnvironmentVariable) (uint, error) {
	return s.db.SetGlobalEnvVar(ctx, ev)
}

func (s Service) UpdateGlobalEnvVar(ctx context.Context, ev entity.EnvironmentVariable) error {
	_, err := s.db.SetGlobalEnvVar(ctx, ev)
	return err
}

func (s Service) DeleteGlobalEnvVar(ctx context.Context, id uint) error {
	return s.db.DeleteGlobalEnvVar(ctx, id)
}

// CalculateEffectiveEnvVars calculates the final environment variables for script execution
// Priority (lowest to highest):
// 1. Parent environment (if withRootEnv is true)
// 2. Static workspace variables (WORKSPACE)
// 3. Global variables
// 4. Project variables
func (s Service) CalculateEffectiveEnvVars(ctx context.Context, project entity.Project, dir string) ([]string, error) {
	var result []string

	envMap := make(map[string]string)

	if project.WithRootEnv {
		for _, raw := range os.Environ() {
			name, value, ok := strings.Cut(raw, "=")
			if ok {
				envMap[name] = value
			}
		}
	}

	envMap["WORKSPACE"] = dir

	globalVars, err := s.db.GlobalEnvVars(ctx)
	if err != nil {
		return nil, fmt.Errorf("get global env vars: %w", err)
	}

	for _, ev := range globalVars {
		value := ev.Value

		if ev.UsesOtherVariables {
			value = s.resolveVariables(envMap, value)
		}

		envMap[ev.Name] = value
	}

	for _, ev := range project.EnvVars {
		value := ev.Value

		if ev.UsesOtherVariables {
			value = s.resolveVariables(envMap, value)
		}

		envMap[ev.Name] = value
	}

	for name, value := range envMap {
		result = append(result, name+"="+value)
	}

	return result, nil
}

func (s Service) resolveVariables(envMap map[string]string, value string) string {
	varRegex := regexp.MustCompile(`\$(\w+)`)

	matches := varRegex.FindAllStringSubmatch(value, -1)
	for _, match := range matches {
		if len(match) == 2 {
			refName := match[1]
			if refValue, exists := envMap[refName]; exists {
				value = strings.ReplaceAll(value, match[0], refValue)
			}
		}
	}

	return value
}
