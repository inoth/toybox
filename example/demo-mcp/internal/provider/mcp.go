package provider

import (
	"context"
	"log"
	"mcp-quickstart/internal/handler"

	"github.com/inoth/toybox/mcpserver"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

func NewMCPServer(ipQuery *handler.IPQueryHandler) *mcpserver.MCPServer {
	return mcpserver.NewMCPServer(
		mcpserver.WithTools(ipQuery.Tools()...),
		mcpserver.WithPrompts(),
		mcpserver.WithResources(),
		mcpserver.WithToolHandlerMiddleware(
			func(next server.ToolHandlerFunc) server.ToolHandlerFunc {
				return func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
					log.Println("Handling tool request")
					return next(ctx, req)
				}
			}),
	)
}
