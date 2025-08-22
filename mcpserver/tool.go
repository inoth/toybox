package mcpserver

import (
	"github.com/mark3labs/mcp-go/server"
)

type MCPHandler interface {
	Tools() []server.ServerTool
	Resources() []server.ServerResource
	Prompts() []server.ServerPrompt
}
