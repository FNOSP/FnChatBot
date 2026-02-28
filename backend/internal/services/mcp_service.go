package services

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"sync"
	"time"

	"fnchatbot/internal/models"

	"github.com/mark3labs/mcp-go/client"
	"github.com/mark3labs/mcp-go/mcp"
)

const (
	defaultTimeoutMs = 5000
)

// DefaultMCPService is set by main; tool_service and handlers use it when non-nil.
var DefaultMCPService *MCPService

// MCPService manages MCP config file, client connections, and status cache.
type MCPService struct {
	filePath string
	mu       sync.RWMutex
	status   map[string]models.MCPStatus
	clients  map[string]*client.Client
}

// NewMCPService creates an MCPService. filePath defaults to "mcp.json" if empty.
func NewMCPService(filePath string) *MCPService {
	if filePath == "" {
		filePath = "mcp.json"
	}
	return &MCPService{
		filePath: filePath,
		status:   make(map[string]models.MCPStatus),
		clients:  make(map[string]*client.Client),
	}
}

// LoadFile reads mcp.json and returns the config. Returns empty servers map if file missing or invalid.
func (s *MCPService) LoadFile() (*models.MCPFile, error) {
	data, err := os.ReadFile(s.filePath)
	if err != nil {
		if os.IsNotExist(err) {
			return &models.MCPFile{Servers: make(map[string]models.MCPServerConfig)}, nil
		}
		return nil, err
	}
	var f models.MCPFile
	if err := json.Unmarshal(data, &f); err != nil {
		return nil, err
	}
	if f.Servers == nil {
		f.Servers = make(map[string]models.MCPServerConfig)
	}
	return &f, nil
}

// SaveFile writes the config to mcp.json.
func (s *MCPService) SaveFile(f *models.MCPFile) error {
	if f == nil || f.Servers == nil {
		f = &models.MCPFile{Servers: make(map[string]models.MCPServerConfig)}
	}
	data, err := json.MarshalIndent(f, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(s.filePath, data, 0644)
}

func (s *MCPService) timeoutMs(cfg models.MCPServerConfig) int {
	if cfg.Timeout > 0 {
		return cfg.Timeout
	}
	return defaultTimeoutMs
}

// envSlice converts map to "KEY=VALUE" slice for stdio client.
func envSlice(env map[string]string) []string {
	if len(env) == 0 {
		return nil
	}
	out := make([]string, 0, len(env))
	for k, v := range env {
		out = append(out, k+"="+v)
	}
	return out
}

// connectServer creates MCP client for the given config, runs Initialize + ListTools, stores client and sets status.
func (s *MCPService) connectServer(ctx context.Context, name string, cfg models.MCPServerConfig) {
	timeout := time.Duration(s.timeoutMs(cfg)) * time.Millisecond
	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	var c *client.Client
	var err error
	switch cfg.Type {
	case models.MCPTypeLocal:
		if len(cfg.Command) == 0 {
			s.setStatus(name, models.MCPStatusFailed, "local MCP requires non-empty command")
			return
		}
		cmd := cfg.Command[0]
		args := cfg.Command[1:]
		env := envSlice(cfg.Env)
		c, err = client.NewStdioMCPClient(cmd, env, args...)
	case models.MCPTypeRemote:
		// Prefer Streamable HTTP; fallback to SSE if URL invalid for one
		c, err = client.NewStreamableHttpClient(cfg.URL, nil)
		if err != nil {
			c, err = client.NewSSEMCPClient(cfg.URL)
		}
	default:
		s.setStatus(name, models.MCPStatusFailed, "unknown type: "+string(cfg.Type))
		return
	}
	if err != nil {
		s.setStatus(name, models.MCPStatusFailed, err.Error())
		return
	}
	// Remote clients need Start then Initialize; Stdio starts automatically but may need Initialize
	if !c.IsInitialized() {
		if err := c.Start(ctx); err != nil {
			c.Close()
			s.setStatus(name, models.MCPStatusFailed, err.Error())
			return
		}
		req := mcp.InitializeRequest{
			Params: mcp.InitializeParams{
				ProtocolVersion: "2024-11-05",
				Capabilities:    mcp.ClientCapabilities{},
				ClientInfo:      mcp.Implementation{Name: "fnchatbot", Version: "0.1.0"},
			},
		}
		if _, err := c.Initialize(ctx, req); err != nil {
			c.Close()
			s.setStatus(name, models.MCPStatusFailed, err.Error())
			return
		}
	}
	_, err = c.ListTools(ctx, mcp.ListToolsRequest{})
	if err != nil {
		c.Close()
		s.setStatus(name, models.MCPStatusFailed, err.Error())
		return
	}
	s.mu.Lock()
	if old := s.clients[name]; old != nil {
		old.Close()
	}
	s.clients[name] = c
	s.status[name] = models.MCPStatus{Status: models.MCPStatusConnected}
	s.mu.Unlock()
}

func (s *MCPService) setStatus(name string, status models.MCPStatusType, errMsg string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	st := models.MCPStatus{Status: status}
	if errMsg != "" {
		st.Error = errMsg
	}
	s.status[name] = st
}

// disconnectServer closes the client for name and removes it from the pool.
func (s *MCPService) disconnectServer(name string) {
	s.mu.Lock()
	c := s.clients[name]
	delete(s.clients, name)
	s.mu.Unlock()
	if c != nil {
		_ = c.Close()
	}
}

// CheckServer checks a single server: if disabled set status disabled; else (re)connect and update status.
func (s *MCPService) CheckServer(ctx context.Context, name string) models.MCPStatus {
	f, err := s.LoadFile()
	if err != nil {
		return models.MCPStatus{Status: models.MCPStatusFailed, Error: err.Error()}
	}
	cfg, ok := f.Servers[name]
	if !ok {
		s.mu.Lock()
		delete(s.status, name)
		s.mu.Unlock()
		return models.MCPStatus{Status: models.MCPStatusFailed, Error: "server not found"}
	}
	if !cfg.Enabled {
		s.disconnectServer(name)
		s.mu.Lock()
		s.status[name] = models.MCPStatus{Status: models.MCPStatusDisabled}
		s.mu.Unlock()
		return models.MCPStatus{Status: models.MCPStatusDisabled}
	}
	s.disconnectServer(name)
	s.connectServer(ctx, name, cfg)
	s.mu.RLock()
	st := s.status[name]
	s.mu.RUnlock()
	return st
}

// CheckAllEnabled runs CheckServer for every enabled server in config, concurrently.
func (s *MCPService) CheckAllEnabled(ctx context.Context) map[string]models.MCPStatus {
	f, err := s.LoadFile()
	if err != nil {
		log.Printf("MCP LoadFile: %v", err)
		return nil
	}
	var wg sync.WaitGroup
	for name, cfg := range f.Servers {
		if !cfg.Enabled {
			s.mu.Lock()
			s.status[name] = models.MCPStatus{Status: models.MCPStatusDisabled}
			s.mu.Unlock()
			continue
		}
		wg.Add(1)
		go func(n string, c models.MCPServerConfig) {
			defer wg.Done()
			s.CheckServer(ctx, n)
		}(name, cfg)
	}
	wg.Wait()
	s.mu.RLock()
	out := make(map[string]models.MCPStatus, len(s.status))
	for k, v := range s.status {
		out[k] = v
	}
	s.mu.RUnlock()
	return out
}

// GetStatus returns a copy of the current status map (all configured servers get a status, default unknown).
func (s *MCPService) GetStatus() map[string]models.MCPStatus {
	f, _ := s.LoadFile()
	if f == nil {
		f = &models.MCPFile{Servers: make(map[string]models.MCPServerConfig)}
	}
	s.mu.RLock()
	defer s.mu.RUnlock()
	out := make(map[string]models.MCPStatus)
	for name := range f.Servers {
		if st, ok := s.status[name]; ok {
			out[name] = st
		} else {
			out[name] = models.MCPStatus{Status: models.MCPStatusUnknown}
		}
	}
	return out
}

// GetClient returns the MCP client for name, or nil if not connected.
func (s *MCPService) GetClient(name string) *client.Client {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.clients[name]
}

// GetConnectedClients returns a snapshot of name -> client for all connected servers.
func (s *MCPService) GetConnectedClients() map[string]*client.Client {
	s.mu.RLock()
	defer s.mu.RUnlock()
	out := make(map[string]*client.Client, len(s.clients))
	for k, v := range s.clients {
		out[k] = v
	}
	return out
}

// Shutdown closes all MCP clients (e.g. stdio subprocesses).
func (s *MCPService) Shutdown() {
	s.mu.Lock()
	clients := s.clients
	s.clients = make(map[string]*client.Client)
	s.mu.Unlock()
	for _, c := range clients {
		if err := c.Close(); err != nil {
			log.Printf("MCP client close: %v", err)
		}
	}
}

// ListServerNames returns all configured server names (from file).
func (s *MCPService) ListServerNames() ([]string, error) {
	f, err := s.LoadFile()
	if err != nil {
		return nil, err
	}
	names := make([]string, 0, len(f.Servers))
	for name := range f.Servers {
		names = append(names, name)
	}
	return names, nil
}

// GetServerConfig returns config for name, or nil if not found.
func (s *MCPService) GetServerConfig(name string) (*models.MCPServerConfig, error) {
	f, err := s.LoadFile()
	if err != nil {
		return nil, err
	}
	cfg, ok := f.Servers[name]
	if !ok {
		return nil, fmt.Errorf("MCP server %q not found", name)
	}
	return &cfg, nil
}

// SetServer adds or updates a server in the config file and optionally runs a check if enabled.
func (s *MCPService) SetServer(name string, cfg models.MCPServerConfig) error {
	f, err := s.LoadFile()
	if err != nil {
		return err
	}
	if f.Servers == nil {
		f.Servers = make(map[string]models.MCPServerConfig)
	}
	f.Servers[name] = cfg
	if err := s.SaveFile(f); err != nil {
		return err
	}
	if cfg.Enabled {
		go func() {
			ctx, cancel := context.WithTimeout(context.Background(), time.Duration(s.timeoutMs(cfg)+500)*time.Millisecond)
			defer cancel()
			s.CheckServer(ctx, name)
		}()
	} else {
		s.disconnectServer(name)
		s.mu.Lock()
		s.status[name] = models.MCPStatus{Status: models.MCPStatusDisabled}
		s.mu.Unlock()
	}
	return nil
}

// DeleteServer removes a server from config and disconnects its client.
func (s *MCPService) DeleteServer(name string) error {
	f, err := s.LoadFile()
	if err != nil {
		return err
	}
	if _, ok := f.Servers[name]; !ok {
		return fmt.Errorf("MCP server %q not found", name)
	}
	delete(f.Servers, name)
	if err := s.SaveFile(f); err != nil {
		return err
	}
	s.disconnectServer(name)
	s.mu.Lock()
	delete(s.status, name)
	s.mu.Unlock()
	return nil
}
