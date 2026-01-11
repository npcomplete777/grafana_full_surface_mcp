package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"

	"github.com/aaronjacobs/grafana-mcp-server/internal/grafana"
	"github.com/aaronjacobs/grafana-mcp-server/internal/mcp"
	"github.com/aaronjacobs/grafana-mcp-server/internal/tools"
)

const (
	serverName    = "grafana-mcp-server"
	serverVersion = "1.0.0"
	protocolVersion = "2024-11-05"
)

// Server handles MCP protocol communication
type Server struct {
	registry *tools.Registry
	reader   *bufio.Reader
	writer   io.Writer
}

func main() {
	// Get configuration from environment
	grafanaURL := os.Getenv("GRAFANA_URL")
	if grafanaURL == "" {
		grafanaURL = "http://localhost:3000"
	}

	apiKey := os.Getenv("GRAFANA_API_KEY")
	if apiKey == "" {
		log.Println("Warning: GRAFANA_API_KEY not set, some operations may fail")
	}

	// Create Grafana client
	client := grafana.NewClient(grafanaURL, apiKey)

	// Create tool registry
	registry := tools.NewRegistry(client)

	// Create server
	server := &Server{
		registry: registry,
		reader:   bufio.NewReader(os.Stdin),
		writer:   os.Stdout,
	}

	log.SetOutput(os.Stderr)
	log.Printf("Starting %s v%s", serverName, serverVersion)
	log.Printf("Grafana URL: %s", grafanaURL)

	// Run the server
	if err := server.Run(); err != nil {
		log.Fatalf("Server error: %v", err)
	}
}

// Run starts the main server loop
func (s *Server) Run() error {
	for {
		line, err := s.reader.ReadBytes('\n')
		if err != nil {
			if err == io.EOF {
				return nil
			}
			return fmt.Errorf("read error: %w", err)
		}

		if len(line) == 0 {
			continue
		}

		var request mcp.Request
		if err := json.Unmarshal(line, &request); err != nil {
			// For parse errors, ID must be null per JSON-RPC 2.0 spec
			s.sendError(json.RawMessage("null"), mcp.ParseError, "Parse error", err.Error())
			continue
		}

		// Skip notifications (requests without an ID)
		if len(request.ID) == 0 || string(request.ID) == "null" {
			// Still handle the method for side effects, but don't respond
			s.handleRequest(&request)
			continue
		}

		s.handleRequest(&request)
	}
}

func (s *Server) handleRequest(req *mcp.Request) {
	// Check if this is a notification (no ID means no response expected)
	isNotification := len(req.ID) == 0 || string(req.ID) == "null"

	switch req.Method {
	case "initialize":
		if !isNotification {
			s.handleInitialize(req)
		}
	case "initialized", "notifications/cancelled":
		// Notifications - no response needed
	case "tools/list":
		if !isNotification {
			s.handleListTools(req)
		}
	case "tools/call":
		if !isNotification {
			s.handleCallTool(req)
		}
	case "ping":
		if !isNotification {
			s.sendResult(req.ID, map[string]string{})
		}
	default:
		if !isNotification {
			s.sendError(req.ID, mcp.MethodNotFound, "Method not found", req.Method)
		}
	}
}

func (s *Server) handleInitialize(req *mcp.Request) {
	result := mcp.InitializeResult{
		ProtocolVersion: protocolVersion,
		Capabilities: mcp.Capabilities{
			Tools: &mcp.ToolsCapability{
				ListChanged: false,
			},
		},
		ServerInfo: mcp.ServerInfo{
			Name:    serverName,
			Version: serverVersion,
		},
	}
	s.sendResult(req.ID, result)
}

func (s *Server) handleListTools(req *mcp.Request) {
	result := mcp.ListToolsResult{
		Tools: s.registry.GetTools(),
	}
	s.sendResult(req.ID, result)
}

func (s *Server) handleCallTool(req *mcp.Request) {
	// Parse params
	paramsJSON, err := json.Marshal(req.Params)
	if err != nil {
		s.sendError(req.ID, mcp.InvalidParams, "Invalid params", err.Error())
		return
	}

	var params mcp.CallToolParams
	if err := json.Unmarshal(paramsJSON, &params); err != nil {
		s.sendError(req.ID, mcp.InvalidParams, "Invalid params", err.Error())
		return
	}

	result, err := s.registry.CallTool(params.Name, params.Arguments)
	if err != nil {
		s.sendError(req.ID, mcp.InternalError, "Tool execution failed", err.Error())
		return
	}

	s.sendResult(req.ID, result)
}

func (s *Server) sendResult(id json.RawMessage, result interface{}) {
	response := mcp.Response{
		JSONRPC: "2.0",
		ID:      id,
		Result:  result,
	}
	s.send(response)
}

func (s *Server) sendError(id json.RawMessage, code int, message, details string) {
	response := mcp.Response{
		JSONRPC: "2.0",
		ID:      id,
		Error: &mcp.Error{
			Code:    code,
			Message: message,
			Data:    &mcp.ErrorData{Details: details},
		},
	}
	s.send(response)
}

func (s *Server) send(response mcp.Response) {
	data, err := json.Marshal(response)
	if err != nil {
		log.Printf("Failed to marshal response: %v", err)
		return
	}
	fmt.Fprintf(s.writer, "%s\n", data)
}
