package handler

import (
	"context"
	"errors"
	"fmt"
	"io"
	"log"

	"net"
	"net/http"

	"github.com/inoth/toybox/mcpserver"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

var _ mcpserver.MCPHandler = (*IPQueryHandler)(nil)

type IPQueryHandler struct{}

func NewIPQueryHandler() *IPQueryHandler {
	return &IPQueryHandler{}
}

func (h *IPQueryHandler) Tools() []server.ServerTool {
	return []server.ServerTool{
		{
			Tool: mcp.NewTool("ip_query",
				mcp.WithDescription("Query IP address information"),
				mcp.WithString("ip",
					mcp.Required(),
					mcp.Description("IP address to query"),
				)),
			Handler: h.IPQuery,
		},
	}
}

func (h *IPQueryHandler) Resources() []server.ServerResource { return nil }

func (h *IPQueryHandler) Prompts() []server.ServerPrompt { return nil }

func (h *IPQueryHandler) IPQuery(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	ip := request.GetString("ip", "")
	if ip == "" {
		return nil, errors.New("IP address is required")
	}

	parsedIP := net.ParseIP(ip)
	if parsedIP == nil {
		log.Printf("invalid IP address: %s", ip)
		return nil, errors.New("invalid IP address")
	}

	resp, err := http.Get(fmt.Sprintf("https://ip.rpcx.io/api/ip?ip=%s", ip))
	if err != nil {
		log.Printf("Error fetching IP information: %v", err)
		return nil, fmt.Errorf("Error fetching IP information: %v", err)
	}
	defer resp.Body.Close()

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Printf("Error reading response body: %v", err)
		return nil, fmt.Errorf("Error reading response body: %v", err)
	}

	return mcp.NewToolResultText(string(data)), nil
}
