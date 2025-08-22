package mcpserver

import "github.com/mark3labs/mcp-go/server"

type Option func(*MCPServer)

func WithName(name string) Option {
	return func(s *MCPServer) {
		s.MCPName = name
	}
}

func WithVersion(version string) Option {
	return func(s *MCPServer) {
		s.Version = version
	}
}

func WithPort(port string) Option {
	return func(s *MCPServer) {
		s.Port = port
	}
}

func WithResources(resources ...server.ServerResource) Option {
	return func(s *MCPServer) {
		s.resources = append(s.resources, resources...)
	}
}

func WithTools(tools ...server.ServerTool) Option {
	return func(s *MCPServer) {
		s.tools = append(s.tools, tools...)
	}
}

func WithPrompts(prompts ...server.ServerPrompt) Option {
	return func(s *MCPServer) {
		s.prompts = append(s.prompts, prompts...)
	}
}

func WithToolFilter(filter ...server.ToolFilterFunc) Option {
	return func(s *MCPServer) {
		s.toolFilters = append(s.toolFilters, filter...)
	}
}

func WithToolHandlerMiddleware(middleware ...server.ToolHandlerMiddleware) Option {
	return func(s *MCPServer) {
		s.toolHandlerMiddlewares = append(s.toolHandlerMiddlewares, middleware...)
	}
}
