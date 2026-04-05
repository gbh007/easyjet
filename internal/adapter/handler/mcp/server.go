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
