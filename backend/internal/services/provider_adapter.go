package services

import (
	"bufio"
	"context"
	"net/http"

	"fnchatbot/internal/models"
)

type ModelInfo struct {
	ID           string                   `json:"id"`
	Name         string                   `json:"name"`
	Object       string                   `json:"object,omitempty"`
	OwnedBy      string                   `json:"owned_by,omitempty"`
	Description  string                   `json:"description,omitempty"`
	Capabilities []models.ModelCapability `json:"capabilities,omitempty"`
}

type ChatMessage struct {
	Role       string     `json:"role"`
	Content    string     `json:"content"`
	ToolCalls  []ToolCall `json:"tool_calls,omitempty"`
	ToolCallID string     `json:"tool_call_id,omitempty"`
	Name       string     `json:"name,omitempty"`
}

type ToolCall struct {
	Index    int      `json:"index"`
	ID       string   `json:"id,omitempty"`
	Type     string   `json:"type"`
	Function ToolFunc `json:"function"`
}

type ToolFunc struct {
	Name      string `json:"name"`
	Arguments string `json:"arguments"`
}

type Tool struct {
	Type     string     `json:"type"`
	Function ToolSchema `json:"function"`
}

type ToolSchema struct {
	Name        string                 `json:"name"`
	Description string                 `json:"description"`
	Parameters  map[string]interface{} `json:"parameters"`
}

type ChatOptions struct {
	APIKey      string  `json:"api_key,omitempty"`
	BaseURL     string  `json:"base_url,omitempty"`
	Temperature float32 `json:"temperature,omitempty"`
	MaxTokens   int     `json:"max_tokens,omitempty"`
	Stream      bool    `json:"stream"`
	Tools       []Tool  `json:"tools,omitempty"`
}

type StreamChunk struct {
	Content      string     `json:"content,omitempty"`
	ToolCalls    []ToolCall `json:"tool_calls,omitempty"`
	FinishReason string     `json:"finish_reason,omitempty"`
	Done         bool       `json:"done"`
}

type ProviderAdapter interface {
	ProviderType() models.ProviderType
	FetchModels(ctx context.Context, baseURL, apiKey string) ([]ModelInfo, error)
	BuildChatRequest(ctx context.Context, model string, messages []ChatMessage, options ChatOptions) (*http.Request, error)
	ParseStreamResponse(ctx context.Context, reader *bufio.Reader) (StreamChunk, error)
	GetBaseURL(baseURL string) string
}
