package models

// MCPType distinguishes local (stdio) vs remote (HTTP/SSE) MCP servers.
type MCPType string

const (
	MCPTypeLocal  MCPType = "local"
	MCPTypeRemote MCPType = "remote"
)

// MCPServerConfig is the per-server config stored in mcp.json.
// Local: Command, Env; Remote: URL, ApiKey, Headers; Common: Enabled, Timeout.
type MCPServerConfig struct {
	Type MCPType `json:"type"`
	// local (stdio)
	Command []string          `json:"command,omitempty"`
	Env     map[string]string `json:"env,omitempty"`
	// remote (HTTP/SSE)
	URL     string            `json:"url,omitempty"`
	ApiKey  string            `json:"api_key,omitempty"`
	Headers map[string]string `json:"headers,omitempty"`
	// common
	Enabled bool `json:"enabled"`
	Timeout int  `json:"timeout,omitempty"` // ms, 0 = use default
}

// MCPFile is the root structure of mcp.json.
type MCPFile struct {
	Servers map[string]MCPServerConfig `json:"servers"`
}

// MCPStatusType is the runtime status of an MCP server.
type MCPStatusType string

const (
	MCPStatusConnected MCPStatusType = "connected"
	MCPStatusDisabled  MCPStatusType = "disabled"
	MCPStatusFailed    MCPStatusType = "failed"
	MCPStatusUnknown   MCPStatusType = "unknown"
)

// MCPStatus holds status and optional error message.
type MCPStatus struct {
	Status MCPStatusType `json:"status"`
	Error  string        `json:"error,omitempty"`
}

// MCPServerInfo is returned by GET /mcp: config + runtime status merged.
type MCPServerInfo struct {
	Name string `json:"name"`
	MCPServerConfig
	MCPStatus
}
