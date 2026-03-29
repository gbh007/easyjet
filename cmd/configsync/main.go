//revive:disable:unhandled-error
package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"

	ht "github.com/ogen-go/ogen/http"

	"github.com/gbh007/easyjet/internal/adapter/handler/httpapi/ogenapi"
)

const (
	defaultServerURL = "http://localhost:8080"
	defaultDumpFile  = "dump.json"
	defaultTimeout   = 1 * time.Minute
)

type dumpData struct {
	Projects      []ogenapi.Project             `json:"projects"`
	GlobalEnvVars []ogenapi.EnvironmentVariable `json:"global_env_vars"`
}

func main() {
	if len(os.Args) < 2 {
		printUsage()
		os.Exit(1)
	}

	switch os.Args[1] {
	case "export":
		if err := runExport(os.Args[2:]); err != nil {
			fmt.Fprintf(os.Stderr, "Export error: %v\n", err)
			os.Exit(1)
		}
	case "import":
		if err := runImport(os.Args[2:]); err != nil {
			fmt.Fprintf(os.Stderr, "Import error: %v\n", err)
			os.Exit(1)
		}
	case "-h", "-help", "--help", "help":
		printUsage()
	default:
		fmt.Fprintf(os.Stderr, "Unknown command: %s\n", os.Args[1])
		printUsage()
		os.Exit(1)
	}
}

func printUsage() {
	fmt.Println("Usage: configsync <command> [options]")
	fmt.Println()
	fmt.Println("Commands:")
	fmt.Println("  export    Export project configurations and environment variables to dump.json")
	fmt.Println("  import    Import project configurations and environment variables from dump.json")
	fmt.Println()
	fmt.Println("Use 'configsync <command> -h' for more information about a command.")
}

func runExport(args []string) error {
	fs := flag.NewFlagSet("export", flag.ExitOnError)
	serverURL := fs.String("url", defaultServerURL, "Server URL")
	username := fs.String("username", "", "Username for basic auth")
	password := fs.String("password", "", "Password for basic auth")
	file := fs.String("file", defaultDumpFile, "Output file path")
	timeout := fs.Duration("timeout", defaultTimeout, "Request timeout")

	if err := fs.Parse(args); err != nil {
		return fmt.Errorf("parse flags: %w", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), *timeout)
	defer cancel()

	client, err := createClient(*serverURL, *username, *password)
	if err != nil {
		return fmt.Errorf("create client: %w", err)
	}

	projects, err := fetchProjects(ctx, client)
	if err != nil {
		return fmt.Errorf("fetch projects: %w", err)
	}

	globalEnvVars, err := fetchGlobalEnvVars(ctx, client)
	if err != nil {
		return fmt.Errorf("fetch global env vars: %w", err)
	}

	data := dumpData{
		Projects:      projects,
		GlobalEnvVars: globalEnvVars,
	}

	jsonData, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		return fmt.Errorf("marshal JSON: %w", err)
	}

	if err := os.WriteFile(*file, jsonData, 0o644); err != nil {
		return fmt.Errorf("write file: %w", err)
	}

	fmt.Printf("Successfully exported %d project(s) and %d global env var(s) to %s\n", len(projects), len(globalEnvVars), *file)
	return nil
}

func runImport(args []string) error {
	fs := flag.NewFlagSet("import", flag.ExitOnError)
	serverURL := fs.String("url", defaultServerURL, "Server URL")
	username := fs.String("username", "", "Username for basic auth")
	password := fs.String("password", "", "Password for basic auth")
	file := fs.String("file", defaultDumpFile, "Input file path")
	timeout := fs.Duration("timeout", defaultTimeout, "Request timeout")

	if err := fs.Parse(args); err != nil {
		return fmt.Errorf("parse flags: %w", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), *timeout)
	defer cancel()

	client, err := createClient(*serverURL, *username, *password)
	if err != nil {
		return fmt.Errorf("create client: %w", err)
	}

	data, err := loadDumpFile(*file)
	if err != nil {
		return fmt.Errorf("load dump file: %w", err)
	}

	for _, project := range data.Projects {
		_, err = client.CreateProject(ctx, &ogenapi.ProjectCreate{
			CronEnabled:    project.CronEnabled,
			CronSchedule:   project.CronSchedule,
			Dir:            project.Dir,
			GitBranch:      project.GitBranch,
			GitURL:         project.GitURL,
			Name:           project.Name,
			RestartAfter:   project.RestartAfter,
			RetentionCount: project.RetentionCount,
			WithRootEnv:    project.WithRootEnv,
			Stages:         project.Stages,
			EnvVars:        project.EnvVars,
		})
		if err != nil {
			return fmt.Errorf("create project %s: %w", project.Name, err)
		}
	}

	for _, envVar := range data.GlobalEnvVars {
		_, err = client.CreateGlobalEnvVar(ctx, &ogenapi.EnvironmentVariableCreate{
			Name:               envVar.Name,
			Value:              envVar.Value,
			UsesOtherVariables: envVar.UsesOtherVariables,
		})
		if err != nil {
			return fmt.Errorf("create global env var %s: %w", envVar.Name, err)
		}
	}

	return nil
}

func createClient(serverURL, username, password string) (*ogenapi.Client, error) {
	var httpClient ht.Client
	if username != "" && password != "" {
		httpClient = &http.Client{
			Transport: &basicAuthTransport{
				username: username,
				password: password,
				next:     http.DefaultTransport,
			},
		}
	}

	client, err := ogenapi.NewClient(serverURL, ogenapi.WithClient(httpClient))
	if err != nil {
		return nil, fmt.Errorf("create ogen client: %w", err)
	}

	return client, nil
}

func fetchProjects(ctx context.Context, client *ogenapi.Client) ([]ogenapi.Project, error) {
	resp, err := client.GetProjects(ctx)
	if err != nil {
		return nil, fmt.Errorf("get projects: %w", err)
	}

	projects := make([]ogenapi.Project, 0, len(resp.Projects.Value))

	for _, p := range resp.Projects.Value {
		resp, err := client.GetProject(ctx, ogenapi.GetProjectParams{
			ProjectID: p.ID,
		})
		if err != nil {
			return nil, fmt.Errorf("get project %d: %w", p.ID, err)
		}

		projects = append(projects, *resp)
	}

	return projects, nil
}

func fetchGlobalEnvVars(ctx context.Context, client *ogenapi.Client) ([]ogenapi.EnvironmentVariable, error) {
	resp, err := client.GetGlobalEnvVars(ctx)
	if err != nil {
		return nil, fmt.Errorf("get global env vars: %w", err)
	}

	return resp.EnvVars.Value, nil
}

func loadDumpFile(path string) (dumpData, error) {
	file, err := os.Open(path)
	if err != nil {
		return dumpData{}, fmt.Errorf("open file: %w", err)
	}
	defer func() { _ = file.Close() }()

	data, err := io.ReadAll(file)
	if err != nil {
		return dumpData{}, fmt.Errorf("read file: %w", err)
	}

	var dump dumpData
	if err := json.Unmarshal(data, &dump); err != nil {
		return dumpData{}, fmt.Errorf("unmarshal JSON: %w", err)
	}

	return dump, nil
}

type basicAuthTransport struct {
	username string
	password string
	next     http.RoundTripper
}

func (t *basicAuthTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	if t.username != "" && t.password != "" {
		req = req.Clone(req.Context())
		req.SetBasicAuth(t.username, t.password)
	}

	return t.next.RoundTrip(req)
}
