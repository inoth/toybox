package mcpserver

import (
	"context"
	"log"
	"net/http"

	"github.com/inoth/toybox/transport"
	"github.com/mark3labs/mcp-go/server"
	"github.com/pkg/errors"
)

const (
	name = "mcp"
)

var _ transport.Transport = (*MCPServer)(nil)

type MCPServer struct {
	MCPName string `toml:"name"`
	Version string `toml:"version"`
	Port    string `toml:"port"`

	resources []server.ServerResource
	tools     []server.ServerTool
	prompts   []server.ServerPrompt

	toolFilters            []server.ToolFilterFunc
	toolHandlerMiddlewares []server.ToolHandlerMiddleware

	mcp *server.MCPServer
	svr *server.StreamableHTTPServer
}

func NewMCPServer(options ...Option) *MCPServer {
	s := &MCPServer{
		MCPName: "Default MCP Server",
		Version: "1.0.0",
		Port:    ":9080",
	}

	for _, opt := range options {
		opt(s)
	}

	serverOpts := make([]server.ServerOption, 0, len(s.toolFilters)+len(s.toolHandlerMiddlewares))
	for _, filter := range s.toolFilters {
		serverOpts = append(serverOpts, server.WithToolFilter(filter))
	}
	for _, mid := range s.toolHandlerMiddlewares {
		serverOpts = append(serverOpts, server.WithToolHandlerMiddleware(mid))
	}
	s.mcp = server.NewMCPServer(s.MCPName, s.Version, serverOpts...)
	return s
}

func (s *MCPServer) Name() string {
	return name
}

func (s *MCPServer) Start(ctx context.Context) error {

	s.mcp.AddTools(s.tools...)
	s.mcp.AddPrompts(s.prompts...)
	s.mcp.AddResources(s.resources...)

	log.Printf("Starting MCP server: %s on port %s", s.MCPName, s.Port)

	s.svr = server.NewStreamableHTTPServer(s.mcp)
	if err := s.svr.Start(s.Port); err != nil && err != context.Canceled && err != http.ErrServerClosed {
		return errors.Wrap(err, "failed to start MCP server")
	}
	return nil
}

func (s *MCPServer) Stop(ctx context.Context) error {
	return s.svr.Shutdown(ctx)
}
