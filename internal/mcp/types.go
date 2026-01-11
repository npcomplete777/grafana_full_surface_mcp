package mcp

import "encoding/json"

// JSON-RPC Request
type Request struct {
	JSONRPC string          `json:"jsonrpc"`
	ID      json.RawMessage `json:"id"`
	Method  string          `json:"method"`
	Params  interface{}     `json:"params,omitempty"`
}

// JSON-RPC Response
type Response struct {
	JSONRPC string          `json:"jsonrpc"`
	ID      json.RawMessage `json:"id"`
	Result  interface{}     `json:"result,omitempty"`
	Error   *Error          `json:"error,omitempty"`
}

// JSON-RPC Error with structured data per JSON-RPC 2.0 spec
type Error struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    *ErrorData  `json:"data,omitempty"`
}

// ErrorData provides structured error details
type ErrorData struct {
	Details string `json:"details,omitempty"`
	Cause   string `json:"cause,omitempty"`
}

// MCP Tool Definition
type Tool struct {
	Name        string            `json:"name"`
	Description string            `json:"description"`
	InputSchema InputSchema       `json:"inputSchema"`
	Annotations *ToolAnnotations  `json:"annotations,omitempty"`
}

type InputSchema struct {
	Type       string              `json:"type"`
	Properties map[string]Property `json:"properties,omitempty"`
	Required   []string            `json:"required,omitempty"`
}

type Property struct {
	Type        string   `json:"type"`
	Description string   `json:"description"`
	Enum        []string `json:"enum,omitempty"`
	Default     interface{} `json:"default,omitempty"`
}

type ToolAnnotations struct {
	ReadOnlyHint     bool `json:"readOnlyHint,omitempty"`
	DestructiveHint  bool `json:"destructiveHint,omitempty"`
	IdempotentHint   bool `json:"idempotentHint,omitempty"`
	OpenWorldHint    bool `json:"openWorldHint,omitempty"`
}

// Tool Call Request
type CallToolParams struct {
	Name      string                 `json:"name"`
	Arguments map[string]interface{} `json:"arguments,omitempty"`
}

// Tool Call Response
type CallToolResult struct {
	Content []ContentBlock `json:"content"`
	IsError bool           `json:"isError,omitempty"`
}

type ContentBlock struct {
	Type string `json:"type"`
	Text string `json:"text,omitempty"`
}

// Initialize Result
type InitializeResult struct {
	ProtocolVersion string       `json:"protocolVersion"`
	Capabilities    Capabilities `json:"capabilities"`
	ServerInfo      ServerInfo   `json:"serverInfo"`
}

type Capabilities struct {
	Tools     *ToolsCapability     `json:"tools,omitempty"`
	Resources *ResourcesCapability `json:"resources,omitempty"`
}

type ToolsCapability struct {
	ListChanged bool `json:"listChanged,omitempty"`
}

type ResourcesCapability struct {
	Subscribe   bool `json:"subscribe,omitempty"`
	ListChanged bool `json:"listChanged,omitempty"`
}

type ServerInfo struct {
	Name    string `json:"name"`
	Version string `json:"version"`
}

// List Tools Result
type ListToolsResult struct {
	Tools []Tool `json:"tools"`
}

// Standard error codes
const (
	ParseError     = -32700
	InvalidRequest = -32600
	MethodNotFound = -32601
	InvalidParams  = -32602
	InternalError  = -32603
)
