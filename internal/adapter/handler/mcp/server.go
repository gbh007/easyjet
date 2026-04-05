package mcp

import (
	"log/slog"
	"net/http"

	"github.com/gbh007/easyjet/internal/core/port"
	"github.com/mark3labs/mcp-go/server"
)

type Config struct {
	Enabled        bool
	AllowRuns      bool
	AllowMutations bool
}

type MCPServer struct {
	logger  *slog.Logger
	cfg     Config
	service port.Service
	server  *server.MCPServer
	mux     *server.StreamableHTTPServer
}

func New(logger *slog.Logger, cfg Config, service port.Service) *MCPServer {
	s := &MCPServer{
		logger:  logger,
		cfg:     cfg,
		service: service,
	}

	s.server = server.NewMCPServer(
		"easyjet",
		"1.0.0",
		server.WithInstructions(`EasyJet is a lightweight CI/CD system designed for home servers and pet projects.
It serves as a minimal alternative to Jenkins, optimized for low resource consumption.

Key capabilities:
- Project management: Create, view, and update CI/CD pipeline configurations
- Pipeline execution: Run multi-stage build/deploy scripts with git integration
- Environment variables: Manage global environment variables and secrets for scripts
- Run history: Track execution results with logs for each pipeline stage
- Docker integration: Build Docker images and run Docker Compose configurations via shell scripts
- Git integration: Track commit changes and display git logs
- Scheduling: Support cron-based automated pipeline execution
- Self-updating: System can update itself via daemon restart

MCP tools available:
- List/get projects: View all configured CI/CD projects and their details
- List/get environment variables: Manage shared variables and secrets
- Get project runs: View execution history for a specific project
- Get run details: View detailed results of a specific pipeline execution
- Run project: Trigger manual pipeline execution (when enabled)
- Create/update/delete projects and variables: Full CRUD operations (when mutations enabled)
`),
	)

	s.registerTools()

	s.mux = server.NewStreamableHTTPServer(s.server)

	return s
}

func (s *MCPServer) registerTools() {
	s.registerListProjects()
	s.registerGetProject()

	s.registerListGlobalEnvVars()
	s.registerGetGlobalEnvVar()

	s.registerGetProjectRuns()
	s.registerGetRun()

	if s.cfg.AllowRuns {
		s.registerRunProject()
	}

	if s.cfg.AllowMutations {
		s.registerCreateProject()
		s.registerUpdateProject()

		s.registerCreateGlobalEnvVar()
		s.registerUpdateGlobalEnvVar()
		s.registerDeleteGlobalEnvVar()
	}
}

func (s *MCPServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.mux.ServeHTTP(w, r)
}
