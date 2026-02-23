package adapters

import (
	"bufio"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"

	"fnchatbot/internal/models"
	"fnchatbot/internal/services"
)

const (
	anthropicVersion = "2023-06-01"
)

type AnthropicAdapter struct{}

func NewAnthropicAdapter() *AnthropicAdapter {
	return &AnthropicAdapter{}
}

func (a *AnthropicAdapter) ProviderType() models.ProviderType {
	return models.ProviderTypeAnthropic
}

func (a *AnthropicAdapter) FetchModels(ctx context.Context, baseURL, apiKey string) ([]services.ModelInfo, error) {
	return []services.ModelInfo{
		{ID: "claude-opus-4-6", Name: "Claude Opus 4.6"},
		{ID: "claude-opus-4-5", Name: "Claude Opus 4.5"},
		{ID: "claude-sonnet-4-5", Name: "Claude Sonnet 4.5"},
		{ID: "claude-haiku-4-5", Name: "Claude Haiku 4.5"},
		{ID: "claude-3-5-sonnet-20241022", Name: "Claude 3.5 Sonnet"},
		{ID: "claude-3-5-haiku-20241022", Name: "Claude 3.5 Haiku"},
		{ID: "claude-3-opus-20240229", Name: "Claude 3 Opus"},
		{ID: "claude-3-sonnet-20240229", Name: "Claude 3 Sonnet"},
		{ID: "claude-3-haiku-20240307", Name: "Claude 3 Haiku"},
	}, nil
}

type anthropicRequest struct {
	Model       string             `json:"model"`
	MaxTokens   int                `json:"max_tokens"`
	Messages    []anthropicMessage `json:"messages"`
	System      string             `json:"system,omitempty"`
	Stream      bool               `json:"stream,omitempty"`
	Temperature float32            `json:"temperature,omitempty"`
	Tools       []anthropicTool    `json:"tools,omitempty"`
}

type anthropicMessage struct {
	Role    string           `json:"role"`
	Content anthropicContent `json:"content"`
}

type anthropicContent interface{}

type anthropicTextContent struct {
	Type string `json:"type"`
	Text string `json:"text"`
}

type anthropicToolResultContent struct {
	Type      string `json:"type"`
	ToolUseID string `json:"tool_use_id"`
	Content   string `json:"content"`
}

type anthropicTool struct {
	Name        string                 `json:"name"`
	Description string                 `json:"description"`
	InputSchema map[string]interface{} `json:"input_schema"`
}

func (a *AnthropicAdapter) BuildChatRequest(ctx context.Context, model string, messages []services.ChatMessage, options services.ChatOptions) (*http.Request, error) {
	var systemPrompt string
	var anthropicMessages []anthropicMessage

	for _, msg := range messages {
		switch msg.Role {
		case "system":
			systemPrompt = msg.Content
		case "user":
			anthropicMessages = append(anthropicMessages, anthropicMessage{
				Role: "user",
				Content: anthropicTextContent{
					Type: "text",
					Text: msg.Content,
				},
			})
		case "assistant":
			content := make([]interface{}, 0)
			if msg.Content != "" {
				content = append(content, anthropicTextContent{
					Type: "text",
					Text: msg.Content,
				})
			}
			for _, tc := range msg.ToolCalls {
				content = append(content, map[string]interface{}{
					"type":  "tool_use",
					"id":    tc.ID,
					"name":  tc.Function.Name,
					"input": json.RawMessage(tc.Function.Arguments),
				})
			}
			anthropicMessages = append(anthropicMessages, anthropicMessage{
				Role:    "assistant",
				Content: content,
			})
		case "tool":
			anthropicMessages = append(anthropicMessages, anthropicMessage{
				Role: "user",
				Content: anthropicToolResultContent{
					Type:      "tool_result",
					ToolUseID: msg.ToolCallID,
					Content:   msg.Content,
				},
			})
		}
	}

	maxTokens := options.MaxTokens
	if maxTokens <= 0 {
		maxTokens = 4096
	}

	reqBody := anthropicRequest{
		Model:     model,
		MaxTokens: maxTokens,
		Messages:  anthropicMessages,
		System:    systemPrompt,
		Stream:    options.Stream,
	}

	if options.Temperature > 0 {
		reqBody.Temperature = options.Temperature
	}

	if len(options.Tools) > 0 {
		tools := make([]anthropicTool, len(options.Tools))
		for i, tool := range options.Tools {
			tools[i] = anthropicTool{
				Name:        tool.Function.Name,
				Description: tool.Function.Description,
				InputSchema: tool.Function.Parameters,
			}
		}
		reqBody.Tools = tools
	}

	jsonBody, err := json.Marshal(reqBody)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	baseURL := a.GetBaseURL(options.BaseURL)
	url := fmt.Sprintf("%s/messages", baseURL)

	req, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewBuffer(jsonBody))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("x-api-key", options.APIKey)
	req.Header.Set("anthropic-version", anthropicVersion)

	return req, nil
}

type anthropicStreamEvent struct {
	Type  string `json:"type"`
	Index int    `json:"index,omitempty"`
	Delta *struct {
		Type string `json:"type"`
		Text string `json:"text,omitempty"`
	} `json:"delta,omitempty"`
	ContentBlock *struct {
		Type string `json:"type"`
		Text string `json:"text,omitempty"`
		ID   string `json:"id,omitempty"`
		Name string `json:"name,omitempty"`
	} `json:"content_block,omitempty"`
	ToolUse *struct {
		ID    string          `json:"id"`
		Name  string          `json:"name"`
		Input json.RawMessage `json:"input"`
	} `json:"tool_use,omitempty"`
	Message *struct {
		ID           string `json:"id"`
		StopReason   string `json:"stop_reason"`
		StopSequence string `json:"stop_sequence"`
	} `json:"message,omitempty"`
}

func (a *AnthropicAdapter) ParseStreamResponse(ctx context.Context, reader *bufio.Reader) (services.StreamChunk, error) {
	for {
		line, err := reader.ReadString('\n')
		if err != nil {
			if err == io.EOF {
				return services.StreamChunk{Done: true}, nil
			}
			return services.StreamChunk{}, fmt.Errorf("error reading stream: %w", err)
		}

		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}

		if !strings.HasPrefix(line, "data: ") {
			continue
		}

		data := strings.TrimPrefix(line, "data: ")
		if data == "" {
			continue
		}

		var event anthropicStreamEvent
		if err := json.Unmarshal([]byte(data), &event); err != nil {
			continue
		}

		switch event.Type {
		case "content_block_delta":
			if event.Delta != nil && event.Delta.Type == "text_delta" {
				return services.StreamChunk{
					Content: event.Delta.Text,
				}, nil
			}

		case "content_block_start":
			if event.ContentBlock != nil && event.ContentBlock.Type == "tool_use" {
				return services.StreamChunk{
					ToolCalls: []services.ToolCall{
						{
							Index: event.Index,
							ID:    event.ContentBlock.ID,
							Type:  "function",
							Function: services.ToolFunc{
								Name: event.ContentBlock.Name,
							},
						},
					},
				}, nil
			}

		case "input_json_delta":
			return services.StreamChunk{
				ToolCalls: []services.ToolCall{
					{
						Index: event.Index,
						Function: services.ToolFunc{
							Arguments: event.Delta.Text,
						},
					},
				},
			}, nil

		case "message_delta":
			if event.Message != nil && event.Message.StopReason != "" {
				return services.StreamChunk{
					FinishReason: event.Message.StopReason,
				}, nil
			}

		case "message_stop":
			return services.StreamChunk{Done: true}, nil

		case "error":
			return services.StreamChunk{}, fmt.Errorf("anthropic API error: %s", data)
		}
	}
}

func (a *AnthropicAdapter) GetBaseURL(baseURL string) string {
	baseURL = strings.TrimSuffix(baseURL, "/")
	if baseURL == "" {
		baseURL = "https://api.anthropic.com"
	}
	if !strings.HasSuffix(baseURL, "/v1") {
		baseURL += "/v1"
	}
	return baseURL
}
